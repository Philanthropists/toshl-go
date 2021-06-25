package toshl

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// HTTPClient is an interface to define the client to access API resources
type HTTPClient interface {
	Get(APIUrl, queryString string) (string, string, error)
	GetMultiple(APIUrl, queryString string) ([]string, error)
	Post(APIUrl, JSONPayload string) (string, error)
	Update(APIUrl, JSONPayload string) (string, error)
	Delete(APIUrl string) error
}

// RestHTTPClient is a real implementation of the HTTPClient
type RestHTTPClient struct {
	BaseURL string
	Token   string
	Client  *http.Client
}

func (c *RestHTTPClient) setAuthenticationHeader(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
}

func (c *RestHTTPClient) setJSONContentTypeHeader(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
}

func (c *RestHTTPClient) setUserAgentHeader(req *http.Request) {
	req.Header.Set("User-Agent", GetUserAgentString())
}

func (c *RestHTTPClient) getIDFromLocationHeader(
	response *http.Response) (string, error) {

	locationHeader := response.Header.Get("Location")

	id, err := c.parseIDFromLocationHeader(locationHeader)
	if err != nil {
		log.Print("Location URL parsing: ", err)
		return "", err
	}

	return id, nil
}

func (c *RestHTTPClient) parseIDFromLocationHeader(
	locationURL string) (string, error) {

	guid, err := url.Parse(locationURL)
	if err != nil {
		log.Print("Location URL parsing: ", err)
		return "", err
	}

	values := strings.Split(guid.Path, "/")

	if len(values) > 1 {
		id := values[len(values)-1]
		return id, nil
	}

	return "", errors.New("cannot parse resource Id")
}

// Get takes an API endpoint and return a JSON string
func (c *RestHTTPClient) Get(APIUrl, queryString string) (string, string, error) {
	url := c.BaseURL + "/" + APIUrl

	if queryString != "" {
		url = url + "?" + queryString
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print("NewRequest: ", err)
		return "", "", err
	}

	// Set authorization token
	c.setAuthenticationHeader(req)

	// Set User-Agent header
	c.setUserAgentHeader(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Print("Do: ", err)
		return "", "", err
	}
	defer resp.Body.Close()

	link := resp.Header.Get("Link")

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("ReadAll: ", err)
		return "", "", err
	}

	return string(bs), link, nil
}

func extractQueryLink(match []byte) []byte {
	re := regexp.MustCompile(`\?[^<>]+`)
	return re.Find(match)
}

func getNextLink(links string) string {
	re := regexp.MustCompile(`<([^<>]*)>; rel="next"`)
	match := re.Find([]byte(links))
	if match == nil {
		return ""
	}

	match = extractQueryLink(match)

	return strings.Trim(string(match), "?")
}

func (c *RestHTTPClient) GetMultiple(APIUrl, queryString string) ([]string, error) {
	response, links, err := c.Get(APIUrl, queryString)
	if err != nil {
		return nil, err
	}

	link := getNextLink(links)
	responses := []string{response}

	for link != "" {
		response, nextLinks, err := c.Get(APIUrl, link)
		if err != nil {
			return nil, err
		}

		responses = append(responses, response)
		link = getNextLink(nextLinks)
	}

	return responses, nil
}

// Post takes an API endpoint and a JSON payload and return string Id
func (c *RestHTTPClient) Post(APIUrl, JSONPayload string) (string, error) {
	url := c.BaseURL + "/" + APIUrl
	var jsonStr = []byte(JSONPayload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Print("NewRequest: ", err)
		return "", err
	}

	// Set authorization token
	c.setAuthenticationHeader(req)

	// Set JSON content type
	c.setJSONContentTypeHeader(req)

	// Set User-Agent header
	c.setUserAgentHeader(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Print("Do: ", err)
		return "", err
	}

	defer resp.Body.Close()

	// Parse Location header to get Id
	id, err := c.getIDFromLocationHeader(resp)
	if err != nil {
		log.Print("Do: ", err)
		return "", err
	}

	return id, nil
}

// Update takes an API endpoint and a JSON payload and update the resource
func (c *RestHTTPClient) Update(APIUrl, JSONPayload string) (string, error) {
	url := c.BaseURL + "/" + APIUrl
	var jsonStr = []byte(JSONPayload)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Print("NewRequest: ", err)
		return "", err
	}

	// Set authorization token
	c.setAuthenticationHeader(req)

	// Set JSON content type
	c.setJSONContentTypeHeader(req)

	// Set User-Agent header
	c.setUserAgentHeader(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Print("Do: ", err)
		return "", err
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("ReadAll: ", err)
		return "", err
	}

	return string(bs), nil
}

// Delete removes the Account having the Id specified in the endpoint
func (c *RestHTTPClient) Delete(APIUrl string) error {
	url := c.BaseURL + "/" + APIUrl

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Print("NewRequest: ", err)
		return err
	}

	// Set authorization token
	c.setAuthenticationHeader(req)

	// Set User-Agent header
	c.setUserAgentHeader(req)

	_, err = c.Client.Do(req)
	if err != nil {
		log.Print("Do: ", err)
		return err
	}

	return nil
}

func (c *RestHTTPClient) SetTimeoutSeconds(timeout int) {
	c.Client.Timeout = time.Duration(timeout) * time.Second
}

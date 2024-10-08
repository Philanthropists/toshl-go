package toshl

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
)

// DefaultBaseURL is ...
const (
	DefaultBaseURL = "https://api.toshl.com"
	ClientVersion  = "0.1"
)

// Client handles API requests
type Client struct {
	client HTTPClient
}

// NewClient returns a new Toshl client
func NewClient(token string, httpClient HTTPClient) *Client {
	baseURL, _ := url.Parse(DefaultBaseURL)

	if httpClient == nil {
		httpClient = &RestHTTPClient{
			Client:  &http.Client{},
			BaseURL: baseURL.String(),
			Token:   token,
		}
	}

	c := &Client{client: httpClient}
	return c
}

// GetHTTPClient returns internal HTTPClient
func (c *Client) GetHTTPClient() HTTPClient {
	return c.client
}

// GetUserAgentString returns the string for UserAgent
func GetUserAgentString() string {
	return fmt.Sprintf(
		"toshl-go %s - %s", ClientVersion, runtime.Version())
}

// Accounts returns the list of Accounts
func (c *Client) Accounts(params *AccountQueryParams) ([]Account, error) {
	queryString := ""

	if params != nil {
		queryString = params.getQueryString()
	}

	res, err := c.client.Get("accounts", queryString)
	if err != nil {
		log.Println("GET /accounts/: ", err)
		return nil, err
	}

	var accounts []Account

	err = json.Unmarshal([]byte(res), &accounts)
	if err != nil {
		log.Println("JSON: ", res)
		return nil, err
	}

	return accounts, nil
}

// GetAccount returns the a specific Account
func (c *Client) GetAccount(accountID string) (*Account, error) {
	res, err := c.client.Get(fmt.Sprintf("accounts/%s", accountID), "")
	if err != nil {
		log.Println(fmt.Sprintf("GET /accounts/%s: ", accountID), err)
		return nil, err
	}

	var account *Account

	err = json.Unmarshal([]byte(res), &account)
	if err != nil {
		log.Println("JSON: ", res)
		return nil, err
	}

	return account, nil
}

type CreateAccountParams struct {
	Name     string   `json:"name"`
	Currency Currency `json:"currency"`
}

// CreateAccount creates a Toshl Account
func (c *Client) CreateAccount(account CreateAccountParams) (string, error) {
	jsonBytes, err := json.Marshal(account)
	if err != nil {
		log.Println("CreateAccount: ", err)
		return "", err
	}

	jsonStr := string(jsonBytes)

	id, err := c.client.Post("accounts", jsonStr)
	if err != nil {
		log.Println("POST /accounts/ ", err)
		return "", err
	}

	return id, nil
}

// SearchAccount search for Account name and return an Account
func (c *Client) SearchAccount(accountName string) (*Account, error) {
	accounts, err := c.Accounts(nil)
	if err != nil {
		log.Println("GET /accounts/: ", err)
		return nil, err
	}

	for _, account := range accounts {
		if account.Name == accountName {
			return &account, nil
		}
	}

	return nil, nil
}

// UpdateAccount updates a Toshl Account
func (c *Client) UpdateAccount(account *Account) error {
	jsonBytes, err := json.Marshal(account)
	if err != nil {
		log.Println("CreateAccount: ", err)
		return err
	}

	jsonStr := string(jsonBytes)

	accountResponse, err := c.client.Update(
		fmt.Sprintf("accounts/%s", account.ID), jsonStr)
	if err != nil {
		log.Println("PUT /accounts/ ", err)
		return err
	}

	err = json.Unmarshal([]byte(accountResponse), account)
	if err != nil {
		log.Println("Cannot decode Account JSON")
		return err
	}

	return nil
}

// DeleteAccount deletes a Toshl Account
func (c *Client) DeleteAccount(account *Account) error {
	err := c.client.Delete(fmt.Sprintf("accounts/%s", account.ID))
	if err != nil {
		log.Print("DELETE /accounts/ ", err)
		return err
	}

	return nil
}

// MoveAccount move a Toshl Account to a different position
func (c *Client) MoveAccount(account *Account, position int) error {
	jsonStr := fmt.Sprintf(`{"position": %s}`, strconv.Itoa(position))

	_, err := c.client.Post(fmt.Sprintf("accounts/%s", account.ID), jsonStr)
	if err != nil {
		log.Print("POST /accounts/ ", err)
		return err
	}

	return nil
}

// ReorderAccounts change the order of Toshl accounts
func (c *Client) ReorderAccounts(order *AccountsOrderParams) error {
	jsonBytes, err := json.Marshal(order)
	if err != nil {
		log.Println("ReorderAccounts: ", err)
		return err
	}

	jsonStr := string(jsonBytes)

	_, err = c.client.Post("accounts/reorder", jsonStr)
	if err != nil {
		log.Print("POST /accounts/reorder ", err)
		return err
	}

	return nil
}

// MergeAccounts merges two ore more Toshl accounts into a single one
func (c *Client) MergeAccounts(order *AccountsMergeParams) error {
	jsonBytes, err := json.Marshal(order)
	if err != nil {
		log.Println("MergeAccounts: ", err)
		return err
	}

	jsonStr := string(jsonBytes)

	_, err = c.client.Post("accounts/merge", jsonStr)
	if err != nil {
		log.Print("POST /accounts/merge ", err)
		return err
	}

	return nil
}

// Budgets returns the list of Budgets
func (c *Client) Budgets(params *BudgetQueryParams) ([]Budget, error) {
	queryString := ""

	if params != nil {
		queryString = params.getQueryString()
	}

	res, err := c.client.Get("budgets", queryString)
	if err != nil {
		log.Print("GET /budgets/: ", err)
		return nil, err
	}

	var budgets []Budget

	err = json.Unmarshal([]byte(res), &budgets)
	if err != nil {
		log.Println("JSON: ", res)
		return nil, err
	}

	return budgets, nil
}

// GetBudget returns the a specific Budget
func (c *Client) GetBudget(budgetID string) (*Budget, error) {
	res, err := c.client.Get(fmt.Sprintf("budgets/%s", budgetID), "")
	if err != nil {
		log.Print(fmt.Sprintf("GET /budgets/%s: ", budgetID), err)
		return nil, err
	}

	var budget *Budget

	err = json.Unmarshal([]byte(res), &budget)
	if err != nil {
		log.Println("JSON: ", res)
		return nil, err
	}

	return budget, nil
}

// Categories returns the list of Categories
func (c *Client) Categories(params *CategoryQueryParams) ([]Category, error) {
	queryString := ""

	if params != nil {
		queryString = params.getQueryString()
	}

	res, err := c.client.Get("categories", queryString)
	if err != nil {
		log.Print("GET /categories/: ", err)
		return nil, err
	}

	var categories []Category

	err = json.Unmarshal([]byte(res), &categories)
	if err != nil {
		log.Println("JSON: ", res)
		return nil, err
	}

	return categories, nil
}

// GetCategory returns the a specific Category
func (c *Client) GetCategory(categoryID string) (*Category, error) {
	res, err := c.client.Get(fmt.Sprintf("categories/%s", categoryID), "")
	if err != nil {
		log.Print(fmt.Sprintf("GET /categories/%s: ", categoryID), err)
		return nil, err
	}

	var category *Category

	err = json.Unmarshal([]byte(res), &category)
	if err != nil {
		log.Println("JSON: ", res)
		return nil, err
	}

	return category, nil
}

// CreateCategory creates a Toshl Category
func (c *Client) CreateCategory(category *Category) error {
	jsonBytes, err := json.Marshal(category)
	if err != nil {
		log.Println("CeateCategory: ", err)
		return err
	}

	jsonStr := string(jsonBytes)

	id, err := c.client.Post("categories", jsonStr)
	if err != nil {
		log.Print("POST /categories/ ", err)
		return err
	}

	category.ID = id

	return nil
}

// UpdateCategory updates a Toshl Category
func (c *Client) UpdateCategory(category *Category) error {
	jsonBytes, err := json.Marshal(category)
	if err != nil {
		log.Println("UpdateCategory: ", err)
		return err
	}

	jsonStr := string(jsonBytes)

	categoryResponse, err := c.client.Update(
		fmt.Sprintf("categories/%s", category.ID), jsonStr)
	if err != nil {
		log.Print("PUT /categories/ ", err)
		return err
	}

	err = json.Unmarshal([]byte(categoryResponse), category)
	if err != nil {
		log.Println("Cannot decode Category JSON")
		return err
	}

	return nil
}

// DeleteCategory deletes a Toshl Category
func (c *Client) DeleteCategory(category *Category) error {
	err := c.client.Delete(fmt.Sprintf("categories/%s", category.ID))
	if err != nil {
		log.Print("DELETE /categories/ ", err)
		return err
	}

	return nil
}

// MergeCategories merges two ore more Toshl categories into a single one
func (c *Client) MergeCategories(order *CategoriesMergeParams) error {
	jsonBytes, err := json.Marshal(order)
	if err != nil {
		log.Println("MergeCategories: ", err)
		return err
	}

	jsonStr := string(jsonBytes)

	_, err = c.client.Post("categories/merge", jsonStr)
	if err != nil {
		log.Print("POST /categories/merge ", err)
		return err
	}

	return nil
}

func (c *Client) Entries(params *EntryQueryParams) ([]Entry, error) {
	queryString := ""
	var err error

	if params != nil {
		queryString, err = params.getQueryString()
		if err != nil {
			log.Fatal(err)
		}
	}

	responses, err := c.client.GetMultiple("entries", queryString)
	if err != nil {
		log.Println("GET /entries/: ", err)
		return nil, err
	}

	var entries []Entry

	for _, response := range responses {
		var responseEntries []Entry
		err = json.Unmarshal([]byte(response), &responseEntries)
		if err != nil {
			return nil, err
		}

		entries = append(entries, responseEntries...)
	}

	return entries, nil
}

func (c *Client) CreateEntry(entry *Entry) error {
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		log.Println("CreateEntry: ", err)
		return err
	}

	jsonStr := string(jsonBytes)

	id, err := c.client.Post("entries", jsonStr)
	if err != nil {
		log.Println("POST /entries/ ", err)
		return err
	}

	entry.Id = &id

	return nil
}

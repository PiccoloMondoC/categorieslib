// sky-categories/pkg/clientlib/categoriesclient/client.go
package categoriesclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Client represents an HTTP client that can be used to send requests to the categories server.
type Client struct {
	BaseURL    string
	HttpClient *http.Client
	Token      string
	ApiKey     string
}

// Category represents the structure of a category.
type Category struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}

// CreateCategoryRequest represents the structure of a category for a create request.
type CreateCategoryRequest struct {
	Name string `json:"name,omitempty"`
}

// AssociateCategoryWithProjectRequest represents the request structure to associate a category with a project.
type AssociateCategoryWithProjectRequest struct {
	CategoryID uuid.UUID `json:"category_id"`
	ProjectID  uuid.UUID `json:"project_id"`
}

// AssociateCategoryWithSkillRequest represents the request structure to associate a category with a skill.
type AssociateCategoryWithSkillRequest struct {
	CategoryID uuid.UUID `json:"category_id"`
	SkillID    uuid.UUID `json:"skill_id"`
}

// DisassociateCategoryFromSkillRequest represents the request structure to disassociate a category from a skill.
type DisassociateCategoryFromSkillRequest struct {
	CategoryID uuid.UUID `json:"category_id"`
	SkillID    uuid.UUID `json:"skill_id"`
}

// GetCategoriesForSkillRequest represents the request structure to get categories for a specific skill.
type GetCategoriesForSkillRequest struct {
	SkillID uuid.UUID `json:"skill_id"`
}

// GetCategoriesForSkillResponse represents the response structure from the get categories request.
type GetCategoriesForSkillResponse struct {
	Categories []Category `json:"categories"`
}

// GetSkillIDsForCategoryRequest represents the request structure to get skill IDs for a specific category.
type GetSkillIDsForCategoryRequest struct {
	CategoryID uuid.UUID `json:"category_id"`
}

// GetSkillIDsForCategoryResponse represents the response structure from the get skill IDs request.
type GetSkillIDsForCategoryResponse struct {
	SkillIDs []uuid.UUID `json:"skill_ids"`
}

func NewClient(baseURL string, token string, apiKey string, httpClient ...*http.Client) *Client {
	var client *http.Client
	if len(httpClient) > 0 {
		client = httpClient[0]
	} else {
		client = &http.Client{
			Timeout: time.Second * 10,
		}
	}

	return &Client{
		BaseURL:    baseURL,
		HttpClient: client,
		Token:      token,
		ApiKey:     apiKey,
	}
}

// CreateCategory creates a new category using the categories microservice.
func (c *Client) CreateCategory(cat *CreateCategoryRequest, authToken string) (*Category, error) {
	// Marshal the Category struct into a JSON string.
	reqBody, err := json.Marshal(cat)
	if err != nil {
		return nil, fmt.Errorf("failed to encode category into JSON: %w", err)
	}

	// Build the request.
	url := fmt.Sprintf("%s/api/categories", c.BaseURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set the request headers.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authToken)
	req.Header.Set("X-API-Key", c.ApiKey)

	// Send the request.
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check the status code of the response.
	if resp.StatusCode != http.StatusCreated {
		// For simplicity, we just return an error here.
		// In a real-world application, you'd likely want to return a more detailed error message.
		return nil, fmt.Errorf("unexpected status code: got %v", resp.StatusCode)
	}

	// Decode the response body.
	var createdCat Category
	if err := json.NewDecoder(resp.Body).Decode(&createdCat); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &createdCat, nil
}

// GetCategory retrieves a category using the categories microservice.
func (c *Client) GetCategory(categoryID uuid.UUID, authToken string) (*Category, error) {
	// Build the request.
	url := fmt.Sprintf("%s/api/categories/%s", c.BaseURL, categoryID.String())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set the request headers.
	req.Header.Set("Authorization", authToken)
	req.Header.Set("X-API-Key", c.ApiKey)

	// Send the request.
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check the status code of the response.
	if resp.StatusCode != http.StatusOK {
		// For simplicity, we just return an error here.
		// In a real-world application, you'd likely want to return a more detailed error message.
		return nil, fmt.Errorf("unexpected status code: got %v", resp.StatusCode)
	}

	// Decode the response body.
	var retrievedCategory Category
	if err := json.NewDecoder(resp.Body).Decode(&retrievedCategory); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &retrievedCategory, nil
}

// AssociateCategoryWithProject associates a category with a project using the categories microservice.
func (c *Client) AssociateCategoryWithProject(categoryID, projectID uuid.UUID, authToken string) error {
	// Build the request body.
	requestBody := AssociateCategoryWithProjectRequest{
		CategoryID: categoryID,
		ProjectID:  projectID,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to encode request body: %w", err)
	}

	// Build the request.
	url := fmt.Sprintf("%s/api/projects/categories/associate", c.BaseURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set the request headers.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authToken)
	req.Header.Set("X-API-Key", c.ApiKey)

	// Send the request.
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check the status code of the response.
	if resp.StatusCode != http.StatusCreated {
		// For simplicity, we just return an error here.
		// In a real-world application, you'd likely want to return a more detailed error message.
		return fmt.Errorf("unexpected status code: got %v", resp.StatusCode)
	}

	return nil
}

// DisassociateCategoryFromProject disassociates a category from a project using the categories microservice.
func (c *Client) DisassociateCategoryFromProject(categoryID, projectID uuid.UUID, authToken string) error {
	// Build the request body.
	requestBody := AssociateCategoryWithProjectRequest{
		CategoryID: categoryID,
		ProjectID:  projectID,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to encode request body: %w", err)
	}

	// Build the request.
	url := fmt.Sprintf("%s/api/projects/categories/disassociate", c.BaseURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set the request headers.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authToken)
	req.Header.Set("X-API-Key", c.ApiKey)

	// Send the request.
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check the status code of the response.
	if resp.StatusCode != http.StatusNoContent {
		// For simplicity, we just return an error here.
		// In a real-world application, you'd likely want to return a more detailed error message.
		return fmt.Errorf("unexpected status code: got %v", resp.StatusCode)
	}

	return nil
}

// GetCategoriesForProject gets the categories associated with a specific project using the categories microservice.
func (c *Client) GetCategoriesForProject(projectID uuid.UUID, authToken string) ([]Category, error) {
	// Build the request.
	url := fmt.Sprintf("%s/api/projects/%s/categories", c.BaseURL, projectID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set the request headers.
	req.Header.Set("Authorization", authToken)
	req.Header.Set("X-API-Key", c.ApiKey)

	// Send the request.
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check the status code of the response.
	if resp.StatusCode != http.StatusOK {
		// For simplicity, we just return an error here.
		// In a real-world application, you'd likely want to return a more detailed error message.
		return nil, fmt.Errorf("unexpected status code: got %v", resp.StatusCode)
	}

	// Decode the response body.
	var categories []Category
	err = json.NewDecoder(resp.Body).Decode(&categories)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return categories, nil
}

// GetProjectIDsForCategory gets the project IDs associated with a specific category using the categories microservice.
func (c *Client) GetProjectIDsForCategory(categoryID uuid.UUID, authToken string) ([]uuid.UUID, error) {
	// Build the request.
	url := fmt.Sprintf("%s/api/categories/%s/projects", c.BaseURL, categoryID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set the request headers.
	req.Header.Set("Authorization", authToken)
	req.Header.Set("X-API-Key", c.ApiKey)

	// Send the request.
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check the status code of the response.
	if resp.StatusCode != http.StatusOK {
		// For simplicity, we just return an error here.
		// In a real-world application, you'd likely want to return a more detailed error message.
		return nil, fmt.Errorf("unexpected status code: got %v", resp.StatusCode)
	}

	// Decode the response body.
	var projectIDs []uuid.UUID
	err = json.NewDecoder(resp.Body).Decode(&projectIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return projectIDs, nil
}

// AssociateCategoryWithSkill associates a category with a skill using the categories microservice.
func (c *Client) AssociateCategoryWithSkill(request AssociateCategoryWithSkillRequest, authToken string) error {
	// Marshal the request into JSON.
	reqBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Build the request.
	url := fmt.Sprintf("%s/api/categories/skills/association", c.BaseURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set the request headers.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authToken)
	req.Header.Set("X-API-Key", c.ApiKey)

	// Send the request.
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check the status code of the response.
	if resp.StatusCode != http.StatusCreated {
		// For simplicity, we just return an error here.
		// In a real-world application, you'd likely want to return a more detailed error message.
		return fmt.Errorf("unexpected status code: got %v", resp.StatusCode)
	}

	return nil
}

// DisassociateCategoryFromSkill disassociates a category from a skill using the categories microservice.
func (c *Client) DisassociateCategoryFromSkill(request DisassociateCategoryFromSkillRequest, authToken string) error {
	// Marshal the request into JSON.
	reqBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Build the request.
	url := fmt.Sprintf("%s/api/categories/skills/disassociation", c.BaseURL)
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set the request headers.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authToken)
	req.Header.Set("X-API-Key", c.ApiKey)

	// Send the request.
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check the status code of the response.
	if resp.StatusCode != http.StatusNoContent {
		// For simplicity, we just return an error here.
		// In a real-world application, you'd likely want to return a more detailed error message.
		return fmt.Errorf("unexpected status code: got %v", resp.StatusCode)
	}

	return nil
}

// GetCategoriesForSkill retrieves the categories associated with a specific skill.
func (c *Client) GetCategoriesForSkill(request GetCategoriesForSkillRequest, authToken string) (GetCategoriesForSkillResponse, error) {
	// Prepare the URL for the request.
	url := fmt.Sprintf("%s/api/skills/%s/categories", c.BaseURL, request.SkillID)

	// Create the HTTP request.
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return GetCategoriesForSkillResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	// Set the request headers.
	req.Header.Set("Authorization", authToken)
	req.Header.Set("X-API-Key", c.ApiKey)

	// Send the request.
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return GetCategoriesForSkillResponse{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check the status code of the response.
	if resp.StatusCode != http.StatusOK {
		// For simplicity, we just return an error here.
		// In a real-world application, you'd likely want to return a more detailed error message.
		return GetCategoriesForSkillResponse{}, fmt.Errorf("unexpected status code: got %v", resp.StatusCode)
	}

	// Decode the response.
	var response GetCategoriesForSkillResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return GetCategoriesForSkillResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return response, nil
}

// GetSkillIDsForCategory retrieves the skill IDs associated with a specific category.
func (c *Client) GetSkillIDsForCategory(request GetSkillIDsForCategoryRequest, authToken string) (GetSkillIDsForCategoryResponse, error) {
	// Prepare the URL for the request.
	url := fmt.Sprintf("%s/api/categories/%s/skills", c.BaseURL, request.CategoryID)

	// Create the HTTP request.
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return GetSkillIDsForCategoryResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	// Set the request headers.
	req.Header.Set("Authorization", authToken)
	req.Header.Set("X-API-Key", c.ApiKey)

	// Send the request.
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return GetSkillIDsForCategoryResponse{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check the status code of the response.
	if resp.StatusCode != http.StatusOK {
		// For simplicity, we just return an error here.
		// In a real-world application, you'd likely want to return a more detailed error message.
		return GetSkillIDsForCategoryResponse{}, fmt.Errorf("unexpected status code: got %v", resp.StatusCode)
	}

	// Decode the response.
	var response GetSkillIDsForCategoryResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return GetSkillIDsForCategoryResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return response, nil
}

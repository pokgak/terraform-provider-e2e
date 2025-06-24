package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	
	
	
	"net/http"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

// CreateMariaDB provisions a new MariaDB instance on E2E Cloud.
//
// It sends a POST request to the DBaaS API with the provided configuration,
// including software/template IDs, VPC attachments, DB config, and optional encryption.
//
// Parameters:
//   - req: Pointer to MariaDBCreateRequest containing user-defined configuration
//   - projectID: Project under which the DB is to be created
//   - location: Cloud region for provisioning
//
// Returns:
//   - Pointer to created MariaDB instance (parsed from API response)
//   - Error if request fails or response is invalid
func (c *Client) CreateMariaDB(req *models.MariaDBCreateRequest, projectID, location string) (*models.MariaDB, error) {
	url := strings.TrimRight(c.Api_endpoint, "/") + "/rds/cluster/"

	// === Step 1: Encode the payload ===
	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(req); err != nil {
		return nil, fmt.Errorf("failed to encode create payload: %v", err)
	}

	// === Step 2: Construct HTTP POST request ===
	httpReq, err := http.NewRequest("POST", url, payloadBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Optional: Debug log of request payload (comment or guard this in production)
	log.Printf("[DEBUG] MariaDB create request payload: %+v", req)

	// === Step 3: Set headers and query params ===
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.Auth_token)
	httpReq.Header.Set("x-api-key", c.Api_key)

	q := httpReq.URL.Query()
	q.Add("project_id", projectID)
	q.Add("location", location)
	httpReq.URL.RawQuery = q.Encode()

	// === Step 4: Execute request ===
	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	// === Step 5: Handle response status ===
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// === Step 6: Decode successful response ===
	var response models.MariaDBResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response.Data, nil
}






// ReadMariaDB retrieves the current state of a MariaDB instance by its ID.
//
// It performs a GET request to the E2E Cloud API and returns the deserialized
// MariaDB object if found.
//
// Parameters:
//   - id: Unique identifier of the MariaDB instance
//   - projectID: Project under which the DB is provisioned
//   - location: Region of the instance
//
// Returns:
//   - Pointer to MariaDB object populated with backend state
//   - Error if the request fails or response is malformed
func (c *Client) ReadMariaDB(id string, projectID string, location string) (*models.MariaDB, error) {
	url := c.Api_endpoint + "/rds/cluster/" + id + "/"

	// === Step 1: Build GET request ===
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// === Step 2: Add query params and headers ===
	q := req.URL.Query()
	q.Add("apikey", c.Api_key)
	q.Add("project_id", projectID)
	q.Add("location", location)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Auth_token)
	req.Header.Set("x-api-key", c.Api_key)

	// Debug logs for troubleshooting (safe to remove or guard in production)
	log.Printf("[DEBUG] ReadMariaDB Request URL: %s", req.URL.String())
	log.Printf("[DEBUG] ReadMariaDB Headers: %+v", req.Header)

	// === Step 3: Execute GET request ===
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	// === Step 4: Handle failure status codes ===
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("read mariadb failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// === Step 5: Decode JSON response ===
	var response models.MariaDBResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response.Data, nil
}



// MariaDBExists checks whether a MariaDB cluster exists on the E2E Cloud platform
// by performing a GET request to the cluster endpoint.
// Returns true if the cluster exists, false if it does not (404), and an error otherwise.
func (c *Client) MariaDBExists(id string, projectID string, location string) (bool, error) {
	// Construct the API endpoint URL for the cluster
	url := strings.TrimRight(c.Api_endpoint, "/") + "/rds/cluster/" + id + "/"

	// Create HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %v", err)
	}

	// Add required query parameters
	q := req.URL.Query()
	q.Add("apikey", c.Api_key)
	q.Add("project_id", projectID)
	q.Add("location", location)
	req.URL.RawQuery = q.Encode()

	// Set necessary headers for authentication and content type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Auth_token)
	req.Header.Set("x-api-key", c.Api_key)

	// Execute the HTTP request
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	// Return false if resource is not found (404)
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	// Return error if any unexpected status code is received
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Resource exists
	return true, nil
}





// DeleteMariaDB deletes a MariaDB cluster by its ID.
func (c *Client) DeleteMariaDB(id string, projectID string, location string) error {
	url := strings.TrimRight(c.Api_endpoint, "/") + "/rds/cluster/" + id + "/"

	// Create DELETE request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %v", err)
	}

	// Set query parameters
	q := req.URL.Query()
	q.Add("apikey", c.Api_key)
	q.Add("project_id", projectID)
	q.Add("location", location)
	req.URL.RawQuery = q.Encode()

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Auth_token)
	req.Header.Set("x-api-key", c.Api_key)

	// Send request
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("DELETE request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check for successful status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("DELETE failed with status: %d", resp.StatusCode)
	}

	return nil
}



// ShutdownMariaDB sends a request to shut down a MariaDB cluster by ID.
// It issues a PUT request to the /rds/cluster/{id}/shutdown endpoint using the client's credentials.
// Returns an error if the request fails or if the response is not a 200 OK.
func (c *Client) ShutdownMariaDB(id string, projectID string, location string) error {
	// Construct the full shutdown endpoint URL
	url := fmt.Sprintf("%s/rds/cluster/%s/shutdown", strings.TrimRight(c.Api_endpoint, "/"), id)

	// Create a new HTTP PUT request to the shutdown endpoint
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create shutdown request: %v", err)
	}

	// Set required query parameters: API key, project ID, and location
	q := req.URL.Query()
	q.Add("apikey", c.Api_key)
	q.Add("project_id", projectID)
	q.Add("location", location)
	req.URL.RawQuery = q.Encode()

	// Set required HTTP headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Auth_token)
	req.Header.Set("x-api-key", c.Api_key)

	// Perform the HTTP request using the client’s HTTP client
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("shutdown request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check if the shutdown request was successful (HTTP 200)
	if resp.StatusCode != http.StatusOK {
		// Read response body for additional error context
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("shutdown failed: %d - %s", resp.StatusCode, string(body))
	}

	// Return nil to indicate successful shutdown
	return nil
}


// ResumeMariaDB sends a PUT request to resume (start) a suspended MariaDB instance.
// It uses the database cluster ID, project ID, and location for the request.
func (c *Client) ResumeMariaDB(id string, projectID string, location string) error {
	url := fmt.Sprintf("%s/rds/cluster/%s/resume", strings.TrimRight(c.Api_endpoint, "/"), id)

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create resume request: %v", err)
	}

	// Add required query parameters
	q := req.URL.Query()
	q.Add("apikey", c.Api_key)
	q.Add("project_id", projectID)
	q.Add("location", location)
	req.URL.RawQuery = q.Encode()

	// Set required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Auth_token)
	req.Header.Set("x-api-key", c.Api_key)

	// Execute request
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("resume request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("resume failed: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}



// RestartMariaDB triggers a restart for a specific MariaDB cluster.
// It performs a PUT request to the /restart endpoint with required query params.
func (c *Client) RestartMariaDB(id string, projectID string, location string) error {
	url := fmt.Sprintf("%s/rds/cluster/%s/restart", strings.TrimRight(c.Api_endpoint, "/"), id)

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create restart request: %v", err)
	}

	// Add required query parameters
	q := req.URL.Query()
	q.Add("apikey", c.Api_key)
	q.Add("project_id", projectID)
	q.Add("location", location)
	req.URL.RawQuery = q.Encode()

	// Set necessary headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Auth_token)
	req.Header.Set("x-api-key", c.Api_key)

	// Execute the HTTP request
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("restart request failed: %v", err)
	}
	defer resp.Body.Close()

	// Expect HTTP 200 on success
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("restart failed: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}



// AttachVPCToMariaDB attaches one or more VPCs to the MariaDB instance by resolving their full metadata.
func (c *Client) AttachVPCToMariaDB(id string, projectID string, location string, vpcIDs []string) error {
	// Expand metadata for all VPC IDs
	vpcMetaList, err := c.ExpandVpcList(vpcIDs, projectID, location)
	if err != nil {
		return fmt.Errorf("failed to expand VPC metadata: %v", err)
	}

	// Build request payload
	payload := models.AttachDetachVPCRequest{
		Action: "attach",
		VPCs:   vpcMetaList,
	}

	// Prepare URL
	url := fmt.Sprintf("%s/rds/cluster/%s/vpc-attach/?apikey=%s&project_id=%s&location=%s",
		strings.TrimRight(c.Api_endpoint, "/"), id, c.Api_key, projectID, location,
	)

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal VPC attach payload: %v", err)
	}

	// Make HTTP PUT request
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Auth_token)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("VPC attach request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("VPC attach failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	return nil
}


// DetachVPCFromMariaDB detaches one or more VPCs from the MariaDB instance.
func (c *Client) DetachVPCFromMariaDB(id string, projectID string, location string, vpcIDs []string) error {
	// Expand metadata for all VPC IDs
	vpcMetaList, err := c.ExpandVpcList(vpcIDs, projectID, location)
	if err != nil {
		return fmt.Errorf("failed to expand VPC metadata: %v", err)
	}

	// Build request payload
	payload := models.AttachDetachVPCRequest{
		Action: "detach",
		VPCs:   vpcMetaList,
	}

	// Prepare URL
	url := fmt.Sprintf("%s/rds/cluster/%s/vpc-detach/?apikey=%s&project_id=%s&location=%s",
		strings.TrimRight(c.Api_endpoint, "/"), id, c.Api_key, projectID, location,
	)

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal VPC detach payload: %v", err)
	}

	// Create HTTP PUT request
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Auth_token)

	// Send request
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("VPC detach request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("VPC detach failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}




// AttachPublicIPToMariaDB attaches a public IP to the specified MariaDB cluster instance.
// This triggers a PUT request to the E2E Cloud API with the "attach" action.
// It is used when the user sets `public_ip_enabled = true` in Terraform.
//
// Parameters:
// - id: Unique identifier of the MariaDB cluster
// - projectID: E2E project ID in which the service resides
// - location: Region/location where the cluster is deployed
//
// Returns an error if the request fails or the API responds with a non-200 status code.
func (c *Client) AttachPublicIPToMariaDB(id string, projectID string, location string) error {
	url := fmt.Sprintf("%s/rds/cluster/%s/public-ip-attach/?apikey=%s&project_id=%s&location=%s",
		strings.TrimRight(c.Api_endpoint, "/"), id, c.Api_key, projectID, location)

	// Prepare JSON payload for the attach action
	payload := map[string]string{"action": "attach"}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal public IP attach payload: %v", err)
	}

	// Construct PUT request
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create attach public IP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Auth_token)

	// Execute request
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("public IP attach request failed: %v", err)
	}
	defer resp.Body.Close()

	// Validate response status
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("public IP attach failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}





// DetachPublicIPFromMariaDB detaches the public IP from a given MariaDB cluster instance.
// It sends a PUT request to the E2E Cloud API with the "detach" action.
// This is triggered when the user sets `public_ip_enabled = false` in Terraform.
//
// Parameters:
// - id: Unique identifier of the MariaDB cluster
// - projectID: E2E project ID in which the service resides
// - location: Region/location where the cluster is deployed
//
// Returns an error if the request fails or the API responds with a non-200 status code.
func (c *Client) DetachPublicIPFromMariaDB(id string, projectID string, location string) error {
	url := fmt.Sprintf("%s/rds/cluster/%s/public-ip-detach/?apikey=%s&project_id=%s&location=%s",
		strings.TrimRight(c.Api_endpoint, "/"), id, c.Api_key, projectID, location)

	// Prepare JSON payload for the detach action
	payload := map[string]string{"action": "detach"}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal public IP detach payload: %v", err)
	}

	// Construct PUT request
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create detach public IP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Auth_token)

	// Execute request
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("public IP detach request failed: %v", err)
	}
	defer resp.Body.Close()

	// Validate response status
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("public IP detach failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}



// AttachParameterGroupToMariaDB attaches a parameter group to a specified MariaDB cluster.
// This method sends a PUT request to the E2E Cloud API to attach a parameter group
// using the action "add". It is used to modify DB runtime configuration by applying a
// predefined parameter group.
//
// Parameters:
//   - clusterID: ID of the target MariaDB cluster.
//   - parameterGroupID: ID of the parameter group to attach.
//   - projectID: ID of the project under which the cluster is provisioned.
//   - location: Geographic location (region) where the cluster resides.
//
// Returns:
//   - error: if the request fails to execute or receives a non-200 HTTP response.
func (c *Client) AttachParameterGroupToMariaDB(clusterID string, parameterGroupID int, projectID, location string) error {
	// Construct the API endpoint URL with required query parameters.
	url := fmt.Sprintf("%s/rds/cluster/%s/parameter-group/%d/add?apikey=%s&project_id=%s&location=%s",
		strings.TrimRight(c.Api_endpoint, "/"), clusterID, parameterGroupID, c.Api_key, projectID, location)

	// Prepare the payload to instruct the API to perform an "add" operation.
	payload := models.ParameterGroupRequest{Action: "add"}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal attach payload: %v", err)
	}

	// Create a new PUT request with the marshaled payload.
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create attach request: %v", err)
	}

	// Set authentication and content headers.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Auth_token)

	// Send the request and handle network errors.
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("attach request failed: %v", err)
	}
	defer resp.Body.Close()

	// Handle non-successful HTTP response.
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("attach failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}



// DetachParameterGroupFromMariaDB removes an attached parameter group from a MariaDB cluster.
// This method sends a PUT request to the E2E Cloud API to detach a parameter group
// using the "detach" endpoint. This is typically done to revert to default DB parameters.
//
// Parameters:
//   - clusterID: ID of the target MariaDB cluster.
//   - parameterGroupID: ID of the parameter group to detach.
//   - projectID: ID of the project under which the cluster is provisioned.
//   - location: Geographic location (region) where the cluster resides.
//
// Returns:
//   - error: if the request fails or the API returns a non-200 HTTP status.
func (c *Client) DetachParameterGroupFromMariaDB(clusterID string, parameterGroupID int, projectID, location string) error {
	// Construct the API endpoint URL for detachment.
	url := fmt.Sprintf("%s/rds/cluster/%s/parameter-group/%d/detach", c.Api_endpoint, clusterID, parameterGroupID)

	// Prepare a minimal empty JSON body as the API expects a payload.
	body := bytes.NewBuffer([]byte("{}"))

	// Create the HTTP PUT request.
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return fmt.Errorf("failed to create detach request: %v", err)
	}

	// Append required query parameters.
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", projectID)
	params.Add("location", location)
	req.URL.RawQuery = params.Encode()

	// Add necessary headers.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Auth_token)

	// Execute the request.
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("detach request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read and handle the response for non-200 cases.
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("detach failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}


// UpgradeMariaDBPlan performs a plan upgrade on the given MariaDB cluster by sending
// a PUT request to the /rds-upgrade/ endpoint with the new template ID.
// This changes the hardware configuration of the cluster (e.g., CPU, RAM).
func (c *Client) UpgradeMariaDBPlan(clusterID, projectID, location string, templateID int) error {
	// Construct the full API URL for plan upgrade.
	url := fmt.Sprintf("%s/rds/cluster/%s/rds-upgrade/", strings.TrimRight(c.Api_endpoint, "/"), clusterID)

	// Define request payload with the new template ID.
	payload := map[string]interface{}{
		"template_id": templateID,
	}

	// Marshal the payload to JSON.
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal upgrade payload: %v", err)
	}

	// Create HTTP PUT request with JSON body.
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create upgrade request: %v", err)
	}

	// Add required query parameters.
	q := req.URL.Query()
	q.Add("apikey", c.Api_key)
	q.Add("project_id", projectID)
	q.Add("location", location)
	req.URL.RawQuery = q.Encode()

	// Set required headers for authentication and content type.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Auth_token)
	req.Header.Set("x-api-key", c.Api_key)

	// Send the request and handle any network error.
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("upgrade request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check for a successful status code (200 OK).
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upgrade failed: %d - %s", resp.StatusCode, string(respBody))
	}

	return nil
}


// ExpandMariaDBDisk performs a disk size upgrade for a given MariaDB cluster.
// It adds the specified `additionalSize` (in GB) to the current disk size.
// This operation requires the cluster to be in STOPPED state prior to the call.
func (c *Client) ExpandMariaDBDisk(clusterID, projectID, location string, additionalSize int) error {
	// Skip disk expansion if additionalSize is zero or negative.
	if additionalSize <= 0 {
		log.Printf("[INFO] Skipping disk expansion: additional size is %d GB (must be > 0)", additionalSize)
		return nil
	}

	log.Printf("[INFO] Initiating disk expansion: cluster=%s, additional_size=%d GB", clusterID, additionalSize)

	// Build the disk upgrade endpoint URL.
	url := fmt.Sprintf("%s/rds/cluster/%s/disk-upgrade/", strings.TrimRight(c.Api_endpoint, "/"), clusterID)

	// Prepare payload for disk expansion.
	payload := models.DiskUpgradeRequest{Size: additionalSize}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal disk upgrade payload: %v", err)
	}

	// Create HTTP PUT request.
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create disk upgrade request: %v", err)
	}

	// Add required headers.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Auth_token)
	req.Header.Set("x-api-key", c.Api_key)

	// Add required query parameters.
	query := req.URL.Query()
	query.Add("apikey", c.Api_key)
	query.Add("project_id", projectID)
	query.Add("location", location)
	req.URL.RawQuery = query.Encode()

	// Send the request.
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("disk upgrade request failed: %v", err)
	}
	defer resp.Body.Close()

	// Handle error if upgrade fails.
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("disk upgrade failed with status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("[INFO] Disk expansion completed: +%d GB added to MariaDB cluster %s", additionalSize, clusterID)
	return nil
}











package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	

	"net/http"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func (c *Client) CreateScalerGroup(req *models.CreateScalerGroupRequest, projectID, location string) (*models.ScalerGroupCreateDetails, error) {
	url := c.Api_endpoint + "/scaler/scalegroups"
	log.Printf("[INFO] Sending request to create Scaler Group at: %s", url)

	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(req); err != nil {
		log.Printf("[ERROR] Failed to encode create payload: %v", err)
		return nil, fmt.Errorf("failed to encode create payload: %v", err)
	}

	httpReq, err := http.NewRequest("POST", url, payloadBuf)
	if err != nil {
		log.Printf("[ERROR] Failed to create HTTP request: %v", err)
		return nil, fmt.Errorf("failed to create POST request: %v", err)
	}
	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)

	log.Printf("[DEBUG] CreateScalerGroup headers: %v", httpReq.Header)

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		log.Printf("[ERROR] HTTP request failed: %v", err)
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("[ERROR] CreateScalerGroup failed: status %d", resp.StatusCode)
		log.Printf("[ERROR] Response body: %s", string(bodyBytes))
		return nil, fmt.Errorf("create scaler group failed: status %d\nresponse: %s", resp.StatusCode, string(bodyBytes))
	}

	var response models.CreateScalerGroupResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		log.Printf("[ERROR] Failed to decode response JSON: %v", err)
		log.Printf("[DEBUG] Raw body: %s", string(bodyBytes))
		return nil, fmt.Errorf("failed to decode create response: %v\nresponse body: %s", err, string(bodyBytes))
	}

	log.Printf("[INFO] Scaler Group created successfully: ID=%s Name=%s", response.Data.ID, response.Data.Name)

	return &response.Data, nil
}

func (c *Client) GetScalerGroup(scaleGroupID, projectID, location string) (*models.ScalerGroupGetDetail, error) {
	url := c.Api_endpoint + "/scaler/scalegroups/" + scaleGroupID + "/"
	log.Printf("[INFO] Fetching Scaler Group details for ID: %s", scaleGroupID)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to create GET request: %v", err)
		return nil, fmt.Errorf("failed to create GET request: %v", err)
	}
	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)

	log.Printf("[DEBUG] GetScalerGroup request URL: %s", httpReq.URL.String())

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		log.Printf("[ERROR] HTTP request failed: %v", err)
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] GetScalerGroup failed: status %d", resp.StatusCode)
		log.Printf("[DEBUG] Response body: %s", string(bodyBytes))
		return nil, fmt.Errorf("get scaler group failed: status %d\nresponse: %s", resp.StatusCode, string(bodyBytes))
	}

	var response models.GetScalerGroupResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		log.Printf("[ERROR] Failed to decode get response: %v", err)
		log.Printf("[DEBUG] Raw body: %s", string(bodyBytes))
		return nil, fmt.Errorf("failed to decode get response: %v\nresponse body: %s", err, string(bodyBytes))
	}

	log.Printf("[INFO] Scaler Group fetched: Name=%s Desired=%d Running=%d", response.Data.Name, response.Data.Desired, response.Data.Running)

	return &response.Data, nil
}

func (c *Client) DeleteScalerGroup(scaleGroupID, projectID, location string) error {
	url := c.Api_endpoint + "/scaler/scalegroups/" + scaleGroupID + "/"
	log.Printf("[INFO] Sending delete request for Scaler Group ID: %s", scaleGroupID)

	httpReq, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to create DELETE request: %v", err)
		return fmt.Errorf("failed to create DELETE request: %v", err)
	}
	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)

	log.Printf("[DEBUG] DeleteScalerGroup request URL: %s", httpReq.URL.String())

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		log.Printf("[ERROR] HTTP request failed: %v", err)
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read delete response body: %v", err)
		return fmt.Errorf("failed to read delete response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		log.Printf("[ERROR] DeleteScalerGroup failed: status %d", resp.StatusCode)
		log.Printf("[DEBUG] Response body: %s", string(bodyBytes))
		return fmt.Errorf("delete scaler group failed: status %d\nresponse: %s", resp.StatusCode, string(bodyBytes))
	}

	var response models.DeleteScalerGroupResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		log.Printf("[ERROR] Failed to decode delete response: %v", err)
		log.Printf("[DEBUG] Raw body: %s", string(bodyBytes))
		return fmt.Errorf("failed to decode delete response: %v", err)
	}

	log.Printf("[INFO] Scaler Group deleted successfully. Message: %s", response.Message)
	return nil
}

func (c *Client) GetSavedImageByName(imageName, projectID, location string) (*models.SavedImage, error) {
	url := c.Api_endpoint + "/images/saved-images/"
	log.Printf("[INFO] Sending request to fetch saved image: %s", imageName)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to create GET request: %v", err)
		return nil, fmt.Errorf("failed to create GET request: %v", err)
	}

	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)

	log.Printf("[DEBUG] GetSavedImageByName request URL: %s", httpReq.URL.String())

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		log.Printf("[ERROR] HTTP request failed: %v", err)
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] GetSavedImageByName failed: status %d", resp.StatusCode)
		log.Printf("[ERROR] Response body: %s", string(bodyBytes))
		return nil, fmt.Errorf("get saved image failed: status %d\nresponse: %s", resp.StatusCode, string(bodyBytes))
	}

	var result models.ListSavedImagesResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		log.Printf("[ERROR] Failed to decode saved-images response: %v", err)
		log.Printf("[DEBUG] Raw body: %s", string(bodyBytes))
		return nil, fmt.Errorf("failed to decode saved-images response: %v\nresponse body: %s", err, string(bodyBytes))
	}

	for _, img := range result.Data {
		if img.Name == imageName {
			log.Printf("[INFO] Found saved image: ID=%s, TemplateID=%d, Distro=%s", img.ImageID, img.TemplateID, img.Distro)
			return &img, nil
		}
	}

	log.Printf("[ERROR] No saved image found with name: %s", imageName)
	return nil, fmt.Errorf("no saved image found with name: %s", imageName)
}
func (c *Client) GetDefaultSecurityGroupID(projectID, location string) (int, error) {
	url := c.Api_endpoint + "security_group/"
	log.Printf("[INFO] Sending request to fetch default Security Group at: %s", url)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to create GET request: %v", err)
		return 0, fmt.Errorf("failed to create GET request: %v", err)
	}

	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)
	log.Printf("[DEBUG] GetDefaultSecurityGroup headers: %v", httpReq.Header)

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		log.Printf("[ERROR] HTTP request failed: %v", err)
		return 0, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read response body: %v", err)
		return 0, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] GetDefaultSecurityGroupID failed: status %d", resp.StatusCode)
		log.Printf("[ERROR] Response body: %s", string(bodyBytes))
		return 0, fmt.Errorf("get default security group failed: status %d\nresponse: %s", resp.StatusCode, string(bodyBytes))
	}

	var response models.GetScalerSecurityGroupsResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		log.Printf("[ERROR] Failed to decode response JSON: %v", err)
		log.Printf("[DEBUG] Raw body: %s", string(bodyBytes))
		return 0, fmt.Errorf("failed to decode security group response: %v\nresponse body: %s", err, string(bodyBytes))
	}

	for _, sg := range response.Data {
		if sg.IsDefault {
			log.Printf("[INFO] Found default Security Group: ID=%d", sg.ID)
			return sg.ID, nil
		}
	}

	log.Printf("[ERROR] No default security group found in response")
	return 0, fmt.Errorf("default security group not found")
}

func (c *Client) GetPlanDetailsFromPlanName(templateID int, planName, projectID, location string) (string, string, error) {
	url := c.Api_endpoint + fmt.Sprintf("/images/upgradeimage/%d/", templateID)
	log.Printf("[INFO] Sending request to fetch plan details for planName=%s, templateID=%d", planName, templateID)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to create GET request: %v", err)
		return "", "", fmt.Errorf("failed to create GET request: %v", err)
	}

	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)
	log.Printf("[DEBUG] GetPlanDetailsFromPlanName request URL: %s", httpReq.URL.String())

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		log.Printf("[ERROR] HTTP request failed: %v", err)
		return "", "", fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read response body: %v", err)
		return "", "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] GetPlanDetailsFromPlanName failed: status %d", resp.StatusCode)
		log.Printf("[ERROR] Response body: %s", string(bodyBytes))
		return "", "", fmt.Errorf("get plan details failed: status %d\nresponse: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Code int `json:"code"`
		Data []struct {
			Name  string `json:"name"` // UI plan name, e.g., "C3.8GB"
			Plan  string `json:"plan"` // slug name
			Specs struct {
				ID string `json:"id"` // plan_id / sku_id
			} `json:"specs"`
		} `json:"data"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		log.Printf("[ERROR] Failed to decode plan details response: %v", err)
		log.Printf("[DEBUG] Raw body: %s", string(bodyBytes))
		return "", "", fmt.Errorf("failed to decode plan details response: %v\nresponse body: %s", err, string(bodyBytes))
	}

	for _, item := range result.Data {
		if item.Name == planName {
			log.Printf("[INFO] Found plan: PlanID=%s, SlugName=%s", item.Specs.ID, item.Plan)
			return item.Specs.ID, item.Plan, nil
		}
	}

	log.Printf("[ERROR] No matching plan found for planName: %s", planName)
	return "", "", fmt.Errorf("plan name %s not found in template %d", planName, templateID)
}




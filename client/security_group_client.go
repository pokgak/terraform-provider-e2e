package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

// GetSecurityGroups fetches the list of security groups
func (c *Client) GetSecurityGroups(projectID int, location string) (*models.SecurityGroupListResponse, error) {
	url := fmt.Sprintf("%ssecurity_group/", c.Api_endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var response models.SecurityGroupListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// CreateSecurityGroup creates a new security group
func (c *Client) CreateSecurityGroup(request *models.SecurityGroupCreateRequest, projectID int, location string) (*models.BaseResponse, error) {
	url := fmt.Sprintf("%ssecurity_group/", c.Api_endpoint)
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var response models.BaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetSecurityGroup fetches a specific security group by ID
func (c *Client) GetSecurityGroup(securityGroupID float64, projectID int, location string) (*models.SecurityGroup, error) {
	groups, err := c.GetSecurityGroups(projectID, location)
	if err != nil {
		return nil, err
	}

	for _, group := range groups.Data {
		if group.Id == securityGroupID {
			return &group, nil
		}
	}
	return nil, fmt.Errorf("security group with ID %.0f not found", securityGroupID)
}

// UpdateSecurityGroup updates the security group details including rules
func (c *Client) UpdateSecurityGroup(securityGroupID float64, request *models.SecurityGroupCreateRequest, projectID int, location string) error {
	url := fmt.Sprintf("%ssecurity_group/%.0f/", c.Api_endpoint, securityGroupID)
	payload, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// RenameSecurityGroup renames the security group
func (c *Client) RenameSecurityGroup(securityGroupID float64, newName string, projectID int, location string) error {
	url := fmt.Sprintf("%ssecurity_group/%.0f/actions/", c.Api_endpoint, securityGroupID)
	payload := map[string]string{"type": "rename", "data": newName}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// UpdateSecurityGroupDescription updates the description
func (c *Client) UpdateSecurityGroupDescription(securityGroupID float64, description string, projectID int, location string) error {
	url := fmt.Sprintf("%ssecurity_group/%.0f/actions/", c.Api_endpoint, securityGroupID)
	payload := map[string]string{"type": "description", "data": description}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// MarkSecurityGroupAsDefault marks the security group as default
func (c *Client) MarkSecurityGroupAsDefault(securityGroupID float64, projectID int, location string) error {
	url := fmt.Sprintf("%ssecurity_group/%.0f/mark-default/", c.Api_endpoint, securityGroupID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// DeleteSecurityGroup deletes the security group
func (c *Client) DeleteSecurityGroup(securityGroupID float64, projectID int, location string) (*models.BaseResponse, error) {
	url := fmt.Sprintf("%ssecurity_group/%.0f/", c.Api_endpoint, securityGroupID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var response models.BaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

func (c *Client) GetSecurityGroupAssociatedNodes(security_group_id float64, project_id int, location string) (*models.SecurityGroupAssociatedNodesResponse, error) {
	url := fmt.Sprintf("%ssecurity_group/%.0f/vms/", c.Api_endpoint, security_group_id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("[INFO] error getting associated nodes for security group %.0f", security_group_id)
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	res := models.SecurityGroupAssociatedNodesResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("[INFO] error unmarshalling associated nodes response")
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetSecurityGroupAssociatedScalars(security_group_id float64, project_id int, location string) (*models.SecurityGroupAssociatedScalarsResponse, error) {
	url := fmt.Sprintf("%ssecurity_group/%.0f/associated-scalars/", c.Api_endpoint, security_group_id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("[INFO] error getting associated scalars for security group %.0f", security_group_id)
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	res := models.SecurityGroupAssociatedScalarsResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("[INFO] error unmarshalling associated scalars response")
		return nil, err
	}

	return &res, nil
}

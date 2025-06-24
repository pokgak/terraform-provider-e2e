package client

import (
	"encoding/json"
	"errors"
	"fmt"
	
	
	"io"
	"log"
	"net/http"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	
)

// GetSoftwareId fetches the software ID for a given engine name and version.
// It queries the /rds/plans/ endpoint and searches for a matching engine.
// Returns the software ID if found, otherwise returns an error.
func (c *Client) GetSoftwareId(projectID string, location string, name string, version string) (int, error) {
	url := c.Api_endpoint + "rds/plans/"

	// Build GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, err
	}

	// Add query parameters
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", projectID)
	params.Add("location", location)
	req.URL.RawQuery = params.Encode()

	// Add authentication headers
	SetBasicHeaders(c.Auth_token, req)

	// Send request
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("[ERROR] GetSoftwareId request failed: %v", err)
		return -1, err
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read GetSoftwareId response body: %v", err)
		return -1, err
	}

	// Parse JSON response into PlanResponse model
	var res models.PlanResponse
	if err := json.Unmarshal(body, &res); err != nil {
		log.Printf("[ERROR] Failed to unmarshal GetSoftwareId response: %v", err)
		return -1, err
	}

	// Search for a matching engine (name + version)
	for _, item := range res.Data.DatabaseEngines {
		if item.EngineName == name && item.EngineVersion == version {
			return item.EngineID, nil
		}
	}

	log.Printf("[INFO] Software ID not found for name: %s, version: %s", name, version)
	return -1, errors.New("matching engine not found")
}

// GetTemplateId fetches the template ID for a given plan name and software ID.
// It queries the /rds/plans/ endpoint using the software ID as a filter.
// Returns the template ID if found, otherwise returns an error.
func (c *Client) GetTemplateId(projectID string, location string, planName string, softwareID string) (int, error) {
	url := c.Api_endpoint + "rds/plans/"

	// Build GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, err
	}

	// Add query parameters
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", projectID)
	params.Add("location", location)
	params.Add("software_id", softwareID)
	req.URL.RawQuery = params.Encode()

	// Add authentication headers
	SetBasicHeaders(c.Auth_token, req)

	// Send request
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("[ERROR] GetTemplateId request failed: %v", err)
		return -1, err
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read GetTemplateId response body: %v", err)
		return -1, err
	}

	// Parse JSON response into PlanResponse model
	var res models.PlanResponse
	if err := json.Unmarshal(body, &res); err != nil {
		log.Printf("[ERROR] Failed to unmarshal GetTemplateId response: %v", err)
		return -1, err
	}

	// Search for a matching plan name
	for _, item := range res.Data.TemplatePlans {
		if item.PlanName == planName {
			return item.PlanTemplateID, nil
		}
	}

	log.Printf("[INFO] Template ID not found for plan name: %s", planName)
	return -1, errors.New("matching plan not found")
}



func (c *Client) ExpandVpcList(vpcIDs []string, projectID, location string) ([]models.VPCMetadata, error) {
	var vpcDetails []models.VPCMetadata

	for _, id := range vpcIDs {
		vpcResp, err := c.GetVpc(id, projectID, location)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch VPC details for ID %s: %v", id, err)
		}

		vpc := vpcResp.Data
		if vpc.State != "Active" {
			return nil, fmt.Errorf("cannot attach VPC %s: VPC is in '%s' state", id, vpc.State)
		}

		// Convert float64 to string (e.g., 123 -> "123")
		networkID := fmt.Sprintf("%.0f", vpc.Network_id)

		vpcDetails = append(vpcDetails, models.VPCMetadata{
			NetworkID: networkID,
			VPCName:   vpc.Name,
			IPv4CIDR:  vpc.Ipv4_cidr,
		})
	}

	return vpcDetails, nil
}



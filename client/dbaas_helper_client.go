package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (c *Client) GetSoftwareId(project_id string, location string, name string, version string) (int, error) {

	url := c.Api_endpoint + "rds/plans/"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, err
	}

	req, err = c.AddParamsAndHeader(req, location, project_id)
	if err != nil {
		return -1, fmt.Errorf(" [ERROR] error while setting parameters and headers =: %s ", err)
	}

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return -1, fmt.Errorf(" [ERROR] error inside GetSoftwareId =: %s ", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return -1, err
	}

	var res models.PlanResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return -1, fmt.Errorf(" [ERROR] inside GetSoftwareId | error while unmarshalling =: %s ", err)
	}

	for _, item := range res.Data.DatabaseEngines {
		if item.EngineName == name && item.EngineVersion == version {
			return item.EngineID, nil
		}
	}

	return -1, errors.New(" [ERROR] matching engine not found ")

}

func (c *Client) GetTemplateId(project_id string, location string, plan string, software_id string) (int, error) {

	url := c.Api_endpoint + "rds/plans/"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	params.Add("location", location)
	params.Add("software_id", software_id)
	req.URL.RawQuery = params.Encode()

	SetBasicHeaders(c.Auth_token, req)

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return -1, fmt.Errorf(" [ERROR] error inside GetTemplateId =: %s ", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return -1, err
	}

	var res models.PlanResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return -1, fmt.Errorf(" [ERROR] inside GetTemplateId | error while unmarshalling =: %s ", err)
	}

	for _, item := range res.Data.TemplatePlans {
		if item.PlanName == plan {
			return item.PlanTemplateID, nil
		}
	}

	return -1, errors.New("[ERROR] matching plan not found ")

}

func (c *Client) ExpandVpcList(d *schema.ResourceData, vpc_list []interface{}) ([]models.VPC, error) {
	var vpc_details []models.VPC

	for _, id := range vpc_list {
		vpc_detail, err := c.GetVpc(strconv.Itoa(id.(int)), d.Get("project_id").(string), d.Get("location").(string))
		if err != nil {
			return nil, err
		}
		data := vpc_detail.Data
		if data.State != "Active" {
			return nil, fmt.Errorf("[ERROR] Can not attach vpc currently, vpc is in %s state", data.State)
		}
		r := models.VPC{
			Network_id: data.Network_id,
			VpcName:    data.Name,
			Ipv4_cidr:  data.Ipv4_cidr,
		}

		vpc_details = append(vpc_details, r)
	}
	return vpc_details, nil
}

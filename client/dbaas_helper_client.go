package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (c *Client) GetSoftwareId(project_id string, location string, name string, version string) (int, error) {

	urlGetReserveIps := c.Api_endpoint + "rds/plans/"
	req, err := http.NewRequest("GET", urlGetReserveIps, nil)
	if err != nil {
		return -1, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	params.Add("location", location)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	if err != nil {
		log.Printf("[INFO] error inside GetSoftwareId")
		return -1, err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	log.Println("[DEBUG] Raw body:", string(body))

	res := models.PlanResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("[INFO] inside GetSoftwareId | error while unmarshlling")
		return -1, err
	}
	data := res.Data.DatabaseEngines
	for _, item := range data {
		if item.EngineName == name && item.EngineVersion == version {
			return item.EngineID, nil
		}
	}

	log.Printf("[INFO] ---- SOFTWARE ID ---- inside GetSoftwareId | error NOT FOUND")
	return -1, errors.New("matching engine not found")

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
		log.Printf("[INFO] error inside GetTemplateId")
		return -1, err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	res := models.PlanResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("[INFO] inside GetTemplateId | error while unmarshlling")
		return -1, err
	}
	data := res.Data.TemplatePlans
	for _, item := range data {
		if item.PlanName == plan {
			return item.PlanTemplateID, nil
		}
	}

	log.Printf("[INFO] ---- Template ID ---- inside GetTemplateId | error NOT FOUND")
	return -1, errors.New("matching plan not found")

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
			return nil, fmt.Errorf("Can not attach vpc currently, vpc is in %s state", data.State)
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

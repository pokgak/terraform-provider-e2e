package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func (c *Client) NewMySqlDb(item *models.MySqlCreate, project_id string) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}

	UrlEndPoint := c.Api_endpoint + "/rds/cluster/"
	log.Printf("[INFO] CLIENT NEWMYSQLDB | Endpoint: %s", UrlEndPoint)

	req, err := http.NewRequest("POST", UrlEndPoint, &buf)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	params.Add("location", item.Location)

	req.URL.RawQuery = params.Encode()

	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}
	err = CheckResponseStatus(response)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) GetSoftwareId(project_id string, location string, name string, version string) (int, error) {
	url := c.Api_endpoint + "rds/plans/"

	req, err := http.NewRequest("GET", url, nil)
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
	if err != nil {
		return -1, err
	}

	var res models.PlanResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("[INFO] inside GetSoftwareId | error while unmarshalling")
		return -1, err
	}

	for _, item := range res.Data.DatabaseEngines {
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
	if err != nil {
		return -1, err
	}

	var res models.PlanResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("[INFO] inside GetTemplateId | error while unmarshalling")
		return -1, err
	}

	for _, item := range res.Data.TemplatePlans {
		if item.PlanName == plan {
			return item.PlanTemplateID, nil
		}
	}

	log.Printf("[INFO] ---- Template ID ---- inside GetTemplateId | error NOT FOUND")
	return -1, errors.New("matching plan not found")
}

func (c *Client) GetMySqlDbaas(mySqlDBaaSId string, project_id string, location string) (*models.ResponseMySql, error) {
	urlGetDBaaSMySQL := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/"
	req, err := http.NewRequest("GET", urlGetDBaaSMySQL, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()

	SetBasicHeaders(c.Auth_token, req)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("[ERROR] client | error making request in GetMySqlDbaas: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] client | error reading response body: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status code: %d - body: %s", resp.StatusCode, string(body))
	}

	var res models.ResponseMySql
	if err := json.Unmarshal(body, &res); err != nil {
		log.Printf("[ERROR] GetMySqlDbaas | error unmarshalling JSON: %v", err)
		return nil, err
	}

	return &res, nil
}

func (c *Client) DeleteMySqlDBaaS(mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/"
	log.Printf("[INFO] %s", urlDBaaSMySql)
	req, err := http.NewRequest("DELETE", urlDBaaSMySql, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) ResumeMySqlDBaaS(mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/resume"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) StopMySqlDBaaS(mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/shutdown"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) RestartMySqlDBaaS(mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/restart"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) AttachVpcToMySql(item *models.AttachDetachVPC, mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}

	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/vpc-attach/"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, &buf)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) DetachVpcFromMySql(item *models.AttachDetachVPC, mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}

	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/vpc-detach/"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, &buf)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) AttachPGToMySqlDBaaS(mySqlDBaaSId string, ParameterGroupId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/parameter-group/" + ParameterGroupId + "/add"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) DetachPGFromMySqlDBaaS(mySqlDBaaSId string, ParameterGroupId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/parameter-group/" + ParameterGroupId + "/detach"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) AttachPublicIPToMySql(mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/public-ip-attach/"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) DetachPublicIPFromMySql(mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/public-ip-detach/"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) UpgradeMySQLPlan(dbaas_id string, template_id int, project_id string, location string) (interface{}, error) {
	dbaas_action := models.MySQlPlanUpgradeAction{
		TemplateID: template_id,
	}

	dbaasAction, err := json.Marshal(dbaas_action)
	if err != nil {
		return nil, err
	}

	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + dbaas_id + "/rds-upgrade/"
	req, err := http.NewRequest("PUT", urlDBaaSMySql, bytes.NewBuffer(dbaasAction))
	if err != nil {
		log.Printf("[INFO] error inside upgrade dbaas MySQL plan: %v", err)
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	params.Add("location", location)
	req.URL.RawQuery = params.Encode()

	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")

	response, err := c.HttpClient.Do(req)

	log.Printf("CLIENT UPGRADE MySQL PLAN | request = %+v", req)
	if response != nil {
		log.Printf("CLIENT UPGRADE MySQL PLAN | STATUS_CODE: %d, response = %+v", response.StatusCode, response)
	}

	if err == nil {
		err = CheckResponseStatus(response)
	}

	if err != nil {
		log.Printf("[INFO] error inside upgrade MySQL plan after CheckResponseStatus: %v", err)
		return nil, err
	}

	return response, nil
}

func (c *Client) ExpnadDisk(dbaas_id string, size int, project_id string, location string) (interface{}, error) {
	dbaas_action := models.MYSQLExpandDisk{
		Size: size,
	}

	dbaasAction, err := json.Marshal(dbaas_action)
	if err != nil {
		return nil, err
	}

	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + dbaas_id + "/disk-upgrade/"
	req, err := http.NewRequest("PUT", urlDBaaSMySql, bytes.NewBuffer(dbaasAction))
	if err != nil {
		log.Printf("[INFO] error inside upgrade dbaas MySQL plan: %v", err)
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	params.Add("location", location)
	req.URL.RawQuery = params.Encode()

	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")

	response, err := c.HttpClient.Do(req)

	log.Printf("CLIENT UPGRADE MySQL PLAN | request = %+v", req)
	if response != nil {
		log.Printf("CLIENT UPGRADE MySQL PLAN | STATUS_CODE: %d, response = %+v", response.StatusCode, response)
	}

	if err == nil {
		err = CheckResponseStatus(response)
	}

	if err != nil {
		log.Printf("[INFO] error inside upgrade MySQL plan after CheckResponseStatus: %v", err)
		return nil, err
	}

	return response, nil
}

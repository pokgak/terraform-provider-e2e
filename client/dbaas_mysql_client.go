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
		return nil, fmt.Errorf("======== [ERROR] error while creating http request =========: %s =========", err)
	}
	err = CheckResponseStatus(response)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] error while checking response code =========: %s =========", err)
	}
	defer response.Body.Close()
	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] error while unmarshling response =========: %s =========", err)
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
		return -1, fmt.Errorf("======== [ERROR] error inside GetSoftwareId =========: %s =========", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return -1, err
	}

	var res models.PlanResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return -1, fmt.Errorf("======== [ERROR] inside GetSoftwareId | error while unmarshalling =: %s =========", err)
	}

	for _, item := range res.Data.DatabaseEngines {
		if item.EngineName == name && item.EngineVersion == version {
			return item.EngineID, nil
		}
	}

	return -1, errors.New("===== [ERROR] matching engine not found ==========")
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
		return -1, fmt.Errorf("======== [ERROR] error inside GetTemplateId =========: %s =========", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return -1, err
	}

	var res models.PlanResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return -1, fmt.Errorf("======== [ERROR] inside GetTemplateId | error while unmarshalling =: %s =========", err)
	}

	for _, item := range res.Data.TemplatePlans {
		if item.PlanName == plan {
			return item.PlanTemplateID, nil
		}
	}

	return -1, errors.New("=======[ERROR] matching plan not found ============")
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
		return nil, fmt.Errorf("======== [ERROR] client | error making request in GetMySqlDbaas: %v=========", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] client | error reading response body: %v=========", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("==== [ERROR] API returned non-200 status code: %d - body: %s =========", resp.StatusCode, string(body))
	}

	var res models.ResponseMySql
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("======== [ERROR] GetMySqlDbaas | error unmarshalling JSON: %s =========", err)
	}

	return &res, nil
}

func (c *Client) DeleteMySqlDBaaS(mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/"
	log.Printf("[INFO] %s", urlDBaaSMySql)
	req, err := http.NewRequest("DELETE", urlDBaaSMySql, nil)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] DeleteMySqlDBaaS | error while creating http request: %s =========", err)
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] DeleteMySqlDBaaS | error while making http request: %s =========", err)
	}
	defer response.Body.Close()
	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] DeleteMySqlDBaaS | error unmarshalling JSON: %s =========", err)
	}
	return jsonRes, nil
}

func (c *Client) ResumeMySqlDBaaS(mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/resume"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, nil)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] ResumeMySqlDBaaS | error while creating http request: %s =========", err)
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] ResumeMySqlDBaaS | error while making http request: %s =========", err)
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] ResumeMySqlDBaaS | error unmarshalling JSON: %s =========", err)
	}
	return jsonRes, nil
}

func (c *Client) StopMySqlDBaaS(mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/shutdown"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, nil)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] StopMySqlDBaaS | error while creating http request: %s =========", err)
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] StopMySqlDBaaS  | error while making http request: %s =========", err)
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] StopMySqlDBaaS | error unmarshalling JSON: %s =========", err)
	}
	return jsonRes, nil
}

func (c *Client) RestartMySqlDBaaS(mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/restart"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, nil)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] RestartMySqlDBaaS | error while creating http request: %s =========", err)
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] RestartMySqlDBaaS | error while making http request: %s =========", err)
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] RestartMySqlDBaaS | error unmarshalling JSON: %s =========", err)
	}
	return jsonRes, nil
}

func (c *Client) AttachVpcToMySql(item *models.AttachDetachVPC, mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] AttachVpcToMySql | error while encoding buffer: %s =========", err)
	}

	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/vpc-attach/"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, &buf)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] AttachVpcToMySql | error while creating http request: %s =========", err)
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] AttachVpcToMySql | error while making http request: %s =========", err)
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] AttachVpcToMySql | error unmarshalling JSON: %s =========", err)
	}
	return jsonRes, nil
}

func (c *Client) DetachVpcFromMySql(item *models.AttachDetachVPC, mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] DetachVpcFromMySql | error while encoding buffer: %s =========", err)
	}

	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/vpc-detach/"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, &buf)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] DetachVpcFromMySql | error while creating http request: %s =========", err)
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] DetachVpcFromMySql | error while making http request: %s =========", err)
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] DetachVpcFromMySql | error unmarshalling JSON: %s =========", err)
	}
	return jsonRes, nil
}

func (c *Client) AttachPGToMySqlDBaaS(mySqlDBaaSId string, ParameterGroupId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/parameter-group/" + ParameterGroupId + "/add"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, nil)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] AttachPGToMySqlDBaaS | error while creating http request: %s =========", err)
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] AttachPGToMySqlDBaaS | error while making http request: %s =========", err)
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] AttachPGToMySqlDBaaS | error unmarshalling JSON: %s =========", err)
	}
	return jsonRes, nil
}

func (c *Client) DetachPGFromMySqlDBaaS(mySqlDBaaSId string, ParameterGroupId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/parameter-group/" + ParameterGroupId + "/detach"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, nil)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] DetachPGFromMySqlDBaaS | error while creating http request: %s =========", err)
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] DetachPGFromMySqlDBaaS | error while making http request: %s =========", err)
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] DetachPGFromMySqlDBaaS | error unmarshalling JSON: %s =========", err)
	}
	return jsonRes, nil
}

func (c *Client) AttachPublicIPToMySql(mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/public-ip-attach/"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, nil)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] AttachPublicIPToMySql | error while creating http request: %s =========", err)
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] AttachPublicIPToMySql | error while making http request: %s =========", err)
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] AttachPublicIPToMySql | error unmarshalling JSON: %s =========", err)
	}
	return jsonRes, nil
}

func (c *Client) DetachPublicIPFromMySql(mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/public-ip-detach/"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, nil)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] DetachPublicIPFromMySql | error while creating http request: %s =========", err)
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	SetBasicHeaders(c.Auth_token, req)
	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] DetachPublicIPFromMySql | error while making http request: %s =========", err)
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] DetachPublicIPFromMySql | error unmarshalling JSON: %s =========", err)
	}
	return jsonRes, nil
}

func (c *Client) UpgradeMySQLPlan(dbaas_id string, template_id int, project_id string, location string) (interface{}, error) {
	dbaas_action := models.MySQlPlanUpgradeAction{
		TemplateID: template_id,
	}

	dbaasAction, err := json.Marshal(dbaas_action)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] UpgradeMySQLPlan | error unmarshalling JSON: %s =========", err)
	}

	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + dbaas_id + "/rds-upgrade/"
	req, err := http.NewRequest("PUT", urlDBaaSMySql, bytes.NewBuffer(dbaasAction))
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] UpgradeMySQLPlan | error while creating http request: %s =========", err)
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	params.Add("location", location)
	req.URL.RawQuery = params.Encode()

	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")

	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] UpgradeMySQLPlan | error while making http request: %s =========", err)
	}

	err = CheckResponseStatus(response)

	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] UpgradeMySQLPlan | error inside upgrade MySQL plan after CheckResponseStatus: %s =========", err)
	}

	return response, nil
}

func (c *Client) ExpandMySQLDBaaSDisk(dbaas_id string, size int, project_id string, location string) (interface{}, error) {
	dbaas_action := models.MYSQLExpandDisk{
		Size: size,
	}

	dbaasAction, err := json.Marshal(dbaas_action)
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] ExpandMySQLDBaaSDisk | error unmarshalling JSON: %s =========", err)
	}

	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + dbaas_id + "/disk-upgrade/"
	req, err := http.NewRequest("PUT", urlDBaaSMySql, bytes.NewBuffer(dbaasAction))
	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] ExpandMySQLDBaaSDisk | error while creating http request: %s =========", err)
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	params.Add("location", location)
	req.URL.RawQuery = params.Encode()

	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")

	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] ExpandMySQLDBaaSDisk | error while creating http request: %s =========", err)
	}

	err = CheckResponseStatus(response)

	if err != nil {
		return nil, fmt.Errorf("======== [ERROR] ExpandMySQLDBaaSDisk | error inside upgrade MySQL plan after CheckResponseStatus: %s =========", err)
	}

	return response, nil
}

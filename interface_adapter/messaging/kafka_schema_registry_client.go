package messaging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hdlproject/es-transaction-service/config"
	"github.com/hdlproject/es-transaction-service/helper"
)

type (
	KSRClient struct {
		BaseURL    *url.URL
		HTTPClient *http.Client
	}

	ksrAddSchemaRequest struct {
		Schema     string `json:"schema"`
		SchemaType string `json:"schemaType"`
	}

	ksrAddSchemaResponse struct {
		Id int `json:"id"`
	}

	KSRJSONSchema struct {
		Type       string                           `json:"type"`
		Properties map[string]KSRJSONSchemaProperty `json:"properties,omitempty"`
	}

	KSRJSONSchemaProperty struct {
		Type string `json:"type"`
	}
)

var ksrClient *KSRClient

func GetKSRClient(config config.KafkaSchemaRegistry) (*KSRClient, error) {
	if ksqldbClient == nil {
		client, err := newKSRClient(config)
		if err != nil {
			return nil, helper.WrapError(err)
		}

		ksrClient = client
	}

	return ksrClient, nil
}

func newKSRClient(config config.KafkaSchemaRegistry) (*KSRClient, error) {
	client := &http.Client{
		Timeout: time.Minute,
	}

	address := fmt.Sprintf("http://%s:%s", config.Host, config.Port)

	baseUrl, err := url.Parse(address)
	if err != nil {
		return nil, helper.WrapError(err)
	}

	return &KSRClient{
		BaseURL:    baseUrl,
		HTTPClient: client,
	}, nil
}

func (instance *KSRClient) AddSchema(subject, schema string) (int, error) {
	endpoint := instance.BaseURL.JoinPath(fmt.Sprintf("/subjects/%s/versions", subject)).String()

	request := ksrAddSchemaRequest{
		Schema:     schema,
		SchemaType: "JSON",
	}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return 0, helper.WrapError(err)
	}

	resp, err := instance.HTTPClient.Post(endpoint, "application/vnd.schemaregistry.v1+json", bytes.NewBuffer(jsonRequest))
	if err != nil {
		return 0, helper.WrapError(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, helper.WrapError(err)
	}
	fmt.Println(string(body))

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error with status code %d", resp.StatusCode)
	}

	var response ksrAddSchemaResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, helper.WrapError(err)
	}

	return response.Id, nil
}

func (instance *KSRClient) GetSchemaById(id int) (KSRJSONSchema, error) {
	endpoint := instance.BaseURL.JoinPath(fmt.Sprintf("/schemas/ids/%d", id)).String()

	resp, err := instance.HTTPClient.Get(endpoint)
	if err != nil {
		return KSRJSONSchema{}, helper.WrapError(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return KSRJSONSchema{}, helper.WrapError(err)
	}
	fmt.Println(string(body))

	if resp.StatusCode != http.StatusOK {
		return KSRJSONSchema{}, fmt.Errorf("error with status code %d", resp.StatusCode)
	}

	var response ksrAddSchemaRequest
	err = json.Unmarshal(body, &response)
	if err != nil {
		return KSRJSONSchema{}, helper.WrapError(err)
	}

	var schemaResponse KSRJSONSchema
	err = json.Unmarshal([]byte(response.Schema), &schemaResponse)
	if err != nil {
		return KSRJSONSchema{}, helper.WrapError(err)
	}

	return schemaResponse, nil
}

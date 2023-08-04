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
	KSQLDBClient struct {
		BaseURL    *url.URL
		HTTPClient *http.Client
	}

	ksqlKSQLRequest struct {
		KSQL string `json:"ksql"`
	}
)

var ksqldbClient *KSQLDBClient

func GetKSQLDBClient(config config.KSQLDB) (*KSQLDBClient, error) {
	if ksqldbClient == nil {
		client, err := newKSQLDBClient(config)
		if err != nil {
			return nil, helper.WrapError(err)
		}

		ksqldbClient = client
	}

	return ksqldbClient, nil
}

func newKSQLDBClient(config config.KSQLDB) (*KSQLDBClient, error) {
	client := &http.Client{
		Timeout: time.Minute,
	}

	address := fmt.Sprintf("http://%s:%s", config.Host, config.Port)

	baseUrl, err := url.Parse(address)
	if err != nil {
		return nil, helper.WrapError(err)
	}

	return &KSQLDBClient{
		BaseURL:    baseUrl,
		HTTPClient: client,
	}, nil
}

func (instance *KSQLDBClient) Insert(query string) error {
	endpoint := instance.BaseURL.JoinPath("/ksql").String()

	request := ksqlKSQLRequest{
		KSQL: query,
	}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return helper.WrapError(err)
	}

	resp, err := instance.HTTPClient.Post(endpoint, "application/vnd.ksql.v1+json", bytes.NewBuffer(jsonRequest))
	if err != nil {
		return helper.WrapError(err)
	}

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println(string(respBody))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error with status code %d", resp.StatusCode)
	}

	return nil
}

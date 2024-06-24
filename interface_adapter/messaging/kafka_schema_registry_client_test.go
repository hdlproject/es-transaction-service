package messaging

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hdlproject/es-transaction-service/config"
)

func TestKSQLDBClient_AddSchema(t *testing.T) {
	ksrConfig := config.KafkaSchemaRegistry{
		Host: "localhost",
		Port: "8081",
	}

	client, err := GetKSRClient(ksrConfig)
	if err != nil {
		t.Fatal(err)
	}

	expectedSchema := KSRJSONSchema{
		Type: "object",
		Properties: map[string]KSRJSONSchemaProperty{
			"message": {
				Type: "string",
			},
		},
	}
	jsonSchema, err := json.Marshal(expectedSchema)
	if err != nil {
		t.Fatal(err)
	}

	topic := "custom-events"
	subject := fmt.Sprintf("%s-value", topic)
	id, err := client.AddSchema(subject, string(jsonSchema))
	if err != nil {
		t.Fatal(err)
	}

	schema, err := client.GetSchemaById(id)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(expectedSchema, schema); diff != "" {
		t.Fatalf("(-want/+got)\n%s", diff)
	}

	expectedSchema = KSRJSONSchema{
		Type: "string",
	}
	jsonSchema, err = json.Marshal(expectedSchema)
	if err != nil {
		t.Fatal(err)
	}

	subject = fmt.Sprintf("%s-key", topic)
	id, err = client.AddSchema(subject, string(jsonSchema))
	if err != nil {
		t.Fatal(err)
	}

	schema, err = client.GetSchemaById(id)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(expectedSchema, schema); diff != "" {
		t.Fatalf("(-want/+got)\n%s", diff)
	}
}

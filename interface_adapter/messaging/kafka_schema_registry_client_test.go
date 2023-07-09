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
			"id": {
				Type: "string",
			},
			"userid": {
				Type: "integer",
			},
			"amount": {
				Type: "integer",
			},
		},
	}

	jsonSchema, err := json.Marshal(expectedSchema)
	if err != nil {
		t.Fatal(err)
	}

	topic := "top-up-events"
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
}

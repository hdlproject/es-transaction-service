package messaging

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/hdlproject/es-transaction-service/config"
)

func TestKSQLDBClient_Insert(t *testing.T) {
	testTimeout := 10 * time.Minute

	ksqldbConfig := config.KSQLDB{
		Host: "localhost",
		Port: "8088",
	}

	client, err := GetKSQLDBClient(ksqldbConfig)
	if err != nil {
		t.Fatal(err)
	}

	kafkaConfig := config.Kafka{
		Host: "localhost",
		Port: "29092",
	}

	consumer, err := GetKafkaConsumer(kafkaConfig)
	if err != nil {
		t.Fatal(err)
	}

	topics := []string{"top-up-events"}
	id := uuid.NewString()
	userId := 1
	amount := 10000
	expectedOutput := fmt.Sprintf(`{"ID":"%s","USER_ID":%d,"AMOUNT":%d}`, id, userId, amount)

	err = client.Insert(fmt.Sprintf("INSERT INTO TOP_UP_EVENTS (id, user_id, amount) VALUES ('%s', %d, %d);", id, userId, amount))
	if err != nil {
		t.Fatal(err)
		return
	}

	result := make(chan string)

	err = consumer.Consume(topics, testTimeout, func(s string, err error) {
		if err != nil {
			t.Errorf("expect nil error")
			return
		}

		if s != expectedOutput {
			t.Errorf("expect message %s but got %s", expectedOutput, s)
			return
		}

		result <- s
	})
	if err != nil {
		t.Fatal(err)
		return
	}

	select {
	case m := <-result:
		fmt.Println(m)
	case <-time.Tick(testTimeout):
		t.Fatalf("test timeout")
	}
}

func TestKSQLDBClient_Publish(t *testing.T) {
	testTimeout := 10 * time.Minute

	ctx, cancelFunc := context.WithTimeout(context.Background(), testTimeout)
	defer cancelFunc()

	kafkaConfig := config.Kafka{
		Host: "localhost",
		Port: "29092",
	}

	consumer, err := GetKafkaConsumer(kafkaConfig)
	if err != nil {
		t.Fatal(err)
	}

	producer, err := GetKafkaProducer(kafkaConfig)
	if err != nil {
		t.Fatal(err)
	}

	topics := []string{"top-up-events"}
	id := uuid.NewString()
	userId := 1
	amount := 10000
	expectedOutput := fmt.Sprintf(`{"ID":"%s","USER_ID":%d,"AMOUNT":%d}`, id, userId, amount)

	result := make(chan string)

	err = consumer.Consume(topics, testTimeout, func(s string, err error) {
		if err != nil {
			t.Errorf("expect nil error")
			return
		}

		if s != expectedOutput {
			t.Errorf("expect message %s but got %s", expectedOutput, s)
			return
		}

		result <- s
	})
	if err != nil {
		t.Fatal(err)
		return
	}

	// wait for the consumer to be ready
	time.Sleep(3 * time.Second)

	err = producer.Produce(ctx, topics[0], expectedOutput, func(s string, err error) {
		if err != nil {
			t.Errorf("expect nil error")
			return
		}

		if s == "" {
			t.Errorf("expect not empty message")
		}
	})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case m := <-result:
		fmt.Println(m)
	case <-time.Tick(testTimeout):
		t.Fatalf("test timeout")
	}
}

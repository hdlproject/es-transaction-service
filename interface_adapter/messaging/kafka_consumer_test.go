package messaging

import (
	"context"
	"testing"
	"time"

	"github.com/hdlproject/es-transaction-service/config"
)

func TestKafkaConsumer_Consume(t *testing.T) {
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

	topics := []string{"custom-events"}

	result := make(chan string)

	err = consumer.Consume(topics, testTimeout, func(s string, err error) {
		if err != nil {
			t.Errorf("expect nil error")
			return
		}

		if s == "" {
			t.Errorf("expect not empty message")
			return
		}

		result <- s
	})
	if err != nil {
		t.Fatal(err)
	}

	// wait for the consumer to be ready
	time.Sleep(3 * time.Second)

	err = producer.Produce(ctx, topics[0], `{"message": "test message 2"}`, func(s string, err error) {
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
	case <-result:
	case <-time.Tick(testTimeout):
		t.Fatalf("test timeout")
	}
}

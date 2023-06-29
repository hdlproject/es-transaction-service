package messaging

import (
	"testing"
	"time"
)

func TestKafkaConsumer_Consume(t *testing.T) {
	testTimeout := 10 * time.Minute

	consumer, err := GetKafkaConsumer()
	if err != nil {
		t.Fatal(err)
	}

	producer, err := GetKafkaProducer()
	if err != nil {
		t.Fatal(err)
	}

	topics := []string{"test"}

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
		t.Error(err)
		return
	}

	// wait for the consumer to be ready
	time.Sleep(3 * time.Second)

	err = producer.Produce(topics[0], "test message 2", func(s string, err error) {
		if err != nil {
			t.Errorf("expect nil error")
			return
		}

		if s == "" {
			t.Errorf("expect not empty message")
		}
	})
	if err != nil {
		t.Error(err)
		return
	}

	select {
	case <-result:
	case <-time.Tick(testTimeout):
		t.Fatalf("test timeout")
	}
}

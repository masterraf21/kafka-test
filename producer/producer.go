package main

import (
	"encoding/json"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	SingleTopic   = "single_consumer_topic"
	MultipleTopic = "multiple_consumer_topic"
)

func main() {
	BROKERS := []string{"localhost:9001"}

	cmdPublishSingle := &cobra.Command{
		Use:   "publish-single",
		Short: "publish message to kafka for single consumer",
		Run: func(cmd *cobra.Command, args []string) {
			prod := NewKafkaProducer(BROKERS)
			message := make(map[string]interface{})
			message["first_key"] = "hello"
			message["second_key"] = "single consumer"
			if err := prod.Publish(SingleTopic, message); err != nil {
				log.Error("Error Publishing Message: %s", err)
			} else {
				log.Info("Success Publishing Message")
			}
		},
	}
	cmdPublishMultiple := &cobra.Command{
		Use:   "publish-multiple",
		Short: "publish message to kafka for multiple consumer",
		Run: func(cmd *cobra.Command, args []string) {
			prod := NewKafkaProducer(BROKERS)
			message := make(map[string]interface{})
			message["first_key"] = "hello"
			message["second_key"] = "multiple consumer"
			if err := prod.Publish(MultipleTopic, message); err != nil {
				log.Error("Error Publishing Message: %s", err)
			} else {
				log.Info("Success Publishing Message")
			}
		},
	}

	rootCmd := &cobra.Command{
		Use: "producer",
	}
	rootCmd.AddCommand(cmdPublishSingle, cmdPublishMultiple)
	rootCmd.SuggestionsMinimumDistance = 1
	rootCmd.Execute()
}

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(brokers []string) KafkaProducer {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 5
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		log.Fatalf("Error Creating Kafka Producer: %s", err)
	}

	return KafkaProducer{
		producer: producer,
	}
}

func (p *KafkaProducer) Publish(topic string, payload map[string]interface{}) error {
	jsonPayload, err := json.Marshal(&payload)
	if err != nil {
		log.Printf("Error marshalling payload: %s", err)
		return err
	}
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(string(jsonPayload)),
	}
	if _, _, err = p.producer.SendMessage(message); err != nil {
		log.Printf("Error publishing message: %s", err)
		return err
	}
	return nil
}

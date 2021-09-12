package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/Shopify/sarama"
)

const (
	GROUP         = "first_group"
	SingleTopic   = "single_consumer_topic"
	MultipleTopic = "multiple_consumer_topic"
)

func main() {
}

type (
	KafkaConsumer struct {
		consumer sarama.ConsumerGroup
	}

	ConsumerHandler struct {
		ready chan bool
		ctx   context.Context
	}
)

func NewKafkaConsumer(brokers []string) KafkaConsumer {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_8_0_0
	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := sarama.NewConsumerGroup(brokers, GROUP, cfg)
	if err != nil {
		log.Fatal("Error Creating Kafk Consumer: %s", err)
	}
	return KafkaConsumer{
		consumer: consumer,
	}
}

func (c *KafkaConsumer) Consume() {
	topics := []string{
		SingleTopic,
		MultipleTopic,
	}
	ctx, cancel := context.WithCancel(context.Background())
	handler := ConsumerHandler{
		ready: make(chan bool),
		ctx:   ctx,
	}

	wg := &sync.WaitGroup{}
	sigterm := make(chan os.Signal, 1)
	defer os.Exit(0)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := c.consumer.Consume(ctx, topics, handler); err != nil {
				log.Panicf("Error from consumer: %v", error)
			}
			if ctx.Err() != nil {
				return
			}
			handler.ready = make(chan bool)
		}
	}()
	<-handler.ready
	log.Println("Sarama consumer running....")

	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-sigterm:
		log.Println("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err := c.consumer.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}

func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	close(h.ready)
	return nil
}

func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
}

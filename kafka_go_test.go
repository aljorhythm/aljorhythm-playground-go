package main

import (
	"context"
	"testing"

	kafka "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

/**
 Implementation
 depend on interface instead of kafka-go's Writer
**/

type KafkaWriter interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

type TheService struct {
	KafkaWriter
}

func (s TheService) doSomething() {
	s.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("Key-A"),
		Value: []byte("Hello World!"),
	})
}

func exampleCaller() {
	w := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "topic-A",
		Balancer: &kafka.LeastBytes{},
	}
	realService := TheService{w}
	realService.doSomething()
}

/**
  Tests
**/

// can use mocking library (mockery etc..) for more convenient expectations
type MockWriter struct {
	messages []kafka.Message
}

func (w MockWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	w.messages = msgs
	return nil
}

func TestTheService(t *testing.T) {
	aMockWriter := MockWriter{}
	service := TheService{aMockWriter}
	t.Run("Sends message when something is done", func(t *testing.T) {
		service.doSomething()
		// dumb assertion for example
		assert.True(t, len(aMockWriter.messages) > 1)
	})
}

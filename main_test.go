package main

import (
	"testing"
	"time"

	"github.com/nats-io/stan.go"
	"github.com/stretchr/testify/assert"
)

func TestLoadOrderAndUpdateUID(t *testing.T) {
	testUID := "testUID"
	expectedDate := time.Now()

	order, err := loadOrderAndUpdateUID(testUID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if order.OrderUID != testUID {
		t.Errorf("Expected OrderUID to be %v, got %v", testUID, order.OrderUID)
	}

	if order.DateCreated.Before(expectedDate) {
		t.Errorf("Expected DateCreated to be after %v, got %v", expectedDate, order.DateCreated)
	}
}

func TestPublishToNATS(t *testing.T) {
	// Установите подключение к NATS Streaming
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("nats://localhost:4222"))
	defer sc.Close()
	assert.NoError(t, err)
}

package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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

func TestIntegration(t *testing.T) {
	// Инициализация сервера
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ваш обработчик
	}))
	defer server.Close()

	// Отправка запроса
	response, err := http.Post(server.URL, "application/x-www-form-urlencoded", strings.NewReader("id=testUID"))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", response.Status)
	}

	// Дополнительные проверки
}

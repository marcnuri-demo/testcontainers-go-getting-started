package main

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestExample(t *testing.T) {
	containerCtx := context.Background()
	chuckNorris, _ := testcontainers.GenericContainer(containerCtx, testcontainers.GenericContainerRequest{
		Started: true,
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "marcnuri/chuck-norris:latest",
			ExposedPorts: []string{"8080/tcp"},
			WaitingFor:   wait.ForLog("Listening on: http://0.0.0.0:8080"),
		},
	})
	defer func() {
		err := chuckNorris.Terminate(containerCtx)
		if err != nil {
			t.Error(err)
		}
	}()
	endpoint, err := chuckNorris.Endpoint(containerCtx, "http")
	if err != nil {
		t.Error(err)
	}
	response, err := http.Get(endpoint)
	if err != nil {
		t.Error(err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
	}
	bodyString := string(body)
	if !strings.Contains(strings.ToLower(bodyString), "chuck") {
		t.Errorf("Expected a Chuck Norris approved response, got \"%s\"", bodyString)
	}
}

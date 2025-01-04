package services

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
)

type EventBus struct {
	Name string `json:"Name"`
}

type EventBusService struct {
	client *eventbridge.Client
}

func NewEventBusService() (*EventBusService, error) {
	// Load the default AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	// Create an EventBridge client
	client := eventbridge.NewFromConfig(cfg)

	return &EventBusService{
		client: client,
	}, nil
}

func (ebs *EventBusService) List() ([]EventBus, error) {
	// Call the ListEventBuses API
	input := &eventbridge.ListEventBusesInput{}
	output, err := ebs.client.ListEventBuses(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to list event buses: %w", err)
	}

	// Map the SDK's EventBus structure to our custom EventBus type
	var eventBuses []EventBus
	for _, eb := range output.EventBuses {
		eventBuses = append(eventBuses, EventBus{
			Name: aws.ToString(eb.Name),
		})
	}

	return eventBuses, nil
}

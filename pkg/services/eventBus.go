package services

import (
	"encoding/json"
	"fmt"

	"github.com/aka-somix/aws-events-gate/internal/aws"
)

type EventBus struct {
	Name string `json:"Name"`
}

type EventBusService struct {
}

func NewEventBusService() *EventBusService {
	return &EventBusService{}
}

func (ebs *EventBusService) List() ([]EventBus, error) {
	// Run the AWS CLI command to list event buses
	out, err := aws.AwsCommand("aws", "events", "list-event-buses", "--output", "json").Output()
	if err != nil {
		fmt.Printf("Internal Error while retrieving Event Buses: %v\n", err)
		return nil, err
	}

	// Parse the JSON output
	var result struct {
		EventBuses []EventBus `json:"EventBuses"`
	}
	if err := json.Unmarshal(out, &result); err != nil {
		fmt.Printf("Error parsing AWS CLI output: %v\n", err)
		return nil, err
	}

	return result.EventBuses, nil
}

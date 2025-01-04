package services

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aka-somix/aws-events-gate/internal/aws"
)

// Global Config for Monitor creations
var cfg *MonitorConfig = &MonitorConfig{
	logGroupPrefix: "/aws-events-gate/watch",
	retentionDays: "1",
	policyName: "allow-logging-from-eventbridge",
	targetID: "aws-events-gate-log-group",
	ruleName: "aws-events-gate-rule",
	eventPattern: `{"source": [ { "wildcard": "*" }]}`,
}

type MonitorService struct {
}

func NewMonitorService() *MonitorService {
	return &MonitorService{}
}


func (ms *MonitorService) Create(eventBus string) error{
	fmt.Printf("Starting resource creation for EventBus: %s\n", eventBus)

	// Create CloudWatch Log Group
	logGroupName := fmt.Sprintf("%s/%s", cfg.logGroupPrefix, eventBus)

	if _, err := aws.AwsCommand("aws", "logs", "create-log-group", "--log-group-name", logGroupName).Output(); err != nil {
		fmt.Printf("Log Group already exists or error occurred: %v\n", err)
	} else {
		fmt.Printf("Log Group %s created successfully\n", logGroupName)
	}

	// Put retention policy
	if _, err := aws.AwsCommand("aws", "logs", "put-retention-policy", "--log-group-name", logGroupName, "--retention-in-days", cfg.retentionDays).Output(); err != nil {
		fmt.Printf("Retention policy already set or error occurred: %v\n", err)
	} else {
		fmt.Printf("Retention policy set for Log Group %s\n", logGroupName)
	}

	// Create Log Resource Policy
	policyDocument := map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Effect":   "Allow",
				"Principal": map[string][]string{"Service": {"events.amazonaws.com"}},
				"Action":   []string{"logs:*"},
				"Resource": fmt.Sprintf("arn:aws:logs:*:*:log-group:*:*"),
			},
		},
	}
	policyJSON, _ := json.Marshal(policyDocument)
    policyJSONStr := strconv.Quote(string(policyJSON)) // Escape quotes
	if err := aws.AwsCommand("aws", "logs", "put-resource-policy", "--policy-name", cfg.policyName, "--policy-document", policyJSONStr); err != nil {
		fmt.Printf("Log Resource Policy already exists or error occurred: %v\n", err)
	} else {
		fmt.Printf("Log Resource Policy %s created successfully\n", cfg.policyName)
	}

	// Create EventBridge Rul
    eventPatternStr := strconv.Quote(cfg.eventPattern) // Escape quotes for the event pattern
	if err := aws.AwsCommand("aws", "events", "put-rule", "--name", cfg.ruleName, "--event-bus-name", eventBus, "--event-pattern", eventPatternStr); err != nil {
		fmt.Printf("EventBridge Rule already exists or error occurred: %v\n", err)
	} else {
		fmt.Printf("EventBridge Rule %s created successfully\n", cfg.ruleName)
	}

	// Add Target to EventBridge Rule
	targetsStr := strconv.Quote(fmt.Sprintf(`[{"Id": "%s", "Arn": "arn:aws:logs:*:*:log-group:%s"}]`, cfg.targetID, logGroupName))
	if err := aws.AwsCommand("aws", "events", "put-targets", "--rule", cfg.ruleName, "--event-bus-name", eventBus, "--targets", targetsStr); err != nil {
		fmt.Printf("Target already exists or error occurred: %v\n", err)
	} else {
		fmt.Printf("Target added to Rule %s successfully\n", cfg.targetID)
	}

	return nil
}

func (ms *MonitorService) List() []string {
	return []string{"abcd"}
}

func (ms *MonitorService) Destroy(eventBus string) error{
	fmt.Printf("Starting resource deletion for EventBus: %s\n", eventBus)

	// Delete Target from EventBridge Rule
	if _, err := aws.AwsCommand("aws", "events", "remove-targets", "--rule", cfg.ruleName, "--event-bus-name", eventBus, "--ids", "CloudWatchLogs").Output(); err != nil {
		fmt.Printf("Target already removed or error occurred: %v\n", err)
	} else {
		fmt.Printf("Target removed from Rule %s successfully\n", cfg.ruleName)
	}

	// Delete EventBridge Rule
	if _, err := aws.AwsCommand("aws", "events", "delete-rule", "--name", cfg.ruleName, "--event-bus-name", eventBus).Output(); err != nil {
		fmt.Printf("Rule already deleted or error occurred: %v\n", err)
	} else {
		fmt.Printf("EventBridge Rule %s deleted successfully\n", cfg.ruleName)
	}

	// Delete CloudWatch Log Group
	logGroupName := fmt.Sprintf("%s/%s", cfg.logGroupPrefix, eventBus)
	if _, err := aws.AwsCommand("aws", "logs", "delete-log-group", "--log-group-name", logGroupName).Output(); err != nil {
		fmt.Printf("Log Group already deleted or error occurred: %v\n", err)
	} else {
		fmt.Printf("Log Group %s deleted successfully\n", logGroupName)
	}

	return nil
}

func (ms *MonitorService) Tail(eventBus string) string{
	return "abcd"
}

package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)


var cfg = &MonitorConfig{
	LogGroupPrefix: "/aws-events-gate/watch",
	RetentionDays:  1,
	PolicyName:     "allow-logging-from-eventbridge",
	TargetID:       "aws-events-gate-log-group",
	RuleName:       "aws-events-gate-rule",
	EventPattern:   `{"source": [ { "wildcard": "*" }]}`,
}

type MonitorService struct {
	cwLogsClient *cloudwatchlogs.Client
	ebClient     *eventbridge.Client
}

func NewMonitorService() (*MonitorService, error) {
	awsCfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &MonitorService{
		cwLogsClient: cloudwatchlogs.NewFromConfig(awsCfg),
		ebClient:     eventbridge.NewFromConfig(awsCfg),
	}, nil
}

func (ms *MonitorService) Create(eventBus string) error {
	fmt.Printf("Starting resource creation for EventBus: %s\n", eventBus)

	ctx := context.TODO()
	logGroupName := fmt.Sprintf("%s/%s", cfg.LogGroupPrefix, eventBus)

	// Create CloudWatch Log Group
	_, err := ms.cwLogsClient.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(logGroupName),
	})
	if err != nil {
		fmt.Printf("Log Group already exists or error occurred: %v\n", err)
	} else {
		fmt.Printf("Log Group %s created successfully\n", logGroupName)
	}

	// Set retention policy
	_, err = ms.cwLogsClient.PutRetentionPolicy(ctx, &cloudwatchlogs.PutRetentionPolicyInput{
		LogGroupName:    aws.String(logGroupName),
		RetentionInDays: &cfg.RetentionDays,
	})
	if err != nil {
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
	_, err = ms.cwLogsClient.PutResourcePolicy(ctx, &cloudwatchlogs.PutResourcePolicyInput{
		PolicyName:     aws.String(cfg.PolicyName),
		PolicyDocument: aws.String(string(policyJSON)),
	})
	if err != nil {
		fmt.Printf("Log Resource Policy already exists or error occurred: %v\n", err)
	} else {
		fmt.Printf("Log Resource Policy %s created successfully\n", cfg.PolicyName)
	}

	// Retrieve Created Log Group
	descOutput, descErr := ms.cwLogsClient.DescribeLogGroups(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(logGroupName),
	})
	if descErr != nil || len(descOutput.LogGroups) == 0 {
		return error(fmt.Errorf("Failed to Create Monitor: %w", descErr))
	}
	var logGroupArn = aws.ToString(descOutput.LogGroups[0].Arn)

	// Create EventBridge Rule
	_, err = ms.ebClient.PutRule(ctx, &eventbridge.PutRuleInput{
		Name:         aws.String(cfg.RuleName),
		EventBusName: aws.String(eventBus),
		EventPattern: aws.String(cfg.EventPattern),
	})
	if err != nil {
		fmt.Printf("EventBridge Rule already exists or error occurred: %v\n", err)
	} else {
		fmt.Printf("EventBridge Rule %s created successfully\n", cfg.RuleName)
	}

	// Add Target to EventBridge Rule
	target := types.Target{
		Id:  aws.String(cfg.TargetID),
		Arn: aws.String(logGroupArn),
	}
	_, err = ms.ebClient.PutTargets(ctx, &eventbridge.PutTargetsInput{
		Rule:         aws.String(cfg.RuleName),
		EventBusName: aws.String(eventBus),
		Targets:      []types.Target{target},
	})
	if err != nil {
		fmt.Printf("Target already exists or error occurred: %v\n", err)
	} else {
		fmt.Printf("Target added to Rule %s successfully\n", cfg.TargetID)
	}

	return nil
}

func (ms *MonitorService) Destroy(eventBus string) error {
	fmt.Printf("Starting resource deletion for EventBus: %s\n", eventBus)

	ctx := context.TODO()
	logGroupName := fmt.Sprintf("%s/%s", cfg.LogGroupPrefix, eventBus)

	// Remove Target from EventBridge Rule
	_, err := ms.ebClient.RemoveTargets(ctx, &eventbridge.RemoveTargetsInput{
		Rule:         aws.String(cfg.RuleName),
		EventBusName: aws.String(eventBus),
		Ids:          []string{cfg.TargetID},
	})
	if err != nil {
		fmt.Printf("Target already removed or error occurred: %v\n", err)
	} else {
		fmt.Printf("Target removed from Rule %s successfully\n", cfg.RuleName)
	}

	// Delete EventBridge Rule
	_, err = ms.ebClient.DeleteRule(ctx, &eventbridge.DeleteRuleInput{
		Name:         aws.String(cfg.RuleName),
		EventBusName: aws.String(eventBus),
	})
	if err != nil {
		fmt.Printf("Rule already deleted or error occurred: %v\n", err)
	} else {
		fmt.Printf("EventBridge Rule %s deleted successfully\n", cfg.RuleName)
	}

	// Delete CloudWatch Log Group
	_, err = ms.cwLogsClient.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
		LogGroupName: aws.String(logGroupName),
	})
	if err != nil {
		fmt.Printf("Log Group already deleted or error occurred: %v\n", err)
	} else {
		fmt.Printf("Log Group %s deleted successfully\n", logGroupName)
	}

	return nil
}

func (ms *MonitorService) List () []string {
	return []string{}
}

func (ms *MonitorService) Tail (eventBus string) []string {
	return []string{}
}

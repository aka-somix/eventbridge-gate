package services

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
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

func (ms *MonitorService) getArnFromLogGroupName(ctx context.Context, logGroupName string) (*cwtypes.LogGroup, error) {
	// Retrieve Created Log Group
	descOutput, descErr := ms.cwLogsClient.DescribeLogGroups(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(logGroupName),
	})
	if descErr != nil || len(descOutput.LogGroups) == 0 {
		return nil, descErr
	}
	var logGroupArn = descOutput.LogGroups[0]

	return &logGroupArn, nil
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

	// Retrieve Log Group Object from name
	logGroupObj, err := ms.getArnFromLogGroupName(ctx, logGroupName);
	if err != nil {
		return err
	}

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
		Arn: aws.String(*logGroupObj.Arn),
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

func (ms *MonitorService) List() ([]string, error) {
    ctx := context.TODO()
    var matchingEventBuses []string

    // List all EventBuses
    listBusesOutput, err := ms.ebClient.ListEventBuses(ctx, &eventbridge.ListEventBusesInput{})
    if err != nil {
        return nil, fmt.Errorf("failed to list event buses: %w", err)
    }

    // Iterate through each EventBus and check for the RuleName
    for _, eventBus := range listBusesOutput.EventBuses {
        rulesOutput, err := ms.ebClient.ListRules(ctx, &eventbridge.ListRulesInput{
            EventBusName: eventBus.Name,
        })
        if err != nil {
            fmt.Printf("Error listing rules for EventBus %s: %v\n", *eventBus.Name, err)
            continue
        }

        for _, rule := range rulesOutput.Rules {
            if *rule.Name == cfg.RuleName {
                matchingEventBuses = append(matchingEventBuses, *eventBus.Name)
                break
            }
        }
    }

    return matchingEventBuses, nil
}

func (ms *MonitorService) Tail(eventBus string) error {
	logGroupName := fmt.Sprintf("%s/%s", cfg.LogGroupPrefix, eventBus)

	var logGroupObj, err = ms.getArnFromLogGroupName(context.TODO(), logGroupName)

	// Create the input for starting live tail
	startLiveTailInput := &cloudwatchlogs.StartLiveTailInput{
		LogGroupIdentifiers: []string{*logGroupObj.LogGroupArn},
	}

	// Start live tailing logs
	liveTail, err := ms.cwLogsClient.StartLiveTail(context.TODO(), startLiveTailInput)
	if err != nil {
		return fmt.Errorf("unable to start live tail, %v", err)
	}

	// Print a message that tailing is starting
	fmt.Println("Tailing logs in real-time... Press 'q' to stop.")


	// Create a channel to listen for interrupt signals to gracefully exit the loop
	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, syscall.SIGINT, syscall.SIGTERM)

	streamingChannel := liveTail.GetStream().Events()

	go func() {
		// Use a range loop to read all messages from the channel
		for streamResult := range streamingChannel {
			// Use a type assertion or type switch to check and cast the type
			if sessionUpdate, ok := streamResult.(*cwtypes.StartLiveTailResponseStreamMemberSessionUpdate); ok {
				for _, message := range sessionUpdate.Value.SessionResults {
					fmt.Println(*message.Message)
				}
			}
		}
		fmt.Println("Closing Live tailing stream.")
	}()

	// Wait for the user to press 'q' to quit
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		if strings.TrimSpace(input) == "q" {
			break
		}
	}

	// Close the stream and handle cleanup
	err = liveTail.GetStream().Close()
	if err != nil {
		fmt.Println("Error closing stream: ", err)
	} else {
		fmt.Println("Live tail stream closed.")
	}

	return nil
}

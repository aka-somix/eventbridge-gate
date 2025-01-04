package services

type MonitorConfig struct {
	logGroupPrefix string	
	retentionDays string
	policyName string
	ruleName string
	eventPattern string
	targetID string
}
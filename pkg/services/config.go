package services

type MonitorConfig struct {
	LogGroupPrefix string
	RetentionDays  int32
	PolicyName     string
	TargetID       string
	RuleName       string
	EventPattern   string
}
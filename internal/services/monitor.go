package services


type MonitorService struct {
}

func NewMonitorService() *MonitorService {
	return &MonitorService{}
}

func (ms *MonitorService) Create(eventBus string) string{
	return "abcd"
}

func (ms *MonitorService) List(eventBus string) []string {
	return []string{"abcd"}
}

func (ms *MonitorService) Destroy(eventBus string) string{
	return "abcd"
}

func (ms *MonitorService) Tail(eventBus string) string{
	return "abcd"
}

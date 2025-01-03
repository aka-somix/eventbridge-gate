package services

import "github.com/aka-somix/aws-events-gate/internal/store"


type MonitorService struct {
	profileStore *store.ProfileStore
}

func NewMonitorService(profileStore *store.ProfileStore) *MonitorService {
	return &MonitorService{
		profileStore,
	}
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

package health

import (
	"github.com/DO-2K23-26/polypass-microservices/search-service/infrastructure"
)

type HealthService struct {
	esClient    infrastructure.ElasticAdapter
	kafkaClient infrastructure.KafkaAdapter
}

func NewHealthService(esClient infrastructure.ElasticAdapter, kafkaClient infrastructure.KafkaAdapter) *HealthService {
	return &HealthService{
		esClient:    esClient,
		kafkaClient: kafkaClient,
	}
}

func (s *HealthService) CheckHealth() HealthResponse {
	// The number of services to check
	serviceNumber := 2
	healthResponse := HealthResponse{}

	results := make(chan struct {
		service string
		status  string
	}, serviceNumber)

	// Run Elasticsearch health check in a goroutine
	go func() {
		esHealth := s.esClient.CheckHealth()
		status := Failure
		if esHealth {
			status = Ok
		}
		results <- struct {
			service string
			status  string
		}{service: "Elasticsearch", status: status}
	}()

	// Run Kafka health check in a goroutine
	go func() {
		kafkaHealth := s.kafkaClient.CheckHealth()
		status := Failure
		if kafkaHealth {
			status = Ok
		}
		results <- struct {
			service string
			status  string
		}{service: "Kafka", status: status}
	}()

	// Collect results
	for range serviceNumber {
		result := <-results
		healthResponse[result.service] = result.status
	}

	return healthResponse
}

package main

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/DO-2K23-26/polypass-microservices/organization/config"
	"github.com/DO-2K23-26/polypass-microservices/organization/infrastructure"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		config.HandleError(err)
		return
	}

	kafka, err := infrastructure.NewKafkaAdapter(cfg.KafkaHost, cfg.ClientId)
	if err != nil {
		log.Fatal(err)
	}
	subject := "create_folder-value"
	encoder, err := infrastructure.NewOrganizationEncoder(cfg.SchemaRegistryURL, subject)
	if err != nil {
		log.Fatal(err)
	}

	var messageCount uint64
	startTime := time.Now()
	numWorkers := 10 // Nombre de workers en parallèle
	var wg sync.WaitGroup

	// Fonction worker qui envoie des messages
	worker := func(workerID int) {
		defer wg.Done()
		for {
			data := map[string]interface{}{
				"id":   fmt.Sprintf("Workder_%d", workerID),
				"name": fmt.Sprintf("User_%d_Worker_%d", atomic.AddUint64(&messageCount, 1), workerID),
				// "age":  int(atomic.LoadUint64(&messageCount) % 100),
			}

			err = kafka.ProduceAvro("create_folder", encoder, data)
			if err != nil {
				log.Printf("Erreur envoi Kafka (Worker %d): %v", workerID, err)
				continue
			}

			if atomic.LoadUint64(&messageCount)%1000000 == 0 {
				elapsed := time.Since(startTime)
				rate := float64(atomic.LoadUint64(&messageCount)) / elapsed.Seconds()
				fmt.Printf("✅ %d messages envoyés (%.2f msg/sec)\n", atomic.LoadUint64(&messageCount), rate)
			}
		}
	}

	// Démarrer les workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i)
	}

	// Attendre indéfiniment (le programme continuera à tourner)
	wg.Wait()
}

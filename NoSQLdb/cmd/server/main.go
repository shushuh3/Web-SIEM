package main

import (
	"encoding/json"
	"log"
	"nosql_db/internal/config"
	"nosql_db/internal/server"
	"nosql_db/internal/storage"
	"os"
)

func main() {
	log.Println("Starting NoSQLdb server...")
	cfg := config.Load()

	// Загрузка начальных данных
	loadInitialData()

	srv := server.New(cfg.Host + ":" + cfg.Port)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}

func loadInitialData() {
	dataFile := "./data/security_events.json"

	// Проверяем существование файла
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		log.Printf("Initial data file not found: %s, skipping data load", dataFile)
		return
	}

	data, err := os.ReadFile(dataFile)
	if err != nil {
		log.Printf("Failed to read initial data file: %v", err)
		return
	}

	var events map[string]map[string]interface{}
	if err := json.Unmarshal(data, &events); err != nil {
		log.Printf("Failed to parse initial data: %v", err)
		return
	}

	if len(events) == 0 {
		log.Println("No events to load")
		return
	}

	// Загружаем события в коллекцию siem_events
	coll, err := storage.GlobalManager.GetCollection("siem_events")
	if err != nil {
		log.Printf("Failed to get collection: %v", err)
		return
	}

	count := 0
	for _, event := range events {
		if _, err := coll.Insert(event); err != nil {
			log.Printf("Failed to insert event: %v", err)
			continue
		}
		count++
	}

	if err := coll.Save(); err != nil {
		log.Printf("Failed to save collection: %v", err)
		return
	}

	if err := coll.SaveAllIndexes(); err != nil {
		log.Printf("Failed to save indexes: %v", err)
		return
	}

	log.Printf("Successfully loaded %d initial events into siem_events collection", count)
}

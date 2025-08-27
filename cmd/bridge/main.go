package main

import (
	"context"
	"log"

	"github.com/example/opcuamqttbridge/internal/bridge"
	dbpkg "github.com/example/opcuamqttbridge/internal/db"
	"github.com/example/opcuamqttbridge/internal/models"
)

func main() {
	db, err := dbpkg.Open()
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	dbpkg.AutoMigrate(db, &models.OPCUAServer{}, &models.MQTTBroker{}, &models.Point{})

	ctx := context.Background()
	if err := bridge.Start(ctx, db); err != nil {
		log.Fatalf("bridge failed: %v", err)
	}
}

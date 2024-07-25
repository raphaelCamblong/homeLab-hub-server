package main

import (
	"context"
	"internal/adapters/api"
	"internal/adapters/database"
	"internal/usecase"
	"log"
)

func main() {
	ctx := context.Background()
	db, err := database.NewDatabase("services.db")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := database.NewServiceRepository(db)
	uc := usecase.NewServiceUsecase(repo)
	handler := api.NewHandler(uc)

	server := api.NewServer(handler)
	server.Start()
}

package main

import (
	"bufio"
	"context"
	"fmt"
	"kororo/internal/adapters/config"
	"kororo/internal/adapters/llm/openrouter"
	mongodb "kororo/internal/adapters/storage/mongo"
	"kororo/internal/adapters/storage/mongo/repository"
	"kororo/internal/core/ports"
	"kororo/internal/core/services"
	"os"
	"strings"
)

func main() {

	var err error
	var config = config.NewConfig()
	var strongLLMAdapter ports.LLMAdapter
	var weakLLMAdapter ports.LLMAdapter

	var mongo *mongodb.Mongo
	mongo, err = mongodb.NewMongo(config)
	if err != nil {
		panic(err)
	}

	strongLLMAdapter = openrouter.New(config, "deepseek/deepseek-chat")
	weakLLMAdapter = openrouter.New(config, "google/gemini-2.0-flash-001")

	// logs
	var logger = services.NewLogService(config)

	// intention
	var intentionRepository = repository.NewIntentionRepository(mongo)
	var targetDectector = services.NewTargetDectector(weakLLMAdapter)
	var intentionService = services.NewIntentionService(intentionRepository, targetDectector, weakLLMAdapter)

	// action
	var actionRepository = repository.NewActionRepository(mongo)
	var actionService = services.NewActionService(actionRepository, weakLLMAdapter, logger)

	// field detector
	var fieldDetectorService = services.NewFieldDetectorService(weakLLMAdapter, logger)

	// intention proccess
	var intentionProccessService = services.NewIntentionProccess(intentionService, fieldDetectorService, actionService, logger)
	var multiIntentionProccessService = services.NewMultiIntentionProcessService(intentionProccessService, strongLLMAdapter, logger)

	// Lector de prompt:
	var reader = bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Prompt: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		var response, err = multiIntentionProccessService.Process(context.Background(), text)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		fmt.Printf("Response: %s\n\n", response)
	}

}

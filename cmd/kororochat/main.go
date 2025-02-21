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

	weakLLMAdapter = openrouter.New(config, "google/gemini-2.0-flash-001")
	//strongLLMAdapter = openrouter.New(config, "deepseek/deepseek-chat")
	strongLLMAdapter = weakLLMAdapter

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
	var multiIntentionChatProccessService = services.NewMultiIntentionChatProcessService(intentionProccessService, strongLLMAdapter, logger)

	// Lector de prompt:
	var reader = bufio.NewReader(os.Stdin)

	fmt.Print("Initial Message: ")
	initialMessage, _ := reader.ReadString('\n')
	initialMessage = strings.TrimSpace(initialMessage)

	multiIntentionChatProccessService.Process(context.TODO(), initialMessage)

}

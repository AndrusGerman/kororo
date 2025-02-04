package main

import (
	"bufio"
	"context"
	"fmt"
	"kororo/internal/adapters/config"
	"kororo/internal/adapters/llm/gemini"
	"kororo/internal/adapters/llm/huggingface"
	"kororo/internal/adapters/rest"
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
	var llmAdapter ports.LLMAdapter
	var mongo *mongodb.Mongo
	var rest = rest.New()

	mongo, err = mongodb.NewMongo(config)
	if err != nil {
		panic(err)
	}

	llmAdapter, err = gemini.New(config)
	if err != nil {
		panic(err)
	}

	llmAdapter = huggingface.New(rest, config, "deepseek-ai/DeepSeek-R1-Distill-Qwen-32B")

	// logs
	var logger = services.NewLogService(config)

	// intention
	var intentionRepository = repository.NewIntentionRepository(mongo)
	var targetDectector = services.NewTargetDectector(llmAdapter)
	var intentionService = services.NewIntentionService(intentionRepository, targetDectector, llmAdapter)

	// action
	var actionRepository = repository.NewActionRepository(mongo)
	var actionService = services.NewActionService(actionRepository, llmAdapter, logger)

	// field detector
	var fieldDetectorService = services.NewFieldDetectorService(llmAdapter)

	// intention proccess
	var intentionProccessService = services.NewIntentionProccess(intentionService, fieldDetectorService, actionService, logger)
	var multiIntentionProccessService = services.NewMultiIntentionProcessService(intentionProccessService, llmAdapter, logger)

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

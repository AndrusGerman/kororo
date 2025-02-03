package main

import (
	"bufio"
	"context"
	"fmt"
	"kororo/internal/adapters/config"
	"kororo/internal/adapters/llm/deepseek"
	"kororo/internal/adapters/llm/gemini"
	"kororo/internal/adapters/rest"
	mongodb "kororo/internal/adapters/storage/mongo"
	"kororo/internal/adapters/storage/mongo/repository"
	"kororo/internal/core/services"
	"os"
	"strings"
)

func main() {

	var err error
	var config = config.NewConfig()

	var mongo *mongodb.Mongo

	mongo, err = mongodb.NewMongo(config)
	if err != nil {
		panic(err)
	}

	var restAdapter = rest.New()

	llmAdapter, err := gemini.New(restAdapter, config)
	if err != nil {
		panic(err)
	}

	deepseek.New(restAdapter, config)

	//var a, b = llmAdapter.ProcessSystemMessage("Saluda a mi amigo", "hola")
	//log.Println("llmAdapter: ", a, b)
	//
	//return

	// intention
	var intentionRepository = repository.NewIntentionRepository(mongo)
	var targetDectector = services.NewTargetDectector(llmAdapter)
	var intentionService = services.NewIntentionService(intentionRepository, targetDectector, llmAdapter)

	// action
	var actionRepository = repository.NewActionRepository(mongo)
	var actionService = services.NewActionService(actionRepository, llmAdapter)

	// field detector
	var fieldDetectorService = services.NewFieldDetectorService(llmAdapter)

	var intentionProccessService = services.NewIntentionProccess(intentionService, fieldDetectorService, actionService)
	var multiIntentionProccessService = services.NewMultiIntentionProcessService(intentionProccessService, llmAdapter)

	// Intenciones:
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

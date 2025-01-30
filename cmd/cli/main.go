package main

import (
	"context"
	"errors"
	"fmt"
	"kororo/internal/adapters/config"
	"kororo/internal/adapters/llm/ollama"
	"kororo/internal/adapters/rest"
	mongodb "kororo/internal/adapters/storage/mongo"
	"kororo/internal/adapters/storage/mongo/repository"
	"kororo/internal/core/domain"
	"kororo/internal/core/domain/models"
	"kororo/internal/core/domain/types"
	"kororo/internal/core/ports"
	"kororo/internal/core/services"
	"log"
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
	var llmAdapter = ollama.New(restAdapter)

	var intentionRepository = repository.NewIntentionRepository(mongo)
	var intentionService = services.NewIntentionService(intentionRepository, llmAdapter)

	var actionRepository = repository.NewActionRepository(mongo)
	var actionService = services.NewActionService(actionRepository, llmAdapter)

	//return
	var fieldDetectorService = services.NewFieldDetectorService(llmAdapter)

	//processIntention("Detecta los latidos de mi amigo", intentionService, fieldDetectorService, actionService)

	processIntention("En mi lista de procesos, cual es el navegador que se esta ejecutando?", intentionService, fieldDetectorService, actionService)

	processIntention("En mi lista de procesos, cual aplicacion de contenedores se esta ejecutando?", intentionService, fieldDetectorService, actionService)

	//processIntention("Quiero saludar a mi amigo, el es manuel", intentionService, fieldDetectorService, actionService)
	//processIntention("Quiero saludar a Jhon", intentionService, fieldDetectorService, actionService)
	//processIntention("Quiero ordenar un libro harry potter", intentionService, fieldDetectorService, actionService)
	//processIntention("Voy a comprar un libro llamado el señor de los anillos", intentionService, fieldDetectorService, actionService)
	//processIntention("Quiero enviar 2$ a mi amigo luis", intentionService, fieldDetectorService, actionService)
	//
}

func processIntention(text string, intentionService ports.IntentionService, fieldDetectorService ports.FieldDetectorService, actionService ports.ActionService) {
	intention, err := intentionService.Detect(context.TODO(), text)
	if err != nil {
		panic(err)
	}

	var fields = make([]string, 0)
	for _, field := range intention.Fields {
		fields = append(fields, field.Description)
	}

	fmt.Printf("\n--La intencion es %s, los campos requeridos son [%s]\n", intention.Description, strings.Join(fields, ", "))

	fieldsValue, err := fieldDetectorService.DetectFields(intention, text)
	if err != nil {
		if errors.Is(err, domain.ErrFieldsRequired) {
			var llmError = new(models.LLMError)
			errors.As(err, &llmError)

			fmt.Println("SystemResponse: ", llmError.UserMessage)
			return
		}
		if errors.Is(err, domain.ErrFieldDescriptionNotFound) {
			fmt.Printf("Error interno llm: %s\n", err.Error())
			return
		}
		log.Println(err)
		return
	}

	for _, fieldValue := range fieldsValue {
		fmt.Printf("El campo '%s' tiene el valor '%s'\n", fieldValue.Field.Description, fieldValue.Value)
	}

	var actionPipelineContext = models.NewActionPipelineContext(fieldsValue)

	for _, actionId := range intention.Actions {
		action, err := actionService.GetAction(context.TODO(), types.Id(actionId))
		if err != nil {
			log.Println("Error al obtener la accion: ", err)
			return
		}
		fmt.Println("---Procesando la accion: ", action.Description)

		actionResponse, err := actionService.ProcessAction(context.TODO(), action, actionPipelineContext)
		if err != nil {
			log.Println("Error al procesar la accion: ", err)
			return
		}

		actionPipelineContext.AddActionResponse(actionResponse)

		if action.ActionProccessType != types.ActionProccessTypeCommand {
			fmt.Println("----Respuesta de la accion: ", actionResponse.Response)
		}
	}

	// respuesta final

	responseAction, err := actionService.GetAction(context.TODO(), types.Id(intention.ResponseAction))
	if err != nil {
		log.Println("Error al obtener la accion de respuesta: ", err)
		return
	}

	responseActionResponse, err := actionService.ProcessAction(context.TODO(), responseAction, actionPipelineContext)
	if err != nil {
		log.Println("Error al procesar la accion de respuesta: ", err)
		return
	}

	fmt.Println("Respuesta final: ", responseActionResponse.Response)
}

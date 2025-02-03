package services

import (
	"context"
	"errors"
	"kororo/internal/core/domain"
	"kororo/internal/core/domain/models"
	"kororo/internal/core/domain/types"
	"kororo/internal/core/ports"
)

type IntentionProccess struct {
	IntentionService     ports.IntentionService
	FieldDetectorService ports.FieldDetectorService
	ActionService        ports.ActionService
}

func (i *IntentionProccess) Process(ctx context.Context, text string) (string, error) {
	intention, err := i.IntentionService.Detect(context.TODO(), text)
	if err != nil {
		return "", err
	}

	var fields = make([]string, 0)
	for _, field := range intention.Fields {
		fields = append(fields, field.Description)
	}

	//fmt.Printf("\n--La intencion es %s, los campos requeridos son [%s]\n", intention.Description, strings.Join(fields, ", "))

	fieldsValue, err := i.FieldDetectorService.DetectFields(intention, text)

	if err != nil {
		if errors.Is(err, domain.ErrFieldsRequired) {
			var llmError = new(models.LLMError)
			errors.As(err, &llmError)
			return "Error al detectar los campos requeridos " + llmError.UserMessage, nil
		}
		if errors.Is(err, domain.ErrFieldDescriptionNotFound) {
			return "Error interno llm: " + err.Error(), nil
		}

		return "", err
	}

	//for _, fieldValue := range fieldsValue {
	//	fmt.Printf("El campo '%s' tiene el valor '%s'\n", fieldValue.Field.Description, fieldValue.Value)
	//}

	var actionPipelineContext = models.NewActionPipelineContext(fieldsValue, text)

	for _, actionId := range intention.Actions {
		action, err := i.ActionService.GetAction(context.TODO(), types.Id(actionId))
		if err != nil {
			return "Error interno: al obtener la accion: " + err.Error(), nil

		}
		//fmt.Println("---Procesando la accion: ", action.Description)

		actionResponse, err := i.ActionService.ProcessAction(context.TODO(), action, actionPipelineContext)
		if err != nil {
			return "Error interno: al procesar la accion: " + err.Error(), nil
		}

		actionPipelineContext.AddActionResponse(actionResponse)
	}

	responseAction, err := i.ActionService.GetAction(context.TODO(), types.Id(intention.ResponseAction))
	if err != nil {
		return "Error interno: al obtener la accion de respuesta: " + err.Error(), nil
	}

	responseActionResponse, err := i.ActionService.ProcessAction(context.TODO(), responseAction, actionPipelineContext)
	if err != nil {
		return "Error interno: al procesar la accion de respuesta: " + err.Error(), nil
	}

	return responseActionResponse.Response, nil
}

func NewIntentionProccess(intentionService ports.IntentionService, fieldDetectorService ports.FieldDetectorService, actionService ports.ActionService) ports.IntentionProccessService {
	return &IntentionProccess{
		IntentionService:     intentionService,
		FieldDetectorService: fieldDetectorService,
		ActionService:        actionService,
	}
}

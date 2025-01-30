package services

import (
	"context"
	"encoding/json"
	"fmt"
	"kororo/internal/core/domain"
	"kororo/internal/core/domain/models"
	"kororo/internal/core/domain/types"
	"kororo/internal/core/ports"
	"os/exec"
	"strings"
)

func NewActionService(actionRepository ports.ActionRepository, llmAdapter ports.LLMAdapter) ports.ActionService {
	return &actionService{actionRepository: actionRepository, llmAdapter: llmAdapter}
}

type actionService struct {
	actionRepository ports.ActionRepository
	llmAdapter       ports.LLMAdapter
}

func (s *actionService) GetAction(ctx context.Context, id types.Id) (*models.Action, error) {
	return s.actionRepository.GetById(ctx, id)
}

func (s *actionService) ProcessAction(ctx context.Context, action *models.Action, actionContext *models.ActionPipelineContext) (*models.ActionResponse, error) {

	if action.ActionProccessType == types.ActionProccessTypeLLMResponse {
		return s.processLLMResponse(ctx, action, actionContext)
	}

	if action.ActionProccessType == types.ActionProccessTypeCommand {
		return s.processCommand(ctx, action, actionContext)
	}

	return nil, domain.ErrActionTypeNotSupported
}

func (s *actionService) processCommand(_ context.Context, action *models.Action, actionContext *models.ActionPipelineContext) (*models.ActionResponse, error) {
	command := action.Command.Command

	cmd := exec.Command("cmd", "/c", command)

	var output, err = cmd.Output()
	if err != nil {
		return nil, err
	}

	return &models.ActionResponse{
		ActionId: types.Id(action.Id),
		Status:   "success",
		Response: string(output),
	}, nil
}

func (s *actionService) processLLMResponse(_ context.Context, action *models.Action, actionContext *models.ActionPipelineContext) (*models.ActionResponse, error) {

	var replacesValues []string

	for _, field := range action.Fields {
		var value, err = actionContext.GetField(field)
		if err != nil {
			fmt.Println("Error al obtener el valor del campo: ", field.Name)
			return nil, err
		}
		replacesValues = append(replacesValues, fmt.Sprintf("${%s}", field.Name), value)
	}

	if listActionResponse := actionContext.GetAllActionResponses(); len(listActionResponse) > 0 {
		var lastResponse = listActionResponse[len(listActionResponse)-1]
		replacesValues = append(replacesValues, "${beforePipeRaw}", lastResponse.Response)
	}

	replacer := strings.NewReplacer(replacesValues...)
	prompt := replacer.Replace(action.ProcessLLMSystemPrompt)

	var systemPrompt = `Eres una asistente de IA encargada de responder exactamente con la acciones que te pide el usuario en su prompt`

	if action.ResponseType == types.ActionResponseTypeJson {
		systemPrompt += `Se te pedira procesar una serie de datos, donde las respuestas deben ser en formato JSON con el siguiente formato: 
		[{
			"name": "nombre del campo",
			"value": "valor del campo"
		}]

		Ejemplos:

		Mensaje: extraer el 'id' de la respuesta  [{"nombre": "luis", "id": "123"}]
		Respuesta: [{"name":"id", "value":"123"}]

		Mensaje: genera el campo 'fruta' con el valor de una fruta aleatoria
		Respuesta: [{"name":"fruta", "value":"manzana"}]

		Tus mensajes deben ser un formato json valido, sin incluir formato markdown o html.
		`
	}

	response, err := s.llmAdapter.ProcessSystemMessage(systemPrompt, prompt)
	if err != nil {
		return nil, err
	}

	var actionResponse = &models.ActionResponse{
		ActionId: types.Id(action.Id),
		Status:   "success",
		Response: response,
	}

	if action.ResponseType == types.ActionResponseTypeJson {
		responseFields, err := s.processPipeline(context.TODO(), response)
		if err != nil {
			return nil, err
		}
		actionResponse.ResponseFields = responseFields
	}

	return actionResponse, nil
}

func (s *actionService) processPipeline(ctx context.Context, response string) ([]*models.ActionsResponseFields, error) {
	var responseFields []*models.ActionsResponseFields
	json.Unmarshal([]byte(response), &responseFields)
	return responseFields, nil
}

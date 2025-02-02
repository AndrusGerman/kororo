package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kororo/internal/core/domain"
	"kororo/internal/core/domain/models"
	"kororo/internal/core/domain/types"
	"kororo/internal/core/ports"
	"net/http"
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

	if action.ActionProccessType == types.ActionProccessTypeRequestHttp {
		return s.processHttp(ctx, action, actionContext)
	}

	return nil, domain.ErrActionTypeNotSupported
}

func (s *actionService) processHttp(_ context.Context, action *models.Action, actionContext *models.ActionPipelineContext) (*models.ActionResponse, error) {

	replacer, err := s.replacer(actionContext, action)
	if err != nil {
		return nil, err
	}
	var respHttp *http.Response
	var method = strings.ToUpper(action.Http.Method)
	var url = replacer.Replace(action.Http.Url)
	if method == "GET" {
		respHttp, err = http.Get(url)
	}

	if method == "POST" {
		action.Http.Body = replacer.Replace(action.Http.Body)
		//log.Println("Body: ", action.Http.Body)
		respHttp, err = http.Post(url, "application/json", strings.NewReader(action.Http.Body))
	}

	if err != nil {
		return nil, err
	}

	defer respHttp.Body.Close()

	body, err := io.ReadAll(respHttp.Body)
	if err != nil {
		return nil, err
	}

	var mapResponse = make(map[string]interface{})
	json.Unmarshal(body, &mapResponse)

	//log.Println("MapResponse: ", mapResponse)

	if action.Http.CheckLLMResponsePrompt != "" {

	}

	var actionResponseFields []*models.ActionsResponseFields
	actionResponseFields = append(actionResponseFields, &models.ActionsResponseFields{
		Name:  action.Http.HttpValueNameResponse,
		Value: string(body),
	})

	for _, formatHttpResponse := range action.Http.FormatHttpResponse {
		actionResponseFields = append(actionResponseFields, &models.ActionsResponseFields{
			Name:  formatHttpResponse.ValueName,
			Value: mapResponse[formatHttpResponse.Src].(string),
		})
	}

	var response = &models.ActionResponse{

		ActionId:       types.Id(action.Id),
		Status:         "success",
		Response:       string(body),
		ResponseFields: actionResponseFields,
	}

	return response, nil
}

func (s *actionService) processCommand(_ context.Context, action *models.Action, actionContext *models.ActionPipelineContext) (*models.ActionResponse, error) {
	replacer, err := s.replacer(actionContext, action)
	if err != nil {
		return nil, err
	}

	command := replacer.Replace(action.Command.Command)
	cmd := exec.Command("cmd", "/c", command)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return &models.ActionResponse{
		ActionId: types.Id(action.Id),
		Status:   "success",
		Response: string(output),
	}, nil
}

func (s *actionService) replacer(actionContext *models.ActionPipelineContext, action *models.Action) (*strings.Replacer, error) {
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

	replacesValues = append(replacesValues, "${userPrompt}", actionContext.GetUserPrompt())

	return strings.NewReplacer(replacesValues...), nil
}
func (s *actionService) processLLMResponse(_ context.Context, action *models.Action, actionContext *models.ActionPipelineContext) (*models.ActionResponse, error) {

	replacer, err := s.replacer(actionContext, action)
	if err != nil {
		return nil, err
	}

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

	response = strings.Replace(response, "```json", "", -1)
	response = strings.Replace(response, "```", "", -1)

	var err = json.Unmarshal([]byte(response), &responseFields)

	return responseFields, err
}

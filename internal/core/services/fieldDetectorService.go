package services

import (
	"encoding/json"
	"fmt"
	"kororo/internal/core/domain"
	"kororo/internal/core/domain/models"
	"kororo/internal/core/ports"
	"strings"
)

type FieldDetectorService struct {
	llmAdapter ports.LLMAdapter
	logger     ports.LogService
}

func NewFieldDetectorService(llmAdapter ports.LLMAdapter, logger ports.LogService) ports.FieldDetectorService {
	return &FieldDetectorService{
		llmAdapter: llmAdapter,
		logger:     logger,
	}
}

func (s *FieldDetectorService) DetectFields(intention *models.Intention, text string) ([]models.FieldValue, error) {
	if len(intention.Fields) == 0 {
		return nil, nil
	}

	type fieldJson struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Type        string `json:"type"`
	}

	type fieldValueJson struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	type IAResponse struct {
		SystemResponse string           `json:"systemResponse"`
		Fields         []fieldValueJson `json:"fields"`
	}

	var fields []fieldJson

	var fieldsNames []string
	for _, field := range intention.Fields {
		fieldsNames = append(fieldsNames, fmt.Sprintf("'%s'", field.Name))
		fields = append(fields, fieldJson{
			Name:        field.Name,
			Description: field.Description,
			Type:        string(field.Type),
		})
	}

	fieldsJson, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}

	var systemMessage = `Eres una IA que detecta los campos requeridos de una intención de un usuario, solo respondes con JSON validos sin
	etiquetas de formato markdown.
	Te enviare una conversación entre un usuario y un asistente de IA y debes detectar los campos requeridos que están en la conversación.
	La respuesta tiene que ser un JSON valido.
	Solo tienes dos maneras de responder:
	1. Si falta algún campo requerido, debes responder con el campo requerido que faltan en un JSON con el siguiente formato:
	{	
		"systemResponse":"Por favor me puedes indicar el campo 'nombre del campo requerido'",
		"fields":[]
	}
	2. Si todos los campos requeridos están presentes, debes responder con el siguiente JSON:
	{
		"systemResponse":"",
		"fields":[{"name": "nombre del campo requerido", "value": "valor del campo requerido"}]
	}

	Tu respuesta jamas debe salir del formato JSON correspondiente y solo tiene que extraer los campos requeridos de cualquier mensaje 
	que envie el usuario.
		
	Esto es una respuesta valida: 
	{
		"systemResponse":"",
		"fields":[{"name":"username", "value":"julio"}]
	}
	


	Esto es una respuesta valida: 
	{
		"systemResponse":"Falta el campo requerido 'nombre de usuario'",
		"fields":[]
	}
	

	Ejemplo con los siguientes datos:
	[{"name":"edad", "description":"edad del usuario", "type":"number"}]
	
	Mensaje: "Hola, me llamo juan, tengo 20 años"
	
	La respuesta debe ser:
	{"systemResponse":"", "fields":[{"name":"edad", "value":"20"}]}`

	var userMessage = "Datos: " + string(fieldsJson) + "\n" + "Conversación: " + text

	response, err := s.llmAdapter.ProcessSystemMessage(systemMessage, userMessage)
	if err != nil {
		return nil, err
	}

	response, err = domain.JSONClear(response)
	if err != nil {
		s.logger.Error("FieldDetectorService.DetectFields.ErrJSONClear", response)
		return nil, err
	}

	var iaResponse IAResponse
	err = json.Unmarshal([]byte(response), &iaResponse)
	if err != nil {
		fmt.Println("iaResponseRaw: ", response)
		return nil, err
	}

	if iaResponse.SystemResponse != "" {
		return nil, fmt.Errorf("%w: %w", domain.ErrFieldsRequired, models.NewLLMError(iaResponse.SystemResponse, ""))
	}

	var fieldsValues []models.FieldValue
	for _, fieldJson := range iaResponse.Fields {
		field := s.getFieldByName(intention.Fields, fieldJson.Name)
		if field == nil {
			return nil, fmt.Errorf("%w: %s", domain.ErrFieldDescriptionNotFound, fieldJson.Name)
		}
		fieldsValues = append(fieldsValues, models.FieldValue{
			Field: field,
			Value: fieldJson.Value,
		})
	}

	return fieldsValues, nil
}

func (s *FieldDetectorService) getFieldByName(fields []*models.Field, name string) *models.Field {
	name = strings.ToLower(strings.TrimSpace(name))
	for _, field := range fields {
		var fieldName = strings.ToLower(strings.TrimSpace(field.Name))
		if fieldName == name {
			return field
		}
	}
	return nil
}

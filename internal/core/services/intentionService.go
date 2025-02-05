package services

import (
	"context"
	"encoding/json"
	"fmt"
	"kororo/internal/core/domain"
	"kororo/internal/core/domain/models"
	"kororo/internal/core/ports"

	"github.com/AndrusGerman/go-criteria"
)

func NewIntentionService(intentionRepository ports.IntentionRepository, targetDectector ports.TargetDectector, llmAdapter ports.LLMAdapter) ports.IntentionService {
	return &intentionService{intentionRepository: intentionRepository, targetDectector: targetDectector, llmAdapter: llmAdapter}
}

type intentionService struct {
	intentionRepository ports.IntentionRepository
	llmAdapter          ports.LLMAdapter
	targetDectector     ports.TargetDectector
}

func (s *intentionService) Detect(ctx context.Context, text string) (*models.Intention, error) {
	intentions, err := s.intentionRepository.Search(ctx, criteria.EmptyCriteria())
	if err != nil {
		return nil, err
	}

	type IntentionJson struct {
		Description string `json:"description"`
		IntentIndex int    `json:"intent_index"`
	}

	var intentionsJson = make([]*IntentionJson, 0)

	for index, intention := range intentions {
		intentionsJson = append(intentionsJson, &IntentionJson{Description: intention.Description, IntentIndex: index})
	}

	jsonString, err := json.Marshal(intentionsJson)
	if err != nil {
		return nil, err
	}

	type ResponseJson struct {
		UserMessage string `json:"user_message"`
		IntentIndex int    `json:"intent_index"`
	}

	var systemMessage = `Eres un asistente de IA que detecta la intención de un usuario y responde con el índice de la intención

	que corresponde al mensaje del usuario. Si el usuario dice algo que no tiene que ver con la intención, responde con -1. Las intenciones son.
	Si el usuario escribe mas de una intencion, escribe el parametro "user_message" con el mensaje de usuario 
	"solo puede mandar al sistema una intención a la vez, 'intencion1', 'intencion2' detectada" y "intent_index" con -1

	Ejemplos de respuestas validas:
	Mensaje: Quiero saludar a alguien
	Intenciones: [{"description": "Saludar a mi amigo", "intent_index": 0}, {"description": "Abrir una cuenta bancaria", "intent_index": 1}]
	Respuesta: {"intent_index": 0}

	Mensaje: Quiero abrir una cuenta bancaria
	Intenciones: [{"description": "Saludar a mi amigo", "intent_index": 0}, {"description": "Abrir una cuenta bancaria", "intent_index": 1}]
	Respuesta: {"intent_index": 1}


	Mensaje: Saluda a mi amigo
	Intenciones: [{"description": "Saludar a mi amigo", "intent_index": 0}, {"description": "Abrir una cuenta bancaria", "intent_index": 1}]
	Respuesta: {"intent_index": 0}
	
	Mensaje: Saluda a mi amigo y abrir una cuenta bancaria
	Intenciones: [{"description": "Saludar a mi amigo", "intent_index": 0}, {"description": "Abrir una cuenta bancaria", "intent_index": 1}]
	Respuesta: {"user_message": "solo puede mandar al sistema una intención a la vez, 'Saludar a mi amigo', 'Abrir una cuenta bancaria' detectada", "intent_index": -1}

	Mensaje: Quiero vender una casa
	Intenciones: [{"description": "Saludar a mi amigo", "intent_index": 0}, {"description": "Abrir una cuenta bancaria", "intent_index": 1}]
	Respuesta: {"intent_index": -1}


	Tus respuestas deben ser solo el json valido, {
		"user_message": "mensaje de usuario",
		"intent_index": "índice de la intención"
	}
	y solo puedes mandar uno de los parametros, no ambos.
	
	Estas son las intenciones: ` + string(jsonString)

	target, err := s.targetDectector.Detect(ctx, text)

	if err != nil {
		return nil, err
	}

	response, err := s.llmAdapter.ProcessSystemMessage(systemMessage, target)
	if err != nil {
		return nil, err
	}

	response = domain.MustJSONClear(response)

	var responseJson ResponseJson
	if err := json.Unmarshal([]byte(response), &responseJson); err != nil {
		return nil, err
	}

	if responseJson.IntentIndex == -1 {
		if responseJson.UserMessage != "" {
			return nil, fmt.Errorf("%w: %w", domain.ErrMultipleIntentionSend, models.NewLLMError(responseJson.UserMessage, responseJson.UserMessage))
		}
		return nil, domain.ErrIntentionNotFound
	}

	return intentions[responseJson.IntentIndex], nil
}

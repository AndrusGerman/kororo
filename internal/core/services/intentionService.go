package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"kororo/internal/core/domain/models"
	"kororo/internal/core/ports"

	"github.com/AndrusGerman/go-criteria"
)

func NewIntentionService(intentionRepository ports.IntentionRepository, llmAdapter ports.LLMAdapter) ports.IntentionService {
	return &intentionService{intentionRepository: intentionRepository, llmAdapter: llmAdapter}
}

type intentionService struct {
	intentionRepository ports.IntentionRepository
	llmAdapter          ports.LLMAdapter
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

	var systemMessage = `Eres un asistente de IA que detecta la intención de un usuario y responde con el índice de la intención
	que corresponde al mensaje del usuario. Si el usuario dice algo que no tiene que ver con la intención, responde con -1. Las intenciones son

	Ejemplos de respuestas validas:
	Mensaje: Quiero saludar a alguien
	Intenciones: [{"description": "Saludar a mi amigo", "intent_index": 0}, {"description": "Abrir una cuenta bancaria", "intent_index": 1}]
	Respuesta: 0

	Mensaje: Quiero abrir una cuenta bancaria
	Intenciones: [{"description": "Saludar a mi amigo", "intent_index": 0}, {"description": "Abrir una cuenta bancaria", "intent_index": 1}]
	Respuesta: 1

	Mensaje: Quiero vender una casa
	Intenciones: [{"description": "Saludar a mi amigo", "intent_index": 0}, {"description": "Abrir una cuenta bancaria", "intent_index": 1}]
	Respuesta: -1


	Tus respuestas deben ser solo el índice de la intención, no debes agregar ningún otro texto.
	las siguientes: ` + string(jsonString)
	response, err := s.llmAdapter.ProcessSystemMessage(systemMessage, text)
	if err != nil {
		return nil, err
	}

	var intentIndex int
	if _, err := fmt.Sscanf(response, "%d", &intentIndex); err != nil {
		return nil, err
	}

	if intentIndex == -1 {
		return nil, errors.New("no se pudo detectar la intención")
	}

	return intentions[intentIndex], nil
}

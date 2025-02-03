package services

import (
	"context"
	"kororo/internal/core/ports"
)

func NewTargetDectector(llmAdapter ports.LLMAdapter) ports.TargetDectector {
	return &TargetDectector{llmAdapter: llmAdapter}
}

type TargetDectector struct {
	llmAdapter ports.LLMAdapter
}

func (s *TargetDectector) Detect(ctx context.Context, text string) (string, error) {

	var systemMessage = `Eres un asistente de IA que detecta el objectivo de un usuario a partir de un mensaje.
	El mensaje solo debe ser el objetivo del usuario, no debes añadir ningun comentario.
	
	Ejemplos:
	
	Usuario: Quiero comprar un coche
	Respuesta: comprar un coche
	
	Usuario: Quiero alquilar un coche
	Respuesta: alquilar un coche
	
	Usuario: Envia un mensaje a luis, que es el jefe de la empresa
	Respuesta: enviar un mensaje
	
	Usuario: Verifica cual es el balance de mi cuenta bancaria
	Respuesta: verificar balance
	
	Usuario: Envia un audio a maria, 'Cual es el balance de mi cuenta bancaria?'
	Respuesta: Enviar un audio`

	return s.llmAdapter.ProcessSystemMessage(systemMessage, text)

}

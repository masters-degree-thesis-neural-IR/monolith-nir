package usecases

import "monolith-nir/service/application/domain"

type DocumentUc interface {
	Create(id string, title string, body string) error
	LoadDocuments(scoreResult []domain.ScoreResult) ([]domain.DocumentResult, error)
}

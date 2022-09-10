package usecases

import "monolith-nir/service/application/domain"

type SearchUc interface {
	LexicalSearchDocument(query string) ([]domain.ScoreResult, error)
	SemanticSearchDocument(query string) ([]domain.ScoreResult, error)
}

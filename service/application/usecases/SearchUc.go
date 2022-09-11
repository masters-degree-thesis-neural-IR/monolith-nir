package usecases

import "monolith-nir/service/application/domain"

type SearchUc interface {
	LexicalSearch(query string) ([]domain.ScoreResult, error)
	SemanticSearch(query string) ([]domain.ScoreResult, error)
}

package usecases

import "monolith-nir/service/application/domain"

type SearchUc interface {
	SearchDocument(query string) ([]domain.QueryResult, error)
	//MakeInvertedIndex(localQuery []string) (domain.InvertedIndex, error)
}

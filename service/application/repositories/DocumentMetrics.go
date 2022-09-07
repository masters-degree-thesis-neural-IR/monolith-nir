package repositories

import "monolith-nir/service/application/domain"

type DocumentMetricsRepository interface {
	FindDocuments(documentIDs map[string]int8) ([]domain.NormalizedDocument, error)
	Save(domain.NormalizedDocument) error
}

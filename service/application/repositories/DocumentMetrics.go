package repositories

import "monolith-nir/service/application/domain"

type DocumentMetricsRepository interface {
	FindByDocumentIDs(documentIDs map[string]int8) ([]domain.NormalizedDocument, error)
	Save(domain.NormalizedDocument) error
}

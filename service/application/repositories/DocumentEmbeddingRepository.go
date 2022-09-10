package repositories

import "monolith-nir/service/application/domain"

type DocumentEmbeddingRepository interface {
	Save(documentEmbedding domain.DocumentEmbedding) error
	FindByDocumentIDs(documentIDs map[string]int8) ([]domain.DocumentEmbedding, error)
}

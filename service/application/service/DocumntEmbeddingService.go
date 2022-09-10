package service

import (
	"monolith-nir/service/application/domain"
	"monolith-nir/service/application/logger"
	"monolith-nir/service/application/ports"
	"monolith-nir/service/application/repositories"
	"monolith-nir/service/application/usecases"
)

type WordEmbeddingService struct {
	Ch                          chan domain.Document
	Logger                      logger.Logger
	WordEmbedding               ports.WordEmbedding
	DocumentEmbeddingRepository repositories.DocumentEmbeddingRepository
}

func NewWordEmbeddingService(logger logger.Logger, documentEmbeddingRepository repositories.DocumentEmbeddingRepository, wordEmbedding ports.WordEmbedding) usecases.WordEmbeddingUc {

	ch := make(chan domain.Document, 100)

	w := WordEmbeddingService{
		Ch:                          ch,
		Logger:                      logger,
		WordEmbedding:               wordEmbedding,
		DocumentEmbeddingRepository: documentEmbeddingRepository,
	}

	go w.CreateDocumentEmbedding(ch)

	return w
}

func (w WordEmbeddingService) CreateDocumentEmbedding(documents chan domain.Document) {

	for document := range documents {
		embedding := w.WordEmbedding.Generate(document.Body)
		documentEmbedding := domain.DocumentEmbedding{
			Id:        document.Id,
			Embedding: embedding,
		}
		err := w.DocumentEmbeddingRepository.Save(documentEmbedding)

		if err != nil {
			w.Logger.Error(err.Error())
		}
	}

}

func (w WordEmbeddingService) CreateEmbedding(id, title, body string) {

	w.Ch <- domain.Document{
		Id:    id,
		Title: title,
		Body:  body,
	}

}

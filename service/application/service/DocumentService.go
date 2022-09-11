package service

import (
	"monolith-nir/service/application/domain"
	"monolith-nir/service/application/exception"
	"monolith-nir/service/application/logger"
	"monolith-nir/service/application/ports"
	"monolith-nir/service/application/repositories"
	"monolith-nir/service/application/usecases"
)

type DocumentService struct {
	Ch                 chan domain.Document
	Logger             logger.Logger
	DocumentEvent      ports.DocumentEvent
	DocumentRepository repositories.DocumentRepository
}

func NewDocumentService(logger logger.Logger, documentEvent ports.DocumentEvent, documentRepository repositories.DocumentRepository) usecases.DocumentUc {

	ch := make(chan domain.Document, 100)

	d := DocumentService{
		Ch:                 ch,
		Logger:             logger,
		DocumentEvent:      documentEvent,
		DocumentRepository: documentRepository,
	}
	go d.SaveDocument(ch)
	return d
}

func (s DocumentService) SaveDocument(documents chan domain.Document) {

	for document := range documents {
		var err = s.DocumentRepository.Save(document)

		if err != nil {
			s.Logger.Error(err.Error())
		}
	}

}

func (s DocumentService) Create(id string, title string, body string) error {

	if id == "" {
		return exception.ThrowValidationError("Invalid id from document")
	}

	s.Ch <- domain.Document{
		Id:    id,
		Title: title,
		Body:  body,
	}

	return nil

}

func (s DocumentService) LoadDocuments(scoreResult []domain.ScoreResult) ([]domain.DocumentResult, error) {

	var documentIDs = make([]string, len(scoreResult))
	for i, result := range scoreResult {
		documentIDs[i] = result.DocumentID
	}

	documents, err := s.DocumentRepository.FindByDocumentIDs(documentIDs)

	if err != nil {
		return []domain.DocumentResult{}, err
	}

	documentResults := make([]domain.DocumentResult, len(scoreResult))

	for i, result := range scoreResult {

		document := documents[result.DocumentID]

		documentResult := domain.DocumentResult{
			Similarity: result.Similarity,
			Document:   document,
		}
		documentResults[i] = documentResult

	}

	return documentResults, nil
}

package service

import (
	"monolith-nir/service/application/exception"
	"monolith-nir/service/application/ports"
	"monolith-nir/service/application/repositories"
	"monolith-nir/service/application/usecases"
)

type DocumentService struct {
	DocumentEvent      ports.DocumentEvent
	DocumentRepository repositories.DocumentRepository
}

func NewDocumentService(documentEvent ports.DocumentEvent, documentRepository repositories.DocumentRepository) usecases.DocumentUc {
	var c usecases.DocumentUc = &DocumentService{
		DocumentEvent:      documentEvent,
		DocumentRepository: documentRepository,
	}
	return c
}

func (s *DocumentService) Create(id string, title string, body string) error {

	if id == "" {
		return *exception.ThrowValidationError("Invalid id from document")
	}

	//if title == "" {
	//	return *exception.ThrowValidationError("Invalid title from document")
	//}

	//if body == "" {
	//	return *exception.ThrowValidationError("Invalid body from document")
	//}
	//
	//document := domain.Document{
	//	Id:    id,
	//	Title: title,
	//	Body:  body,
	//}
	//
	//var err = s.DocumentRepository.Save(document)
	//
	//if err != nil {
	//	return err
	//}

	//return err //s.DocumentEvent.Created(document)

	return nil

}

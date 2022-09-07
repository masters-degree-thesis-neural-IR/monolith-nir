package controller

import (
	"monolith-nir/service/application/domain"
	"monolith-nir/service/application/usecases"
	"monolith-nir/service/infraestructure/dto"
)

type Controller struct {
	DocumentService usecases.DocumentUc
	IndexService    usecases.CreateIndexUc
	Search          usecases.SearchUc
}

func NewController(documentService usecases.DocumentUc, indexService usecases.CreateIndexUc, search usecases.SearchUc) Controller {
	return Controller{
		DocumentService: documentService,
		IndexService:    indexService,
		Search:          search,
	}
}

func (c *Controller) SearchDocuments(query string) ([]domain.QueryResult, error) {
	return c.Search.SearchDocument(query)
}

func (c *Controller) CreateDocument(document dto.Document) error {

	err := c.DocumentService.Create(document.Id, document.Title, document.Body)

	if err != nil {
		return err
	}

	err = c.IndexService.CreateIndex(document.Id, document.Title, document.Body)

	if err != nil {
		return err
	}

	return nil

}

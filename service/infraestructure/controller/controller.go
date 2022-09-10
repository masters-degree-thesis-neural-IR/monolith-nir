package controller

import (
	"monolith-nir/service/application/domain"
	"monolith-nir/service/application/usecases"
	"monolith-nir/service/infraestructure/dto"
)

type Controller struct {
	WordEmbeddingUc usecases.WordEmbeddingUc
	DocumentService usecases.DocumentUc
	IndexService    usecases.CreateIndexUc
	Search          usecases.SearchUc
}

func NewController(wordEmbeddingUc usecases.WordEmbeddingUc, documentService usecases.DocumentUc, indexService usecases.CreateIndexUc, search usecases.SearchUc) Controller {
	return Controller{
		WordEmbeddingUc: wordEmbeddingUc,
		DocumentService: documentService,
		IndexService:    indexService,
		Search:          search,
	}
}

func (c *Controller) SearchDocuments(query string) ([]domain.ScoreResult, error) {
	//c.Search.LexicalSearchDocument(query)
	return c.Search.SemanticSearchDocument(query)
}

func (c *Controller) CreateDocument(document dto.Document) error {

	//go func() {
	//	err := c.DocumentService.Create(document.Id, document.Title, document.Body)
	//	if err != nil {
	//		log.Fatalln(err.Error())
	//	}
	//}()

	go func() {
		c.WordEmbeddingUc.CreateEmbedding(document.Id, document.Title, document.Body)
	}()

	err := c.IndexService.CreateIndex(document.Id, document.Title, document.Body)

	if err != nil {
		return err
	}

	return nil

}

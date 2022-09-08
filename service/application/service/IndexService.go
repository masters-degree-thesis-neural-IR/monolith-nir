package service

import (
	"monolith-nir/service/application/domain"
	"monolith-nir/service/application/nlp"
	"monolith-nir/service/application/repositories"
	"monolith-nir/service/application/usecases"
)

type IndexService struct {
	Ch                        chan domain.NormalizedDocument
	IndexMemoryRepository     repositories.IndexMemoryRepository
	DocumentMetricsRepository repositories.DocumentMetricsRepository
}

func NewIndexService(documentMetricsRepository repositories.DocumentMetricsRepository, indexMemoryRepository repositories.IndexMemoryRepository) usecases.CreateIndexUc {

	ch := make(chan domain.NormalizedDocument, 100)

	var c usecases.CreateIndexUc = IndexService{
		Ch:                        ch,
		IndexMemoryRepository:     indexMemoryRepository,
		DocumentMetricsRepository: documentMetricsRepository,
	}
	go c.CreateNormalizedDocument(ch)
	return c
}

func (i IndexService) CreateNormalizedDocument(ch chan domain.NormalizedDocument) {

	for document := range ch {
		i.DocumentMetricsRepository.Save(document)
	}

}

func (i IndexService) CreateIndex(id string, title string, body string) error {

	tokens := nlp.Tokenizer(body, true)
	normalizedTokens, err := nlp.RemoveStopWords(tokens, "en")

	if err != nil {
		return err
	}

	i.Ch <- domain.NormalizedDocument{
		Id:     id,
		Length: len(normalizedTokens),
		Tf:     nlp.TermFrequency(normalizedTokens),
	}

	for _, term := range normalizedTokens {

		var documentList, err = i.IndexMemoryRepository.FindByTerm(term)

		if err != nil {
			return err
		}

		if documentList != nil && len(documentList) > 0 {
			if nlp.NotContains(id, documentList) {
				documentList = append(documentList, id)
				i.IndexMemoryRepository.Update(term, documentList)
			}
		} else {
			i.IndexMemoryRepository.Save(term, []string{id})
		}
	}

	return nil
}

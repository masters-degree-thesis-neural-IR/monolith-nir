package service

import (
	"monolith-nir/service/application/domain"
	"monolith-nir/service/application/logger"
	"monolith-nir/service/application/nlp"
	"monolith-nir/service/application/repositories"
	"monolith-nir/service/application/usecases"
)

type IndexService struct {
	Logger                    logger.Logger
	Ch                        chan domain.Document
	IndexMemoryRepository     repositories.IndexMemoryRepository
	DocumentMetricsRepository repositories.DocumentMetricsRepository
}

func NewIndexService(logger logger.Logger, documentMetricsRepository repositories.DocumentMetricsRepository, indexMemoryRepository repositories.IndexMemoryRepository) usecases.CreateIndexUc {

	ch := make(chan domain.Document, 100)

	i := IndexService{
		Ch:                        ch,
		Logger:                    logger,
		IndexMemoryRepository:     indexMemoryRepository,
		DocumentMetricsRepository: documentMetricsRepository,
	}
	go i.CreateNormalizedDocument(ch)
	return i
}

func (i IndexService) CreateNormalizedDocument(ch chan domain.Document) {

	for document := range ch {

		tokens := nlp.Tokenizer(document.Body, true)
		normalizedTokens, err := nlp.RemoveStopWords(tokens, "en")

		if err != nil {
			i.Logger.Error(err.Error())
		}

		normalizedDocument := domain.NormalizedDocument{
			Id:     document.Id,
			Length: len(normalizedTokens),
			Tf:     nlp.TermFrequency(normalizedTokens),
		}

		i.DocumentMetricsRepository.Save(normalizedDocument)
	}

}

func (i IndexService) CreateIndex(id string, title string, body string) error {

	tokens := nlp.Tokenizer(body, true)
	normalizedTokens, err := nlp.RemoveStopWords(tokens, "en")

	if err != nil {
		i.Logger.Error(err.Error())
		return err
	}

	i.Ch <- domain.Document{
		Id:    id,
		Title: title,
		Body:  body,
	}

	for _, term := range normalizedTokens {

		var documentList, err = i.IndexMemoryRepository.FindByTerm(term)

		if err != nil {
			i.Logger.Error(err.Error())
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

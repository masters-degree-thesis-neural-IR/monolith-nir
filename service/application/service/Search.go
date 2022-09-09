package service

import (
	"monolith-nir/service/application/domain"
	"monolith-nir/service/application/exception"
	"monolith-nir/service/application/logger"
	"monolith-nir/service/application/nlp"
	"monolith-nir/service/application/repositories"
	"monolith-nir/service/application/usecases"
)

type Search struct {
	Logger                    logger.Logger
	DocumentRepository        repositories.DocumentRepository
	IndexMemoryRepository     repositories.IndexMemoryRepository
	DocumentMetricsRepository repositories.DocumentMetricsRepository
}

func NewSearch(logger logger.Logger, documentMetricsRepository repositories.DocumentMetricsRepository, indexMemoryRepository repositories.IndexMemoryRepository, documentRepository repositories.DocumentRepository) usecases.SearchUc {
	return Search{
		Logger:                    logger,
		DocumentRepository:        documentRepository,
		IndexMemoryRepository:     indexMemoryRepository,
		DocumentMetricsRepository: documentMetricsRepository,
	}
}

func (s Search) MakeInvertedIndex(localQuery []string, foundDocuments map[string]int8) (domain.InvertedIndex, error) {

	normalizedDocuments, err := s.DocumentMetricsRepository.FindDocuments(foundDocuments)

	invertedIndex := domain.InvertedIndex{
		Df:                      map[string]int{},
		NormalizedDocumentFound: make(map[string]domain.NormalizedDocument, 0),
		CorpusSize:              0,
	}

	if err != nil {
		s.Logger.Error(err.Error())
		return domain.InvertedIndex{}, err
	}

	for _, term := range localQuery {
		for _, document := range normalizedDocuments {
			qtd := document.Tf[term]
			//TODO: Remover a condição
			if qtd > 0 {
				invertedIndex.Df[term] += 1
			}
		}
	}

	for _, document := range normalizedDocuments {
		invertedIndex.NormalizedDocumentFound[document.Id] = document
	}

	invertedIndex.CorpusSize = len(invertedIndex.NormalizedDocumentFound)
	invertedIndex.Idf = nlp.CalcIdf(invertedIndex.Df, invertedIndex.CorpusSize)

	return invertedIndex, nil

}

func (s Search) FindDocuments(localQuery []string) map[string]int8 {

	var foundDocuments = make(map[string]int8, 0)

	for _, term := range localQuery {
		documents, _ := s.IndexMemoryRepository.FindByTerm(term)
		for _, docID := range documents {
			foundDocuments[docID] = 0
		}
	}

	return foundDocuments
}

func (s Search) SearchDocument(query string) ([]domain.QueryResult, error) {

	localQuery := nlp.Tokenizer(query, true) //nlp.RemoveStopWords(nlp.Tokenizer(query, true), "en")
	foundDocuments := s.FindDocuments(localQuery)
	invertedIndex, err := s.MakeInvertedIndex(localQuery, foundDocuments)

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, exception.ThrowValidationError(err.Error())
	}

	queryResults := nlp.SortDesc(nlp.ScoreBM25(localQuery, &invertedIndex), 10)

	return queryResults, nil
}

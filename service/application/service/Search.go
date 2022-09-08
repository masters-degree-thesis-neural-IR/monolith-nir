package service

import (
	"log"
	"monolith-nir/service/application/domain"
	"monolith-nir/service/application/exception"
	"monolith-nir/service/application/nlp"
	"monolith-nir/service/application/repositories"
	"monolith-nir/service/application/usecases"
)

type Search struct {
	DocumentRepository        repositories.DocumentRepository
	IndexMemoryRepository     repositories.IndexMemoryRepository
	DocumentMetricsRepository repositories.DocumentMetricsRepository
}

func NewSearch(documentMetricsRepository repositories.DocumentMetricsRepository, indexMemoryRepository repositories.IndexMemoryRepository, documentRepository repositories.DocumentRepository) usecases.SearchUc {
	return Search{
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
		log.Fatalln("Error....: ", err)
		return domain.InvertedIndex{}, err
	}

	for _, term := range localQuery {
		for _, document := range normalizedDocuments {
			qtd := document.Tf[term]
			invertedIndex.Df[term] += qtd
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

	localQuery := nlp.Tokenizer(query, true) //:= nlp.RemoveStopWords(nlp.Tokenizer(query, true), "en")
	foundDocuments := s.FindDocuments(localQuery)
	invertedIndex, err := s.MakeInvertedIndex(localQuery, foundDocuments)

	if err != nil {
		return nil, *exception.ThrowValidationError(err.Error())
	}

	queryResults := nlp.SortDesc(nlp.ScoreBM25(localQuery, &invertedIndex), 10)
	//tempQueryResults := make([]domain.QueryResult, len(queryResults))
	//
	//for i, queryResult := range queryResults {
	//
	//	doc, err := s.DocumentRepository.FindById(queryResult.NormalizedDocument.Id)
	//
	//	if doc == nil {
	//		continue
	//	}
	//
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	tempQueryResults[i].Similarity = queryResult.Similarity
	//	tempQueryResults[i].NormalizedDocument = queryResult.NormalizedDocument
	//	tempQueryResults[i].Document = *doc
	//
	//}

	return queryResults, nil
}

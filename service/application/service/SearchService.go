package service

import (
	"monolith-nir/service/application/domain"
	"monolith-nir/service/application/exception"
	"monolith-nir/service/application/logger"
	"monolith-nir/service/application/nlp"
	"monolith-nir/service/application/ports"
	"monolith-nir/service/application/repositories"
	"monolith-nir/service/application/usecases"
)

type Search struct {
	Logger                      logger.Logger
	WordEmbedding               ports.WordEmbedding
	DocumentRepository          repositories.DocumentRepository
	IndexMemoryRepository       repositories.IndexMemoryRepository
	DocumentMetricsRepository   repositories.DocumentMetricsRepository
	DocumentEmbeddingRepository repositories.DocumentEmbeddingRepository
}

func NewSearchService(wordEmbedding ports.WordEmbedding, logger logger.Logger, documentEmbeddingRepository repositories.DocumentEmbeddingRepository, documentMetricsRepository repositories.DocumentMetricsRepository, indexMemoryRepository repositories.IndexMemoryRepository, documentRepository repositories.DocumentRepository) usecases.SearchUc {
	return Search{
		Logger:                      logger,
		WordEmbedding:               wordEmbedding,
		DocumentRepository:          documentRepository,
		IndexMemoryRepository:       indexMemoryRepository,
		DocumentMetricsRepository:   documentMetricsRepository,
		DocumentEmbeddingRepository: documentEmbeddingRepository,
	}
}

func (s Search) LoadDocumentsEmbedding(foundDocuments map[string]int8) ([]domain.DocumentEmbedding, error) {

	documentsEmbedding, err := s.DocumentEmbeddingRepository.FindByDocumentIDs(foundDocuments)

	if err != nil {
		s.Logger.Error(err.Error())
		return []domain.DocumentEmbedding{}, err
	}

	return documentsEmbedding, nil

}

func (s Search) MakeInvertedIndex(localQuery []string, foundDocuments map[string]int8) (domain.InvertedIndex, error) {

	normalizedDocuments, err := s.DocumentMetricsRepository.FindByDocumentIDs(foundDocuments)

	if err != nil {
		s.Logger.Error(err.Error())
		return domain.InvertedIndex{}, err
	}

	invertedIndex := domain.InvertedIndex{
		Df:                      map[string]int{},
		NormalizedDocumentFound: make(map[string]domain.NormalizedDocument, 0),
		CorpusSize:              0,
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

func (s Search) LexicalSearch(query string) ([]domain.ScoreResult, error) {

	localQuery := nlp.Tokenizer(query, true) //nlp.RemoveStopWords(nlp.Tokenizer(query, true), "en")
	foundDocuments := s.FindDocuments(localQuery)
	invertedIndex, err := s.MakeInvertedIndex(localQuery, foundDocuments)

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, exception.ThrowValidationError(err.Error())
	}

	scoreResult := nlp.SortDesc(nlp.ScoreBM25(localQuery, &invertedIndex), 10)

	return scoreResult, nil
}

func (s Search) SemanticSearch(query string) ([]domain.ScoreResult, error) {

	localQuery := nlp.Tokenizer(query, true)
	foundDocuments := s.FindDocuments(localQuery)
	documentsEmbedding, err := s.LoadDocumentsEmbedding(foundDocuments)

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, exception.ThrowValidationError(err.Error())
	}

	localQueryEmbedding := s.WordEmbedding.Generate(query)
	scoreResult := nlp.SortDesc(nlp.ScoreCosineSimilarity(localQueryEmbedding, documentsEmbedding), 10)

	return scoreResult, nil

}

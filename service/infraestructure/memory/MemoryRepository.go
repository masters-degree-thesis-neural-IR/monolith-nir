package memory

import (
	"monolith-nir/service/application/exception"
	"monolith-nir/service/application/logger"
	"monolith-nir/service/application/repositories"
	"monolith-nir/service/infraestructure/files"
	"time"
)

type IndexRepository struct {
	Chan  chan map[string][]string
	Index map[string][]string
}

func NewMemoryIndexRepository(logger logger.Logger) repositories.IndexMemoryRepository {

	ch := make(chan map[string][]string, 100)

	start := time.Now()
	index := files.LoadIndex()
	elapsed := time.Since(start)
	logger.Info("LOAD INDEX", "Total..:", len(index), "Time ..: ", elapsed)

	i := IndexRepository{
		Chan:  ch,
		Index: index,
	}

	go files.DumpIndex(ch)

	return &i
}

func (i *IndexRepository) FindByTerm(term string) ([]string, error) {
	documents := i.Index[term]
	return documents, nil
}

func (i *IndexRepository) Update(term string, documents []string) error {
	localDocuments := i.Index[term]
	if len(localDocuments) > 0 {
		i.Index[term] = documents
		i.Chan <- map[string][]string{term: documents}
	} else {
		return exception.ThrowValidationError("Not have term indexed. Use Save function")
	}

	return nil
}

func (i *IndexRepository) Save(term string, documents []string) error {
	localDocuments := i.Index[term]
	if len(localDocuments) > 0 {
		return exception.ThrowValidationError("Term already exists. Use Update function")
	}

	i.Index[term] = documents
	i.Chan <- map[string][]string{term: documents}
	return nil
}

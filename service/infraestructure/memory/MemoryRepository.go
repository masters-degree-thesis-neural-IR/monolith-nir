package memory

import (
	"monolith-nir/service/application/exception"
	"monolith-nir/service/application/repositories"
)

type IndexRepository struct {
	Index map[string][]string
}

func NewMemoryIndexRepository() repositories.IndexMemoryRepository {
	i := IndexRepository{
		Index: make(map[string][]string),
	}

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
	} else {
		return *exception.ThrowValidationError("Not have term indexed. Use Save function")
	}

	return nil
}

func (i *IndexRepository) Save(term string, documents []string) error {
	localDocuments := i.Index[term]
	if len(localDocuments) > 0 {
		return *exception.ThrowValidationError("Term already exists. Use Update function")
	}

	i.Index[term] = documents
	return nil
}

package repositories

import "monolith-nir/service/application/domain"

type IndexRepository interface {
	FindByTerm(term string) (*domain.Index, error)
	Update(index domain.Index) error
	Save(index domain.Index) error
}

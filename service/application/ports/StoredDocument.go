package ports

import "monolith-nir/service/application/domain"

type Store interface {
	StoreDocument(document domain.Document) error
}

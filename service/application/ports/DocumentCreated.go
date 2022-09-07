package ports

import "monolith-nir/service/application/domain"

type DocumentEvent interface {
	Created(document domain.Document) error
}

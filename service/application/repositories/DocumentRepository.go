package repositories

import (
	"monolith-nir/service/application/domain"
)

type DocumentRepository interface {
	Save(document domain.Document) error
}

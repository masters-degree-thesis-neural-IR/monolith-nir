package usecases

import "monolith-nir/service/application/domain"

type CreateIndexUc interface {
	CreateIndex(id string, title string, body string) error
	CreateNormalizedDocument(ch chan domain.NormalizedDocument)
}

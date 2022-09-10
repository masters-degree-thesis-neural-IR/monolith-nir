package mysq

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"monolith-nir/service/application/domain"
	"monolith-nir/service/application/exception"
	"monolith-nir/service/application/repositories"
)

type DocumentEmbeddingRepository struct {
	DB    *sql.DB
	Cache map[string]*domain.DocumentEmbedding
}

func NewDocumentEmbeddingRepository(db *sql.DB) repositories.DocumentEmbeddingRepository {
	return DocumentEmbeddingRepository{
		DB:    db,
		Cache: make(map[string]*domain.DocumentEmbedding),
	}
}

func (d DocumentEmbeddingRepository) FindByDocumentIDs(documentIDs map[string]int8) ([]domain.DocumentEmbedding, error) {
	var documentsEmbedding []domain.DocumentEmbedding

	nocache := make([]string, 100)

	//verifica quais documentos estão no cache
	for id, _ := range documentIDs {
		document := d.Cache[id]
		if document != nil {
			documentsEmbedding = append(documentsEmbedding, *document)
		} else {
			nocache = append(nocache, id)
		}
	}

	var in string
	for _, id := range nocache {
		in += fmt.Sprintf("'%s',", id)
	}
	in += "''"

	query := fmt.Sprintf("SELECT * FROM tb_document_embedding WHERE id IN (%s)", in)
	result, err := d.DB.Query(query)

	if err != nil {
		return nil, exception.ThrowUnexpectedError(err.Error())
	}

	for result.Next() {

		var id string
		var embedding string

		err = result.Scan(&id, &embedding)
		if err != nil {
			return nil, exception.ThrowUnexpectedError(err.Error())
		}

		var documentEmbedding domain.DocumentEmbedding
		json.Unmarshal([]byte(embedding), &documentEmbedding)

		d.Cache[id] = &documentEmbedding

		documentsEmbedding = append(documentsEmbedding, documentEmbedding)
	}

	defer result.Close()

	return documentsEmbedding, nil
}

func (d DocumentEmbeddingRepository) Save(documentEmbedding domain.DocumentEmbedding) error {

	jsonDocument, err := json.Marshal(documentEmbedding)

	if err != nil {
		return exception.ThrowUnexpectedError(err.Error())
	}

	insert, err := d.DB.Prepare("INSERT INTO tb_document_embedding VALUES(?,?)")

	if err != nil {
		return exception.ThrowUnexpectedError(err.Error())
	}

	insert.Exec(documentEmbedding.Id, string(jsonDocument))

	defer insert.Close()

	return nil
}

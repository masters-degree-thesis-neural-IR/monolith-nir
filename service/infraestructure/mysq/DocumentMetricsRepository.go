package mysq

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"monolith-nir/service/application/domain"
	"monolith-nir/service/application/exception"
	"monolith-nir/service/application/repositories"
)

type DocumentMetricsRepository struct {
	DB    *sql.DB
	Cache map[string]*domain.NormalizedDocument
}

func NewDocumentMetricsRepository(db *sql.DB) repositories.DocumentMetricsRepository {
	return DocumentMetricsRepository{
		DB:    db,
		Cache: make(map[string]*domain.NormalizedDocument),
	}
}

func (d DocumentMetricsRepository) FindByDocumentIDs(documentIDs map[string]int8) ([]domain.NormalizedDocument, error) {

	var normalizedDocuments []domain.NormalizedDocument

	nocache := make([]string, 0)

	//verifica quais documentos estão no cache
	for id, _ := range documentIDs {
		document := d.Cache[id]
		if document != nil {
			normalizedDocuments = append(normalizedDocuments, *document)
		} else {
			nocache = append(nocache, id)
		}
	}

	var in string
	for _, id := range nocache {
		in += fmt.Sprintf("'%s',", id)
	}
	in += "''"
	//println(in)

	query := fmt.Sprintf("SELECT * FROM tb_document_metrics WHERE id IN (%s)", in)
	result, err := d.DB.Query(query)

	if err != nil {
		return nil, exception.ThrowUnexpectedError(err.Error())
	}

	for result.Next() {

		var id string
		var metrics string

		err = result.Scan(&id, &metrics)
		if err != nil {
			return nil, exception.ThrowUnexpectedError(err.Error())
		}

		var normalizedDocument domain.NormalizedDocument
		json.Unmarshal([]byte(metrics), &normalizedDocument)

		d.Cache[id] = &normalizedDocument

		normalizedDocuments = append(normalizedDocuments, normalizedDocument)
	}

	defer result.Close()

	return normalizedDocuments, nil

}

func (d DocumentMetricsRepository) Save(document domain.NormalizedDocument) error {

	jsonDocument, err := json.Marshal(document)

	if err != nil {
		return exception.ThrowUnexpectedError(err.Error())
	}

	insert, err := d.DB.Prepare("INSERT INTO tb_document_metrics VALUES(?,?)")

	if err != nil {
		return exception.ThrowUnexpectedError(err.Error())
	}

	insert.Exec(document.Id, string(jsonDocument))

	defer insert.Close()

	return nil

}

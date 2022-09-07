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
	DB *sql.DB
}

func NewDocumentMetricsRepository(db *sql.DB) repositories.DocumentMetricsRepository {
	return DocumentMetricsRepository{
		DB: db,
	}
}

func (d DocumentMetricsRepository) FindDocuments(documentIDs map[string]int8) ([]domain.NormalizedDocument, error) {

	var in string
	for id, _ := range documentIDs {
		in += fmt.Sprintf("'%s',", id)
	}
	in += "''"

	query := fmt.Sprintf("SELECT * FROM tb_document_metrics WHERE id IN (%s)", in)
	result, err := d.DB.Query(query)

	if err != nil {
		return nil, *exception.ThrowUnexpectedError(err.Error())
	}

	var normalizedDocuments []domain.NormalizedDocument

	for result.Next() {

		var document domain.NormalizedDocument

		err = result.Scan(&document)
		if err != nil {
			return nil, *exception.ThrowUnexpectedError(err.Error())
		}

		normalizedDocuments = append(normalizedDocuments, document)
	}

	defer result.Close()

	return normalizedDocuments, nil

}

func (d DocumentMetricsRepository) Save(document domain.NormalizedDocument) error {

	jsonDocument, err := json.Marshal(document)

	if err != nil {
		return *exception.ThrowUnexpectedError(err.Error())
	}

	insert, err := d.DB.Prepare("INSERT INTO tb_document_metrics VALUES(?,?)")

	if err != nil {
		return *exception.ThrowUnexpectedError(err.Error())
	}

	insert.Exec(document.Id, string(jsonDocument))

	defer insert.Close()

	return nil

}

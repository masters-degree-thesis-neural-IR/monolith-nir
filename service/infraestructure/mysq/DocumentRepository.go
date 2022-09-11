package mysq

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"monolith-nir/service/application/domain"
	"monolith-nir/service/application/exception"
	"monolith-nir/service/application/repositories"
)

type DocumentRepository struct {
	DB    *sql.DB
	Cache map[string]*domain.Document
}

func NewDocumentRepository(db *sql.DB) repositories.DocumentRepository {
	return DocumentRepository{
		DB:    db,
		Cache: make(map[string]*domain.Document),
	}
}

func (d DocumentRepository) FindByDocumentIDs(documentIDs []string) (map[string]domain.Document, error) {

	documents := make(map[string]domain.Document)
	nocache := make([]string, 100)

	//verifica quais documentos estão no cache
	for _, id := range documentIDs {
		document := d.Cache[id]
		if document != nil {
			documents[id] = *document
		} else {
			nocache = append(nocache, id)
		}
	}

	var in string
	for _, id := range nocache {
		in += fmt.Sprintf("'%s',", id)
	}
	in += "''"

	query := fmt.Sprintf("SELECT * FROM tb_document WHERE id IN (%s)", in)
	result, err := d.DB.Query(query)

	if err != nil {
		return nil, exception.ThrowUnexpectedError(err.Error())
	}

	for result.Next() {

		document := domain.Document{}
		err = result.Scan(&document.Id, &document.Title, &document.Body)
		if err != nil {
			return nil, exception.ThrowUnexpectedError(err.Error())
		}

		d.Cache[document.Id] = &document
		documents[document.Id] = document
	}

	defer result.Close()

	return documents, nil
}

func (d DocumentRepository) Save(document domain.Document) error {

	//query := fmt.Sprintf("INSERT INTO tb_document VALUES('%s','%s','%s')",
	//	document.Id, document.Title, document.Body)

	//insert, err := d.DB.Query(query)

	insert, err := d.DB.Prepare("INSERT INTO tb_document VALUES(?,?,?)")

	if err != nil {
		return exception.ThrowUnexpectedError(err.Error())
	}

	insert.Exec(document.Id, document.Title, document.Body)

	defer insert.Close()

	return nil
}

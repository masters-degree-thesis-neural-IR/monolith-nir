package mysq

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"monolith-nir/service/application/domain"
	"monolith-nir/service/application/exception"
	"monolith-nir/service/application/repositories"
)

type DocumentRepository struct {
	DB *sql.DB
}

func NewDocumentRepository(db *sql.DB) repositories.DocumentRepository {
	return DocumentRepository{
		DB: db,
	}
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

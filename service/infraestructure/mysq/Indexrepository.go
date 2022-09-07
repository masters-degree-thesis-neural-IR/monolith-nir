package mysq

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"monolith-nir/service/application/domain"
	"monolith-nir/service/application/exception"
	"monolith-nir/service/application/repositories"
	"strings"
)

type IndexRepository struct {
	DB *sql.DB
}

func NewIndexRepository(db *sql.DB) repositories.IndexRepository {
	return IndexRepository{
		DB: db,
	}
}

type Index struct {
	Term      string
	Documents string
}

func (i IndexRepository) FindByTerm(term string) (*domain.Index, error) {

	query := fmt.Sprintf("SELECT * FROM tb_index WHERE term = '%s'", strings.ReplaceAll(term, "'", ""))
	result, err := i.DB.Query(query)

	if err != nil {
		return nil, *exception.ThrowUnexpectedError(err.Error())
	}

	var index Index

	for result.Next() {
		err = result.Scan(&index.Term, &index.Documents)
		if err != nil {
			return nil, *exception.ThrowUnexpectedError(err.Error())
		}
	}

	if index.Term == "" {
		return nil, nil
	}

	defer result.Close()

	var normalizedDocument []domain.NormalizedDocument
	json.Unmarshal([]byte(index.Documents), &normalizedDocument)

	return &domain.Index{
		Term:      index.Term,
		Documents: normalizedDocument,
	}, nil

}

func (i IndexRepository) Update(index domain.Index) error {

	documents, err := json.Marshal(index.Documents)

	if err != nil {
		return *exception.ThrowUnexpectedError(err.Error())
	}

	update, err := i.DB.Prepare("UPDATE tb_index SET documents =? WHERE term =?")

	if err != nil {
		panic(err.Error())
	}

	update.Exec(documents, index.Term)

	defer update.Close()

	return nil

}

func (i IndexRepository) Save(index domain.Index) error {

	documents, err := json.Marshal(index.Documents)

	if err != nil {
		return *exception.ThrowUnexpectedError(err.Error())
	}

	//query := fmt.Sprintf("INSERT INTO tb_index VALUES('%s','%s')",
	//	index.Term, string(documents))

	//insert, err := i.DB.Query(query)
	insert, err := i.DB.Prepare("INSERT INTO tb_index VALUES(?,?)")

	if err != nil {
		return *exception.ThrowUnexpectedError(err.Error())
	}

	insert.Exec(index.Term, string(documents))

	defer insert.Close()

	return nil
}

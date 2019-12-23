package mysql

import (
	"database/sql"
	"github.com/elitenomad/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error)  {
	/*
	'	Add a insert statement to insert the data passed onto the table.
	 */
	statement := `INSERT INTO SNIPPETS (title, content, created, expires) VALUES(?, ?, UTC_TIMESTAMP(),DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	/*
		Use the exec method to execute the statement
	 */
	result, err := m.DB.Exec(statement, title, content, expires)
	if err != nil {
		return 0, err
	}

	/*
		Collect the Last insert ID and return it to convery the user the record is
		successfully inserted.
	 */
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	/*
		If id collected is valid return converted value of the id.
	 */
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error)  {
	return nil, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error)  {
	return nil, nil
}



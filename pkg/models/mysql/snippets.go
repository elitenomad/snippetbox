package mysql

import (
	"database/sql"
	"errors"
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
	/*
		Write an sql statement which fetches the data based on arguments
	 */
	statement := `SELECT ID, TITLE, CONTENT, CREATED FROM SNIPPETS WHERE EXPIRES > UTC_TIMESTAMP() AND ID = ?`

	/*
		Use the QUERY ROW Method to execute the statemnt
	 */
	row := m.DB.QueryRow(statement, id)

	/*
		Initialize a pointer to a new Zeroed snippet struct
	 */
	s := &models.Snippet{}

	/*
	 Use row.Scan to copy the values of each field onto pointer
	 */
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}else {
			return nil, err
		}
	}

	// If everything went OK
	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error)  {
	/*
		Statement
	 */
	stmt := `SELECT ID, TITLE, CONTENT, CREATED FROM SNIPPETS WHERE EXPIRES > UTC_TIMESTAMP() ORDER BY CREATED DESC LIMIT 10`

	/*
		Execute the statment using DB.Query
	 */
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return  nil, err
	}

	/*
		Very important to close the result set
	 */
	defer rows.Close()

	/*
		Initialze array of pointer
	 */
	snippets := []*models.Snippet{}
	
	/*
		Loop through the rows and append snippet object to snippets
	 */
	for rows.Next() {
		s := &models.Snippet{}

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}


	return snippets, nil
}



package mysql

import (
	"database/sql"
	"errors"
	"github.com/elitenomad/snippetbox/pkg/models"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserModel struct {
	DB *sql.DB
}


func (m *UserModel) Insert(name, email, password string) error {
	/*
		Create a bcrypt hash for the password before pushing it into
		the DB
	 */
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	statement := `INSERT INTO USERS (name, email, hashed_password, created) VALUES(?, ?, ?, UTC_TIMESTAMP())`

	/*
		Use the Exec method to run the above statement
	 */
	_, err = m.DB.Exec(statement, name, email, string(hashedPassword))

	/*
		Look for duplicated email if Yes return models.ErrDuplicateEmail error
	 */
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			/*
				Looking at DB code and string to return an error sounds very bad. TODO: Look for improvements here
			 */
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	/*
		GET the user record by email
		Un-hash the password
	 */
	statement := `SELECT id, hashed_password from USERS where email = ?`

	/*
		Use the QUERY ROW Method to execute the statement
	*/
	row := m.DB.QueryRow(statement, email)

	/*
		Initialize a pointer to a Zeroed user struct
	*/
	var id int
	var hashedPassword []byte

	/*
	 Use row.Scan to copy the values of each field onto pointer
	*/
	err := row.Scan(&id, &hashedPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		}else {
			return 0, err
		}
	}

	/*
		Check for the password correctness
	 */
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	/*
		If everything went OK
	 */
	return id, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	/*
		Zeroed user struct from models
	 */
	u := &models.User{}

	/*
		statement to fetch the user info by ID
	 */
	stmt := `SELECT id, name, email, created, active FROM users WHERE id = ?`

	/*
		As there is one row per ID, use DB.QueryRow
	*/
	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Email, &u.Created, &u.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	/*
		Everything OK then return below
	 */
	return u, nil
}
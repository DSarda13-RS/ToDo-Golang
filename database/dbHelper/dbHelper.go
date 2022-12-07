package dbHelper

import (
	"ToDo/database"
	"ToDo/models"
	"ToDo/utils"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
)

func IsUserExist(email string) (bool, error) {
	SQL := `SELECT id FROM users where email = $1 AND archived_at IS NULL`
	var id string
	err := database.Todo.Get(&id, SQL, email)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return true, nil
}

func CreateUser(db *sqlx.Tx, name, email, password string) (string, error) {
	SQL := `INSERT INTO users (name, email, password) VALUES ($1,$2,$3) RETURNING id`
	var userID string
	err := db.QueryRowx(SQL, name, email, password).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func CreateUserSession(db sqlx.Ext, userID, sessionToken string) error {
	SQL := `INSERT INTO user_session (user_id, session_token) VALUES ($1,$2) RETURNING id`
	_, err := db.Exec(SQL, userID, sessionToken)
	return err
}

func GetUserIDByPassword(email, password string) (string, error) {
	SQL := `SELECT id, password FROM users WHERE archived_at IS NULL AND email = $1`
	var userID string
	var passwordHash string
	err := database.Todo.QueryRowx(SQL, email).Scan(&userID, &passwordHash)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	if err == sql.ErrNoRows {
		return "", nil
	}
	if passwordErr := utils.CheckPassword(password, passwordHash); passwordErr != nil {
		return "", passwordErr
	}
	return userID, nil
}

func GetUserBySession(sessionToken string) (*models.User, error) {
	// language=SQL
	SQL := `SELECT 
       			u.id, 
       			u.name, 
       			u.email, 
       			u.created_at 
			FROM users u
			JOIN user_session us on u.id = us.user_id
			WHERE u.archived_at IS NULL AND us.session_token = $1`
	var user models.User
	err := database.Todo.Get(&user, SQL, sessionToken)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, nil
}

func DeleteSessionToken(userID, token string) error {
	// language=SQL
	SQL := `DELETE FROM user_session WHERE user_id = $1 AND session_token = $2`
	_, err := database.Todo.Exec(SQL, userID, token)
	return err
}

func CreateTask(db sqlx.Ext, userID, name, description string, pendingAt time.Time) error {
	SQL := `INSERT INTO todo (user_id, name, description, pending_at) VALUES ($1,$2,$3,$4) RETURNING id`
	_, err := db.Exec(SQL, userID, name, description, pendingAt)
	return err
}

func UpdateTask(db sqlx.Ext, name, description, userID, taskID string, pendingAt time.Time, isCompleted bool) error {
	SQL := `UPDATE todo SET name = $1, description = $2, pending_at = $3, mark_completed = $4, updated_at = NOW() WHERE user_id = $5 AND id = $6 AND archived_at IS NULL`
	_, err := db.Exec(SQL, name, description, pendingAt, isCompleted, userID, taskID)
	return err
}

func GetTask(userID, taskID string) (*models.Task, error) {
	SQL := `SELECT name, description, pending_at as DueDate, mark_completed as IsCompleted FROM todo WHERE user_id = $1 AND id = $2 AND archived_at IS NULL`
	var task models.Task
	err := database.Todo.Get(&task, SQL, userID, taskID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &task, nil
}

func DeleteTask(db sqlx.Ext, userID, taskID string) error {
	SQL := `UPDATE todo SET archived_at = NOW() WHERE user_id = $1 AND id = $2 AND archived_at IS NULL`
	_, err := db.Exec(SQL, userID, taskID)
	return err
}

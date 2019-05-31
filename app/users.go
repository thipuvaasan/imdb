package main

import (
	"database/sql"

	"github.com/raazcrzy/imdb/models"
	"github.com/raazcrzy/imdb/utils"
)

// fetchEmailForUser function fetches the email for a particular username
func fetchEmailForUser(userID, password string) (string, error) {
	row := utils.PgDB.QueryRow(`SELECT email FROM imdb.users WHERE user_id=$1 AND user_password=$2`, userID, password)

	var email string
	err := row.Scan(&email)
	if err != nil {
		Log.Errorln(err)
		return "", err
	}

	return email, nil
}

// isSuperAdmin function checks if the user is super admin
func isSuperAdmin(email string) bool {
	for i := 0; i < len(utils.Admins); i++ {
		if email == utils.Admins[i] {
			return true
		}
	}
	return false
}

// isAdmin function checks if the user has admin role
func isAdmin(email string) bool {
	row := utils.PgDB.QueryRow(`SELECT role FROM imdb.users WHERE email=$1`, email)

	var role sql.NullString
	err := row.Scan(&role)
	if err != nil {
		Log.Errorln(err)
		return false
	}

	if role.String == "admin" {
		return true
	}
	return false
}

// isAuthorizedUser function checks if the email and username matches
func isAuthorizedUser(email, userID string) bool {
	row := utils.PgDB.QueryRow(`SELECT COUNT(*) FROM imdb.users WHERE email=$1 AND user_id=$2`, email, userID)

	var num int
	err := row.Scan(&num)
	if err != nil {
		Log.Errorln(err)
		return false
	}

	if num > 0 {
		return true
	}
	return false
}

// createUser function creates a new user in the postgres database
func createUser(user models.User) (map[string]interface{}, error) {
	_, err := utils.PgDB.Exec(`INSERT INTO imdb.users(email, name, created_at, user_id, user_password, role) VALUES($1, $2, $3, $4, $5, $6);`, user.Email, user.Name, user.CreatedAt, user.UserName, user.UserPassword, user.Role)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_pkey"` {
			return map[string]interface{}{
				"message": "User already exists",
				"status":  400,
			}, nil
		}
		if err.Error() == `pq: duplicate key value violates unique constraint "users_user_password_key"` {
			return map[string]interface{}{
				"message": "user_name not unique",
				"status":  400,
			}, nil
		}
		Log.Errorln(err)
		return nil, err
	}
	return map[string]interface{}{
		"message": "user created successfully",
		"status":  201,
	}, nil
}

// deleteUser function deletes a user from the postgres database
func deleteUser(user string) (map[string]interface{}, error) {
	_, err := utils.PgDB.Exec(`DELETE FROM imdb.users WHERE email=$1;`, user)
	if err != nil {
		Log.Errorln(err)
		return nil, err
	}
	return map[string]interface{}{
		"message": "user deleted successfully",
		"status":  200,
	}, nil
}

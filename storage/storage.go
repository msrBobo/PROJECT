package storage

import (
	"PROJECT/models"
	"database/sql"

	_ "github.com/lib/pq"
)

func CreateUser(user *models.User) (*models.User, error) {

	db, err := Connect()

	if err != nil {
		return nil, err
	}
	defer db.Close()
	respUser := models.User{}

	err = db.QueryRow(`INSERT INTO users (id, name, last_name) VALUES ($1, $2, $3) RETURNING id, name, last_name`, user.Id, user.FirstName, user.LastName).Scan(&respUser.Id, &respUser.FirstName, &respUser.LastName)

	if err != nil {

		return nil, err
	}

	return &respUser, nil
}

func GetUser(userId string) (*models.User, error) {

	db, err := Connect()

	if err != nil {
		return nil, err
	}

	respUser := models.User{}

	err = db.QueryRow(`SELECT id, name, last_name FROM users WHERE id = $1`, userId).Scan(&respUser.Id, &respUser.FirstName, &respUser.LastName)

	if err != nil {

		return nil, err
	}

	return &respUser, nil
}

func GetAll(page, limit int) (users []*models.User, err error) {

	db, err := Connect()

	if err != nil {
		return nil, err
	}

	offset := limit * (page - 1)

	rows, err := db.Query(`SELECT id, name, last_name FROM users LIMIT $1 OFFSET $2`, limit,
		offset)

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.FirstName, &user.LastName)

		if err != nil {

			return nil, err
		}
		users = append(users, &user)

	}

	return users, nil
}


func Connect() (*sql.DB, error) {
	dbInfo := "user=bobo password=1234 dbname=bobo sslmode=disable"
	
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func UpdateUser(userId string, updatedUser *models.User) (*models.User, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	respUser := models.User{}

	err = db.QueryRow(`UPDATE users SET name = $1, last_name = $2 WHERE id = $3 RETURNING id, name, last_name`,
		updatedUser.FirstName, updatedUser.LastName, userId).
		Scan(&respUser.Id, &respUser.FirstName, &respUser.LastName)

	if err != nil {
		return nil, err
	}

	return &respUser, nil
}

func DeleteUser(userId string) (*models.User, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	respUser := models.User{}

	err = db.QueryRow(`DELETE FROM users WHERE id = $1 RETURNING id, name, last_name`, userId).
		Scan(&respUser.Id, &respUser.FirstName, &respUser.LastName)

	if err != nil {
		return nil, err
	}

	return &respUser, nil
}

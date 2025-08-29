package user

import (
	"database/sql"
	"fmt"
)

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) GetAll() ([]User, error) {
	rows, err := s.db.Query("SELECT id, username, password_hash, role, created, deleted FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.Id, &user.Username, &user.PasswordHash, &user.Role, &user.Created, &user.Deleted)
		if err != nil {
			panic(err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *UserStore) GetOne(id int64) (*User, error) {
	user := User{}
	row := s.db.QueryRow("SELECT * FROM users WHERE id = $1", id)
	err := row.Scan(&user.Id, &user.Username, &user.PasswordHash, &user.Role, &user.Created, &user.Deleted)
	if err != nil {
		panic(err)
	}

	fmt.Println(user)

	return &user, err
}

func (s *UserStore) Add(data User) (*User, error) {
	return nil, nil
}

func Delete(id int64, data User) (*User, error) {
	return nil, nil
}

func (s *UserStore) Update(int int64, data User) (*User, error) {
	return nil, nil
}

package models

import (
	"database/sql"

	"github.com/makcik45/jwt-go/config"
	"github.com/makcik45/jwt-go/entities"
)

type Usermodel struct {
	db *sql.DB
}

func NuserModel() *Usermodel {
	conn, err := config.DBconn()

	if err != nil {
		panic(err)
	}
	return &Usermodel{
		db: conn,
	}
}

func (u Usermodel) Where(user *entities.User, fildName, fildValue string) error {

	rows, err := u.db.Query("select * from user2 where "+fildName+"=? limit 1", fildValue)

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&user.Id, &user.NamaLengkap, &user.Email, &user.Username, &user.Password)
	}
	return nil
}

func (u Usermodel) Create(user entities.User) (int64, error) {

	result, err := u.db.Exec("insert into user2 (nama_lengkap, email, username, password) values (?, ?, ?, ?)",
		user.NamaLengkap, user.Email, user.Username, user.Password)

	if err != nil {
		return 0, err
	}
	lastid, _ := result.LastInsertId()

	return lastid, nil
}

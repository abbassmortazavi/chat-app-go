package models

import (
	"backend/internal/db"
	"time"
)

type User struct {
	ID                   int64     `json:"id"`
	Email                string    `json:"email"`
	Name                 string    `json:"name"`
	Password             string    `json:"-"`
	RefreshTokenWeb      string    `json:"-"`
	RefreshTokenMobile   string    `json:"-"`
	RefreshTokenWebAt    time.Time `json:"-"`
	RefreshTokenMobileAt time.Time `json:"-"`
	CreatedAt            time.Time `json:"created_at"`
}

func GetUserByEmail(email string) (*User, error) {
	u := &User{}
	err := db.DB.QueryRow(`select
   		 id,
   		 email,
   		 name,
   		 password,
   		 refresh_token_web,
   		 refresh_token_mobile,
   		 refresh_token_web_at,
   		 refresh_token_mobile_at,
   		 created_at from users where email=?`, email).Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&u.Password,
		&u.RefreshTokenWeb,
		&u.RefreshTokenMobile,
		&u.RefreshTokenWebAt,
		&u.RefreshTokenMobileAt,
		&u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func CreateUser(email, name, hashedPassword string) (*User, error) {

	res, err := db.DB.Exec(`insert into users (name,email,password) values (?,?,?)`, name, email, hashedPassword)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()

	return &User{
		ID:        id,
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}, nil

}

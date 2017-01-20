package server

import (
	"database/sql"
	"errors"
	"time"

	context "golang.org/x/net/context"

	"github.com/arbarlow/account_service/account"
	"github.com/jinzhu/gorm"
)

type Account struct {
	Id        string `sql:"id;type:uuid;primary_key;default:gen_random_uuid()"`
	Name      sql.NullString
	Email     sql.NullString `gorm:"not null;unique" valid:"email"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AccountServer struct {
	account.AccountServiceServer
	db *gorm.DB
}

func (as *AccountServer) DBConnect(conn string) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	res := db.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto;")
	if res.Error != nil {
		return nil, res.Error
	}
	res = db.AutoMigrate(&Account{})
	if res.Error != nil {
		return nil, res.Error
	}

	as.db = db

	return db, nil
}

func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func (as AccountServer) Create(ctx context.Context, r *account.AccountCreateRequest) (*account.AccountDetails, error) {
	a := Account{
		Name:  nullString(r.Name),
		Email: nullString(r.Email),
	}

	err := as.db.Create(&a).Error

	return &account.AccountDetails{
		Id:    a.Id,
		Name:  a.Name.String,
		Email: a.Email.String,
	}, err
}

func (as AccountServer) Read(ctx context.Context, r *account.AccountRequest) (*account.AccountDetails, error) {
	var a Account
	as.db.Where("id = ?", r.Id).First(&a)

	if a.Id == "" {
		return nil, errors.New("No account found")
	}

	return &account.AccountDetails{
		Id:    a.Id,
		Name:  a.Name.String,
		Email: a.Email.String,
	}, nil
}

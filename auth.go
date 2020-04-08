package main

import (
	"github.com/VolantMQ/vlapi/vlauth"
	_ "github.com/lib/pq"
	"database/sql"
	"fmt"
)

type UserModel struct {
	Username           string   
	Password           string
	SubscriptionList   []string 
	PublishList        []string
}

type authProvider struct {
	cfg        config
	db_connection pq.Conn
}

func (p *authProvider) Connect() error {
	connStr := p.cfg.postgresUrl
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	p.db_connection = db
	ping_err = db.Ping()
	if ping_err != nil {
		return ping_err
	}
	return nil
}

func (p *authProvider) Init() error {
	err := p.Connect()
	return err
}

func (p *authProvider) Finduser(username string, password string) (user UserModel, error error) {
	db_query := fmt.Sprintf("SELECT * FROM %s WHERE username = $1 AND password = $2", p.cfg.postgresUserTable)
	rows, err := p.connection.Query(db_query, username, password)
	if err != nil {
		error = err
		return user, error
	}
	return user, nil
}
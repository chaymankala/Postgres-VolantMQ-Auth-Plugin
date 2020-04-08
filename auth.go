package main

import (
	"github.com/VolantMQ/vlapi/vlauth"
	"fmt"
	"github.com/jackc/pgx/v4"
	"context"
)

type UserModel struct {
	Username           string   
	Password           string
	SubscriptionList   []string 
	PublishList        []string
}

type authProvider struct {
	cfg        config
	db_connection *pgx.Conn
}

func (p *authProvider) Connect() error {
	connStr := p.cfg.postgresUrl
	db, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return err
	}
	// defer db.Close(context.Background())
	p.db_connection = db
	return nil
}

func (p *authProvider) Init() error {
	err := p.Connect()
	return err
}

func (p *authProvider) Finduser(username string, password string) (user UserModel, error error) {
	db_query := fmt.Sprintf("select Username,Password,SubscriptionList,PublishList FROM %s WHERE username = $1 AND password = $2", p.cfg.postgresUserTable)
	err := p.connection.Query(db_query, username, password).Scan(&user.Username, &user.Password, &user.SubscriptionList, &user.PublishList)
	if err != nil {
		error = err
		return user, err
	}
	return user, nil
}

func (p *authProvider) Password(clientID, username, password string) error {
	_, err := p.Finduser(username, password)
	if err != nil {
		return vlauth.StatusDeny
	}
	return vlauth.StatusAllow
}

func (p *authProvider) FindUserByUsername(username string) (user UserModel, err error) {
	db_query := fmt.Sprintf("select Username,Password,SubscriptionList,PublishList FROM %s WHERE username = $1", p.cfg.postgresUserTable)
	err := p.connection.Query(db_query, username).Scan(&user.Username, &user.Password, &user.SubscriptionList, &user.PublishList)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (p *authProvider) ACL(clientID, username, topic string, access vlauth.AccessType) error {
	user, err := p.FindUserByUsername(username)
	if err != nil {
		return vlauth.StatusDeny
	}
	permission := access.Type()
	allowed := false
	if permission == "write" {
		allowed = IsTopicAllowed(topic, user.PublishList)
	} else {
		allowed = IsTopicAllowed(topic, user.SubscriptionList)
	}
	if err == nil && allowed {
		return vlauth.StatusAllow
	}
	return vlauth.StatusDeny
}
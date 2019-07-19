package service

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
)

var g_service *Service

type Service struct {
	db             *sql.DB
	ProductService ProductService
}

func SetDB(db *sql.DB) {
	g_service = &Service{
		db:             db,
		ProductService: productService{db: db},
	}
}

func Share() *Service {
	if g_service.db == nil {
		panic(errors.New("service 没有初始化 db"))
	}
	return g_service
}

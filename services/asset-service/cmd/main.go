package main

import (
	_ "github.com/go-sql-driver/mysql"
)

const (
	mysqlDSN  = "user:password@tcp(localhost:3306)/dbname"
	redisAddr = "localhost:6379"
)

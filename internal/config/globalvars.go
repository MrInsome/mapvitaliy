package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

const (
	AccountID   = 1
	StatusConf  = http.StatusConflict
	StatusNotAl = http.StatusMethodNotAllowed
)

var CurrentAccount int
var Dsn string
var SecretKey string
var RedirectURL string
var Cmd string

// var UniToken = "6o1huty4e7se8tcmfcwr3fhcrx7cdh65aixjgsby"
func CallVarsFromENV() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("env файл не найден")
	}
	CurrentAccount = 1
	Dsn = os.Getenv("DSN")
	SecretKey = os.Getenv("SECRET_KEY")
	RedirectURL = os.Getenv("REDIRECT_URL")
	return nil
}

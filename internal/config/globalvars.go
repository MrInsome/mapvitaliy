package config

import (
	"net/http"
	"os"
)

const (
	AccountID   = 1
	StatusConf  = http.StatusConflict
	StatusNotAl = http.StatusMethodNotAllowed
)

var CurrentAccount = os.Getenv("CURRENT_ACCOUNT")
var CurrentIntegration = os.Getenv("CURRENT_INTEGRATION")
var Dsn = os.Getenv("DSN")
var SecretKey = os.Getenv("SECRET_KEY")
var RedirectURL = os.Getenv("REDIRECT_URL")

//var CurrentAccount = 1
//var CurrentIntegration = 0
//var Dsn = "steven:here@tcp(127.0.0.1:3306)/fullstack_api?charset=utf8mb4&parseTime=True&loc=Local"
//var SecretKey = "lQQtMWEtSWCvm0VFnhxkVv3PSk5D8sfFI54302pQTKN2k8GGmXIxCsB7tXyEfRou"
//var RedirectURL = "https://732b-212-46-197-210.eu.ngrok.io/vidget"

//var UniToken = "6o1huty4e7se8tcmfcwr3fhcrx7cdh65aixjgsby"

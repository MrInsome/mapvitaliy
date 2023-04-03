package config

import (
	"net/http"
)

const (
	AccountID   = 1
	StatusConf  = http.StatusConflict
	StatusNotAl = http.StatusMethodNotAllowed
)

var CurrentAccount = 1
var Dsn = "steven:here@tcp(127.0.0.1:3306)/fullstack_api?charset=utf8mb4&parseTime=True&loc=Local"
var SecretKey = "FNfIkCKHli5vA5mEKVlHrAjSH44Z1KfYE4G96MmqTW5JcxlAIxZwF4vPg2dDVYad"
var RedirectURL = "https://def0-173-233-147-68.eu.ngrok.io/vidget"

//var UniToken = "6o1huty4e7se8tcmfcwr3fhcrx7cdh65aixjgsby"

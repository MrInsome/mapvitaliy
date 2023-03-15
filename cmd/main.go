package main

import (
	. "apitraning/internal"
	. "apitraning/pkg"
	"log"
	"net/http"
)

func main() {
	repo := NewRepository()
	//Тестовый аккаунт
	var integration1 []Integration
	integration1 = append(integration1, Integration{
		SecretKey:          "SHYXnIWnNbViplLAS7obxKNZreqWkwemnQ1RFfCz5vvIxdiz9KKtvRinNj6GV9vD",
		ClientID:           "1c65fcfb-11d0-4bd6-931f-28b4d89892cf",
		RedirectURL:        "https://907a-212-46-197-210.eu.ngrok.io",
		AuthenticationCode: "def5020013652b86dbfaa8e7e4fae5628595a4a041f8d554e8e97bacbdf0088c8791d13b79fafcb4b6a88ba512679d261817fe075a29c55673ed67ca9687560ddc84a709d07e253b86d034ae15bda43c8d368e771be2e8cc456685107f4a5a4afda403093b1e21fa20395d804127f6c13aca780dcc1996265941c4d62ebaaecf1f1bbd56cdbdbe6b980382291beab3716d615eb7f065e428b30403ee04888077b3ddc2a1ec2935b327a6da6b7d1ab0c2e826c44b7de076ca7d18ed32e216cb1a9be5dc82c5607447b9c6241e170884d7ef7cc730e851d26487f4fff028536dcdefcfe12678656aa84105d1b7e6527932be2493ec4433f719c2760ee7d947b23c7d25a631b8100b219f922bccf647d8bef11f83a834e2bfca98551dadcc26e81186cca92f4682a988f6cf9f599cfdcf45ef526654d1455a302998b5600e7d96fd59de4b2c2d5531a056a29e558eea4e589f36cfd28990435902febcd4fb37d9c8edfa49e93d026dc4ba73a2ba4eb8380bd877581753a19a90bd625d42ca644120ec626397d4a0c3f3c65598e92770c6517d16f93d6b1426917be597bc57cbf54db0af76ee31ae9092a829b12db934088d6f218c8e045210591eea2d00df1d129258a7739b01b5974fd49c2984a71d24490b334847074efdd1f9fd77f3ebc8e8818d020918d61cb13c"})
	account1 := Account{
		AccessToken:  "",
		RefreshToken: "",
		Expires:      86400,
		AccountID:    1,
		Integration:  []Integration{},
	}
	repo.AddAccount(account1)
	repo.AddIntegration(account1.AccountID, integration1[0])

	handler := AccountsHandler(repo)
	integrationHandler := AccountIntegrationsHandler(repo)
	auth := AuthHandler(repo)
	requestHandler := AmoContat(repo)

	router := http.NewServeMux()
	router.Handle("/accounts", handler)
	router.Handle("/access_token", auth)
	router.Handle("/request", requestHandler)
	router.Handle("/accounts/integrations", integrationHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

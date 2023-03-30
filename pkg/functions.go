package pkg

import (
	"apitraning/internal"
	"apitraning/internal/config"
	"apitraning/internal/types"
	"bytes"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"net/mail"
	"net/url"
	"strconv"
	"strings"
)

type Repository struct {
	accounts     map[int]internal.Account
	integrations []internal.Integration
	contacts     []internal.Contacts
	data         map[int]types.DataToAccess
	referer      types.Referer
	db           *gorm.DB
}

func NewRepository() *Repository {
	return &Repository{
		accounts:     make(map[int]internal.Account),
		integrations: []internal.Integration{},
		contacts:     []internal.Contacts{},
		data:         make(map[int]types.DataToAccess),
		referer:      types.Referer{},
	}
}

func (r *Repository) RefererAdd(ref types.Referer) {
	r.referer = ref
}
func (r *Repository) RefererGet() types.Referer {
	return r.referer
}
func (r *Repository) AddAccount(account internal.Account) {
	r.accounts[account.AccountID] = account
}
func (r *Repository) DelAccount(account internal.Account) {
	delete(r.accounts, account.AccountID)
}
func (r *Repository) AddAuthData(accountID int) {
	data := types.DataToAccess{
		ClientID:     r.accounts[accountID].Integration[config.CurrentIntegration].ClientID,
		ClientSecret: r.accounts[accountID].Integration[config.CurrentIntegration].SecretKey,
		GrantType:    "authorization_code",
		Code:         r.accounts[accountID].Integration[config.CurrentIntegration].AuthenticationCode,
		RedirectUri:  r.accounts[accountID].Integration[config.CurrentIntegration].RedirectURL,
	}
	r.data[accountID] = data
}

func (r *Repository) AuthData(accountID int) types.DataToAccess {
	return r.data[accountID]
}

func (r *Repository) AddIntegration(accountID int, integration internal.Integration) {
	account, ok := r.accounts[accountID]
	if !ok {
		return
	}
	account.Integration = append(account.Integration, integration)
	r.accounts[accountID] = account
}
func (r *Repository) DelIntegration(accountID int, integration internal.Integration) {
	account, ok := r.accounts[accountID]
	if !ok {
		return
	}
	for i, el := range account.Integration {
		if el == integration {
			account.Integration[i] = account.Integration[len(account.Integration)-1]
			account.Integration[len(account.Integration)-1] = internal.Integration{}
			account.Integration = account.Integration[:len(account.Integration)-1]
		}
	}
	r.accounts[accountID] = account
}

func (r *Repository) UpdateIntegration(accountID int, integration internal.Integration, replaced internal.Integration) {
	account, ok := r.accounts[accountID]
	if !ok {
		return
	}
	for i, el := range account.Integration {
		if el == integration {
			account.Integration[i] = replaced
		}
	}
	r.accounts[accountID] = account

}

func (r *Repository) GetAccountIntegrations(accountID int) []internal.Integration {
	account, ok := r.accounts[accountID]
	if !ok {
		return nil
	}
	return account.Integration
}

func (r *Repository) GetAllAccounts() []internal.Account {
	accounts := make([]internal.Account, 0, len(r.accounts))
	for _, account := range r.accounts {
		accounts = append(accounts, account)
	}
	return accounts
}

func (r *Repository) GetAccount(accountID int) (internal.Account, error) {
	if accountID >= len(r.accounts) {
		return internal.Account{}, fmt.Errorf("Aккаунт %d не найден в нашей системе", accountID)
	}
	return r.accounts[accountID], nil
}

func (r *Repository) ContactsResp(n types.ContactResponce) []internal.Contacts {
	for _, v := range n.Response.Contacts {
		id := v.ID
		name := v.Name
		customFields := v.EmailValues
		for _, cf := range customFields {
			if cf.FieldCode == "EMAIL" {
				if cf.Values[0].Value != "" {
					_, err := mail.ParseAddress(cf.Values[0].Value)
					if err == nil {
						r.contacts = append(r.contacts, internal.Contacts{Name: name, ContactID: id, Email: cf.Values[0].Value})
					}
				}
			}
		}
	}
	return r.contacts
}
func (r *Repository) UnsubscribeAccount(db *gorm.DB, accountID int) error {
	account := r.accounts[accountID]
	if account.AccountID == 0 {
		return fmt.Errorf("")
	}
	r.DelAccount(account)
	db.Where("account_id = ?", account.AccountID).Delete(&internal.Contacts{})
	db.Where("account_id = ?", account.AccountID).Delete(&internal.Integration{})
	db.Delete(account)
	config.CurrentAccount = 1
	return nil
}
func (r *Repository) DBReturn() *gorm.DB {
	return r.db
}

func (r *Repository) GormOpen() {
	db, err := gorm.Open(mysql.Open(config.Dsn), &gorm.Config{})
	if err != nil {
		panic("Невозможно подключится к БД")
	}
	r.db = db
	err = r.db.AutoMigrate(&internal.Account{}, &internal.Integration{}, &internal.Contacts{})
	if err != nil {
		panic("Невозможно провести миграцию в БД")
	}
	r.SynchronizeDB(r.db)
}

func (r *Repository) SynchronizeDB(db *gorm.DB) {
	var account []internal.Account
	var contacts []internal.Contacts
	var integrations []internal.Integration
	db.Find(&account)
	for i, el := range account {
		db.Where("account_id = ?", el.AccountID).Find(&integrations)
		account[i].Integration = integrations
		db.Where("account_id = ?", el.AccountID).Find(&contacts)
		account[i].Contacts = contacts
		r.AddAccount(account[i])
	}
	db.Find(&contacts)
	db.Find(&integrations)
	r.contacts = contacts
	r.integrations = integrations
}

func GetARTokens(repo AccountAuth, db *gorm.DB, w http.ResponseWriter) error {
	account, err := repo.GetAccount(config.CurrentAccount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	ref := repo.RefererGet()
	var respToken types.TokenResponse
	repo.AddAuthData(config.CurrentAccount)
	a, err := json.Marshal(repo.AuthData(config.CurrentAccount))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	resp, err := http.Post("https://"+ref.Referer+"/oauth2/access_token",
		"application/json", bytes.NewBuffer(a))
	err = json.NewDecoder(resp.Body).Decode(&respToken)
	err = json.NewEncoder(w).Encode(respToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	account.AccessToken = respToken.AccessToken
	account.RefreshToken = respToken.RefreshToken
	account.Expires = respToken.ExpiresIn
	repo.AddAccount(account)
	db.Updates(account)
	return err
}
func importUni(apiKey string, repo AccountRepo) error {
	account, err := repo.GetAccount(config.CurrentAccount)
	if err != nil {
		return err
	}
	contacts := account.Contacts
	apiUrl := "https://api.unisender.com/ru/api/importContacts"
	data := url.Values{}
	data.Set("format", "json")
	data.Set("api_key", apiKey)
	data.Set("field_names[0]", "email")
	data.Set("field_names[1]", "Name")

	for i, el := range contacts {
		data.Set("data["+strconv.Itoa(i)+"][0]", el.Email)
		data.Set("data["+strconv.Itoa(i)+"][1]", el.Name)
	}

	resp, err := http.Post(apiUrl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func Router(repo *Repository) *http.ServeMux {
	Handler := AccountsHandler(repo)
	IntegrationHandler := AccountIntegrationsHandler(repo)
	Auth := AuthHandler(repo)
	RequestHandler := AmoContact(repo)
	GetFromAmoVidget := FromAMOVidget(repo)
	FromAmoUniKey := UnisenKey(repo)
	ImportUni := UnisenderImport(repo)

	router := http.NewServeMux()

	router.Handle("/vidget", GetFromAmoVidget)
	router.Handle("/accounts", Handler)
	router.Handle("/access_token", Auth)
	router.Handle("/request", RequestHandler)
	router.Handle("/accounts/integrations", IntegrationHandler)
	router.Handle("/vidget/unisender", FromAmoUniKey)
	router.Handle("/import", ImportUni)
	return router
}

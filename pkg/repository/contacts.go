package repository

import (
	"apitraning/internal/types"
	"fmt"
	"net/mail"
)

func (r *Repository) ParseContactsResponse(parseModel types.ContactResponce) []types.Contacts {
	for _, v := range parseModel.Response.Contacts {
		id := v.ID
		name := v.Name
		customFields := v.EmailValues
		for _, cf := range customFields {
			if cf.FieldCode == "EMAIL" {
				if cf.Values[0].Value != "" {
					_, err := mail.ParseAddress(cf.Values[0].Value)
					if err == nil {
						r.contacts = append(r.contacts, types.Contacts{Name: name, ContactID: id, Email: cf.Values[0].Value})
					}
				}
			}
		}
	}
	return r.contacts
}

func (r *Repository) AddUnsyncCon(id int, contact types.UnsyncContacts) {
	r.unSyncCon[id] = contact
}
func (r *Repository) GetUnsyncCon() ([]types.UnsyncContacts, error) {
	unsyncContacts := make([]types.UnsyncContacts, 0, len(r.accounts))
	for _, account := range r.unSyncCon {
		unsyncContacts = append(unsyncContacts, account)
	}
	if len(unsyncContacts) == 0 {
		return unsyncContacts, fmt.Errorf("В базе отсутствуют данные о несихронизированных контактах")
	}
	return unsyncContacts, nil
}

func (r *Repository) GetSyncCon() ([]types.Contacts, error) {
	Contacts := make([]types.Contacts, 0, len(r.accounts))
	for _, account := range r.contacts {
		Contacts = append(Contacts, account)
	}
	if len(Contacts) == 0 {
		return Contacts, fmt.Errorf("В базе отсутствуют данные о сихронизированных контактах")
	}
	return Contacts, nil
}

func (r *Repository) GetContact(conID int) (types.Contacts, error) {
	if r.contacts[conID].ContactID == 0 {
		return types.Contacts{}, fmt.Errorf("контакт %d не найден в нашей системе", conID)
	}
	return r.contacts[conID], nil
}
func (r *Repository) AddContact(contact types.Contacts) {
	r.accounts[contact.AccountID].Contacts[contact.ContactID] = contact
}

func (r *Repository) DelContact(account types.Account, contact types.Contacts) {
	account, ok := r.accounts[account.AccountID]
	if !ok {
		return
	}
	for i, el := range account.Contacts {
		if el == contact {
			account.Contacts[i] = account.Contacts[len(account.Contacts)-1]
			account.Contacts[len(account.Contacts)-1] = types.Contacts{}
			account.Contacts = account.Contacts[:len(account.Contacts)-1]
		}
	}
	r.accounts[account.AccountID] = account
}

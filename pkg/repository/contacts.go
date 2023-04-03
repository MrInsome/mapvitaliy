package repository

import (
	"apitraning/internal/types"
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

func (r *Repository) AddSyncCon(id int, contact types.Contacts) {
	r.contacts[id] = contact
}

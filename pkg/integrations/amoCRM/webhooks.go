package amoCRM

import (
	"apitraning/internal/types"
	"apitraning/pkg/repository"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func WebhookProducer(repo repository.BStalkWH) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Невозможно считать webhook", http.StatusBadRequest)
			return
		}
		accid := r.FormValue("account[id]")
		add := r.FormValue("contacts[add][0][id]")
		if add != "" {
			id, err := repo.Put([]byte(accid+" "+add+" "+
				r.FormValue("contacts[add][0][custom_fields][0][values][0][value]")+" "+
				r.FormValue("contacts[add][0][name]")), 1, 10, 100*time.Second)
			fmt.Fprintf(w, "Добавлено в очередь add %d", id)
			if err != nil {
				return
			}
		}

		update := r.FormValue("contacts[update][0][id]")
		if update != "" {
			_, err := repo.Put([]byte(accid+" "+update+" "+
				r.FormValue("contacts[update][0][custom_fields][0][values][0][value]")+" "+
				r.FormValue("contacts[update][0][name]")), 1, 10, 100*time.Second)
			fmt.Fprint(w, "Добавлено в очередь upd")
			if err != nil {
				return
			}
		}
		del := r.FormValue("contacts[delete][0][id]")
		if del != "" {
			_, err := repo.Put([]byte(accid+" "+del), 1, 10, 100*time.Second)
			fmt.Fprint(w, "Добавлено в очередь del")
			if err != nil {
				return
			}
		}
		return
	}
}

func WebhookWorker(repo repository.AccountContacts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for {
			id, body, err := repo.Reserve(5 * time.Second)
			if body == nil {
				fmt.Fprint(w, "Работа выполнена. В очереди пусто")
				break
			}
			splited := strings.Split(string(body), " ")
			accid, err := strconv.Atoi(splited[0])
			conid, err := strconv.Atoi(splited[1])
			contact := types.Contacts{
				AccountID: accid,
				ContactID: conid,
			}
			acc, err := repo.GetAccount(accid)
			if err != nil {
				fmt.Fprint(w, err)
				return
			}

			if strings.Contains(string(body), "add") {
				contact = types.Contacts{Name: splited[3], Email: splited[2]}
				acc.Contacts = append(acc.Contacts, contact)
				repo.DBReturn().Where("account_id = ?", accid).Updates(&types.Contacts{})
			}
			if strings.Contains(string(body), "update") {
				contact = types.Contacts{Name: splited[3], Email: splited[2]}
				repo.AddContact(contact)
				repo.DBReturn().Where("contactID = ?", conid).Updates(&types.Contacts{})
			}
			if strings.Contains(string(body), "delete") {
				repo.DelContact(acc, contact)
				repo.DBReturn().Where("contactID = ?", conid).Delete(&types.Contacts{})
			}
			err = repo.Delete(id)
			if err != nil {
				fmt.Errorf("ошибка воркера", id)
				return
			}
		}
	}
}

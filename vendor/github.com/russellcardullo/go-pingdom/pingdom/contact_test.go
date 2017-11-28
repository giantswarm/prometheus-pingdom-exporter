package pingdom

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestContactServiceList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/2.0/notification_contacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
      "contacts": [
        {
          "id": 11111111,
          "name": "John Doe",
          "type": "Notification contact"
        },
        {
          "id": 22222222,
          "name": "Jane Doe",
          "type": "Notification contact"
        }
      ]
    }`)
	})

	contacts, err := client.Contacts.List()
	if err != nil {
		t.Errorf("ListContacts returned error: %v", err)
	}

	want := []ContactResponse{
		ContactResponse{
			ID:   11111111,
			Name: "John Doe",
			Type: "Notification contact",
		},
		ContactResponse{
			ID:   22222222,
			Name: "Jane Doe",
			Type: "Notification contact",
		},
	}

	if !reflect.DeepEqual(contacts, want) {
		t.Errorf("ListContacts returned %+v, want %+v", contacts, want)
	}
}

func TestContactServiceCreate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/2.0/notification_contacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
      "contact":{
        "id":33333333,
        "name":"Michael Doe"
      }
    }`)
	})

	newContact := Contact{
		Name:      "Michael Doe",
		Email:     "michael.doe@site.com",
		Cellphone: "76543210",
	}
	contact, err := client.Contacts.Create(&newContact)
	if err != nil {
		t.Errorf("CreateContact returned error: %v", err)
	}

	want := &ContactResponse{ID: 33333333, Name: "Michael Doe"}
	if !reflect.DeepEqual(contact, want) {
		t.Errorf("CreateContact returned %+v, want %+v", newContact, want)
	}
}

func TestContactServiceUpdate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/2.0/notification_contacts/11111111", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{"message":"Modification of contact was successful!"}`)
	})

	updateContact := Contact{Name: "John Doe", Email: "jdoe@site.com"}
	msg, err := client.Contacts.Update(11111111, &updateContact)
	if err != nil {
		t.Errorf("UpdateContact returned error: %v", err)
	}

	want := &PingdomResponse{Message: "Modification of contact was successful!"}
	if !reflect.DeepEqual(msg, want) {
		t.Errorf("UpdateContact returned %+v, want %+v", msg, want)
	}
}

func TestContactServiceDelete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/2.0/notification_contacts/33333333", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{"message":"Deletion of contact was successful!"}`)
	})

	msg, err := client.Contacts.Delete(33333333)
	if err != nil {
		t.Errorf("DeleteContact returned error: %v", err)
	}

	want := &PingdomResponse{Message: "Deletion of contact was successful!"}
	if !reflect.DeepEqual(msg, want) {
		t.Errorf("DeleteContact returned %+v, want %+v", msg, want)
	}
}

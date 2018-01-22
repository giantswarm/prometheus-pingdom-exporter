package pingdom

import (
	"reflect"
	"testing"
)

func TestContactPutParams(t *testing.T) {
	contact := Contact{
		Name:               "John Doe",
		Cellphone:          "76543210",
		CountryCode:        "44",
		CountryISO:         "GB",
		DefaultSMSProvider: "nexmo",
		TwitterUser:        "real-john-doe",
		Email:              "john.doe@site.com",
	}
	params := contact.PutParams()

	want := map[string]string{
		"name":               contact.Name,
		"email":              contact.Email,
		"cellphone":          contact.Cellphone,
		"countryiso":         contact.CountryISO,
		"countrycode":        contact.CountryCode,
		"defaultsmsprovider": contact.DefaultSMSProvider,
		"directtwitter":      "false",
		"twitteruser":        contact.TwitterUser,
	}

	if !reflect.DeepEqual(params, want) {
		t.Errorf("Contact.PutParams() returned %+v, want %+v", params, want)
	}
}

func TestContactPostParams(t *testing.T) {
	contact := Contact{
		Name:               "John Doe",
		Cellphone:          "76543210",
		CountryCode:        "44",
		CountryISO:         "GB",
		DefaultSMSProvider: "nexmo",
		TwitterUser:        "real-john-doe",
		Email:              "john.doe@site.com",
	}
	params := contact.PostParams()

	want := map[string]string{
		"name":               contact.Name,
		"email":              contact.Email,
		"cellphone":          contact.Cellphone,
		"countryiso":         contact.CountryISO,
		"countrycode":        contact.CountryCode,
		"defaultsmsprovider": contact.DefaultSMSProvider,
		"directtwitter":      "false",
		"twitteruser":        contact.TwitterUser,
	}

	if !reflect.DeepEqual(params, want) {
		t.Errorf("Contact.PutParams() returned %+v, want %+v", params, want)
	}
}

func TestContactValid(t *testing.T) {
	contact := Contact{
		Name:      "John Doe",
		Cellphone: "76543210",
	}

	if err := contact.Valid(); err != nil {
		t.Errorf("Valid with valid contact returned error %+v", err)
	}

	contact = Contact{
		Name:      "",
		Cellphone: "76543210",
	}
	if err := contact.Valid(); err == nil {
		t.Errorf("Valid with invalid contact expected error, returned nil")
	}
}

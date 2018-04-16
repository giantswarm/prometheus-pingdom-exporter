package pingdom

import (
	"reflect"
	"testing"
)

func TestHttpCheckPutParams(t *testing.T) {
	check := HttpCheck{
		Name:     "fake check",
		Hostname: "example.com",
		Url:      "/foo",
		RequestHeaders: map[string]string{
			"User-Agent": "Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)",
			"Pragma":     "no-cache",
		},
		Username:       "user",
		Password:       "pass",
		ContactIds:     []int{11111111, 22222222},
		IntegrationIds: []int{33333333, 44444444},
	}
	params := check.PutParams()
	want := map[string]string{
		"name":                     "fake check",
		"host":                     "example.com",
		"paused":                   "false",
		"resolution":               "0",
		"sendtoemail":              "false",
		"sendtosms":                "false",
		"sendtotwitter":            "false",
		"sendtoiphone":             "false",
		"sendtoandroid":            "false",
		"sendnotificationwhendown": "0",
		"notifyagainevery":         "0",
		"notifywhenbackup":         "false",
		"use_legacy_notifications": "false",
		"url":              "/foo",
		"requestheader0":   "Pragma:no-cache",
		"requestheader1":   "User-Agent:Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)",
		"auth":             "user:pass",
		"encryption":       "false",
		"shouldnotcontain": "",
		"postdata":         "",
		"contactids":       "11111111,22222222",
		"integrationids":   "33333333,44444444",
		"tags":             "",
		"probe_filters":    "",
	}

	if !reflect.DeepEqual(params, want) {
		t.Errorf("Check.PutParams() returned %+v, want %+v", params, want)
	}
}

func TestHttpCheckPostParams(t *testing.T) {
	check := HttpCheck{
		Name:     "fake check",
		Hostname: "example.com",
		Url:      "/foo",
		RequestHeaders: map[string]string{
			"User-Agent": "Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)",
			"Pragma":     "no-cache",
		},
		Username:       "user",
		Password:       "pass",
		ContactIds:     []int{11111111, 22222222},
		IntegrationIds: []int{33333333, 44444444},
	}
	params := check.PostParams()
	want := map[string]string{
		"name":                     "fake check",
		"host":                     "example.com",
		"paused":                   "false",
		"resolution":               "0",
		"sendtoemail":              "false",
		"sendtosms":                "false",
		"sendtotwitter":            "false",
		"sendtoiphone":             "false",
		"sendtoandroid":            "false",
		"sendnotificationwhendown": "0",
		"notifyagainevery":         "0",
		"notifywhenbackup":         "false",
		"use_legacy_notifications": "false",
		"type":           "http",
		"url":            "/foo",
		"requestheader0": "Pragma:no-cache",
		"requestheader1": "User-Agent:Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)",
		"auth":           "user:pass",
		"encryption":     "false",
		"contactids":     "11111111,22222222",
		"integrationids": "33333333,44444444",
	}

	if !reflect.DeepEqual(params, want) {
		t.Errorf("Check.PostParams() returned %+v, want %+v", params, want)
	}
}

func TestHttpCheckValid(t *testing.T) {
	check := HttpCheck{Name: "fake check", Hostname: "example.com", Resolution: 15}
	if err := check.Valid(); err != nil {
		t.Errorf("Valid with valid check returned error %+v", err)
	}

	check = HttpCheck{Name: "fake check", Hostname: "example.com"}
	if err := check.Valid(); err == nil {
		t.Errorf("Valid with invalid check (`Resolution` == 0) expected error, returned nil")
	}

	check = HttpCheck{
		Name:             "fake check",
		Hostname:         "example.com",
		Resolution:       15,
		ShouldContain:    "foo",
		ShouldNotContain: "bar",
	}
	if err := check.Valid(); err == nil {
		t.Errorf("Valid with invalid check (`ShouldContain` and `ShouldNotContain` defined) expected error, returned nil")
	}

}

func TestPingCheckPostParams(t *testing.T) {
	check := PingCheck{Name: "fake check", Hostname: "example.com", ContactIds: []int{11111111, 22222222}, IntegrationIds: []int{33333333, 44444444}}
	params := check.PostParams()
	want := map[string]string{
		"name":                     "fake check",
		"host":                     "example.com",
		"paused":                   "false",
		"resolution":               "0",
		"sendtoemail":              "false",
		"sendtosms":                "false",
		"sendtotwitter":            "false",
		"sendtoiphone":             "false",
		"sendtoandroid":            "false",
		"sendnotificationwhendown": "0",
		"notifyagainevery":         "0",
		"notifywhenbackup":         "false",
		"use_legacy_notifications": "false",
		"type":           "ping",
		"contactids":     "11111111,22222222",
		"integrationids": "33333333,44444444",
		"probe_filters":  "",
	}

	if !reflect.DeepEqual(params, want) {
		t.Errorf("Check.PostParams() returned %+v, want %+v", params, want)
	}
}

func TestPingCheckValid(t *testing.T) {
	check := PingCheck{Name: "fake check", Hostname: "example.com", Resolution: 15}
	if err := check.Valid(); err != nil {
		t.Errorf("Valid with valid check returned error %+v", err)
	}

	check = PingCheck{Name: "fake check", Hostname: "example.com"}
	if err := check.Valid(); err == nil {
		t.Errorf("Valid with invalid check expected error, returned nil")
	}
}

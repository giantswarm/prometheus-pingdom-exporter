package pingdom

import (
	"encoding/json"
	"testing"
)

var detailedCheckJson = `
{
	"id" : 85975,
	"name" : "My check 7",
	"resolution" : 1,
	"sendtoemail" : false,
	"sendtosms" : false,
	"sendtotwitter" : false,
	"sendtoiphone" : false,
	"sendnotificationwhendown" : 0,
	"notifyagainevery" : 0,
	"notifywhenbackup" : false,
	"created" : 1240394682,
	"type" : {
		"http" : {
			"url" : "/",
			"port" : 80,
			"requestheaders" : {
				"User-Agent" : "Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)",
				"Prama" : "no-cache"
			}
		}
	},
	"hostname" : "s7.mydomain.com",
	"status" : "up",
	"lasterrortime" : 1293143467,
	"lasttesttime" : 1294064823
}
`

func TestPingdomError(t *testing.T) {
	pe := PingdomError{StatusCode: 400, StatusDesc: "Bad Request", Message: "Missing param foo"}
	want := "400 Bad Request: Missing param foo"
	if e := pe.Error(); e != want {
		t.Errorf("Error() returned '%+v', want '%+v'", e, want)
	}

}

func TestCheckResponseUnmarshal(t *testing.T) {
	var ck CheckResponse

	err := json.Unmarshal([]byte(detailedCheckJson), &ck)
	if err != nil {
		t.Errorf("Error running json.Unmarshal for CheckResponse: '%+v'", err)
	}
	if ck.Type.Name != "http" {
		t.Errorf("CheckResponse.Type.Name should be populated. returned '%+v', want '%+v'", ck.Type.Name, "http")
	}

	if ck.Type.HTTP == nil {
		t.Errorf("CheckResponse.Type.HTTP should be populated.")
		return
	}
	rhl := len(ck.Type.HTTP.RequestHeaders)
	if rhl != 2 {
		t.Errorf("CheckResponse.Type.HTTP.RequestHeaders should be populated. length returned '%+v', want '%+v'", rhl, 2)
	}
}

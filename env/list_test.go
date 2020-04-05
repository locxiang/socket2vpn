package env

import (
	"gopkg.in/ffmt.v1"
	"testing"
)

func TestSetSocket5Port(t *testing.T) {

	user := "hello"

	SavePPTP(user, &PPTP{
		PPTPName:    "asdf",
		ConnectName: "asdfsad",
	})

	SetSocket5Port(user, 315)

	pptp := GetPPTP(user)

	ffmt.Print(pptp)
}

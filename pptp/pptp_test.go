package pptp

import "testing"

func TestGetConnectName(t *testing.T) {

	name, err := GetConnectName("dsfasdfsdfppp3 slksjalfj")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(name)
	}
}

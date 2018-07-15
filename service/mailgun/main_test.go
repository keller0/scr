package mailgun

import (
	"testing"
)

func TestSendSimpleMessage(t *testing.T) {

	id, e := SimpleMessage("nice to have you", "test content", "test@test.com")
	if e != nil {
		t.Error(e)
	}
	println(id)

}

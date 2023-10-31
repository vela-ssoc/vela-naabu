package naabu

import (
	"encoding/json"
	"testing"
)

func TestParam(t *testing.T) {
	param := NewTaskParam()

	chunk, err := json.Marshal(param)
	if err != nil {
		t.Logf("%v", err)
		return
	}
	t.Logf("%s", string(chunk))
}

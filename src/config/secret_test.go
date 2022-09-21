package config

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSecretString(t *testing.T) {
	// secret in string
	ss1 := SecretString("qwerty")
	assert.Equal(t, "*****", ss1.String())
	assert.Equal(t, "qwerty", ss1.Get())
	assert.Equal(t, "*****", fmt.Sprint(ss1))

	// secret in json
	sj1 := struct {
		Secret SecretString `json:"secret"`
	}{Secret: SecretString("qwerty")}
	j, _ := json.Marshal(sj1)
	assert.Equal(t, "{\"secret\":\"*****\"}", string(j))
}

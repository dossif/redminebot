package redmine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatTask(t *testing.T) {
	s1, _ := formatTask("qwerty123", 20)
	s2, _ := formatTask("qwerty123", 5)
	s3, _ := formatTask("Заголовок123", 20)
	s4, _ := formatTask("Заголовок123", 5)
	assert.Equal(t, "qwerty123", s1)
	assert.Equal(t, "qwert...", s2)
	assert.Equal(t, "Заголовок123", s3)
	assert.Equal(t, "Загол...", s4)
}

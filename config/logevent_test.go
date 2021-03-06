package config

import (
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_Format(t *testing.T) {
	var (
		assert   = assert.New(t)
		logevent LogEvent
		out      string
	)

	log.SetLevel(log.DebugLevel)

	logevent = LogEvent{
		Timestamp: time.Now(),
		Message:   "Test Message",
		Extra: map[string]interface{}{
			"int":    123,
			"float":  1.23,
			"string": "Test String",
			"time":   time.Now(),
		},
	}

	out = logevent.Format("%{message}")
	assert.Equal("Test Message", out)

	out = logevent.Format("%{@timestamp}")
	assert.NotEmpty(out)
	assert.NotEqual("%{@timestamp}", out)

	out = logevent.Format("%{int}")
	assert.Equal("123", out)

	out = logevent.Format("%{float}")
	assert.Equal("1.23", out)

	out = logevent.Format("%{string}")
	assert.Equal("Test String", out)

	out = logevent.Format("%{null}")
	assert.Equal("%{null}", out)

	logevent.AddTag("tag1", "tag2", "tag3")
	assert.Len(logevent.Tags, 3)
	assert.Contains(logevent.Tags, "tag1")

	logevent.AddTag("tag1", "tag%{int}")
	assert.Len(logevent.Tags, 4)
	assert.Contains(logevent.Tags, "tag123")
}

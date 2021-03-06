package config

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

type LogEvent struct {
	Timestamp time.Time              `json:"timestamp"`
	Message   string                 `json:"message"`
	Tags      []string               `json:"tags,omitempty"`
	Extra     map[string]interface{} `json:"-"`
}

const timeFormat = `2006-01-02T15:04:05.999999999Z`

func appendIfMissing(slice []string, s string) []string {
	for _, ele := range slice {
		if ele == s {
			return slice
		}
	}
	return append(slice, s)
}

func (self *LogEvent) AddTag(tag ...string) {
	for _, t := range tag {
		t = self.Format(t)
		self.Tags = appendIfMissing(self.Tags, t)
	}
}

func (self *LogEvent) Marshal() (raw []byte, err error) {
	event := map[string]interface{}{
		"@timestamp": self.Timestamp.UTC().Format(timeFormat),
		"message":    self.Message,
	}
	for key, value := range self.Extra {
		event[key] = value
	}
	if raw, err = json.Marshal(event); err != nil {
		log.Errorf("Marshal failed: %v\n%v", self, err)
		return
	}
	return
}

func (self *LogEvent) MarshalIndent() (raw []byte, err error) {
	event := map[string]interface{}{
		"@timestamp": self.Timestamp.UTC().Format(timeFormat),
		"message":    self.Message,
	}
	for key, value := range self.Extra {
		event[key] = value
	}
	if raw, err = json.MarshalIndent(event, "", "\t"); err != nil {
		log.Errorf("MarshalIndent failed: %v\n%v", self, err)
		return
	}
	raw = append(raw, '\n')
	return
}

func (self *LogEvent) Get(field string) (v interface{}) {
	switch field {
	case "@timestamp":
		v = self.Timestamp
	case "message":
		v = self.Message
	default:
		v = self.Extra[field]
	}
	return
}

func (self *LogEvent) GetString(field string) (v string) {
	switch field {
	case "@timestamp":
		v = self.Timestamp.UTC().Format(timeFormat)
	case "message":
		v = self.Message
	default:
		if value, ok := self.Extra[field]; ok {
			v = fmt.Sprintf("%v", value)
		}
	}
	return
}

func (self *LogEvent) Format(format string) (out string) {
	revar := regexp.MustCompile(`%{([\w@]+)}`)
	out = format
	matches := revar.FindAllStringSubmatch(out, -1)
	for _, submatches := range matches {
		field := submatches[1]
		value := self.GetString(field)
		if value != "" {
			out = strings.Replace(out, submatches[0], value, -1)
		}
	}
	return
}

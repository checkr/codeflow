package transistor

import (
	"encoding/json"
	"errors"
	"math/rand"
	"reflect"
	"time"
)

func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func SliceContains(name string, list []string) bool {
	for _, b := range list {
		if b == name {
			return true
		}
	}
	return false
}

func MapPayload(name string, event *Event) error {
	if typ, ok := ApiRegistry[name]; ok {
		d, _ := json.Marshal(event.Payload)
		val := reflect.New(reflect.TypeOf(typ))
		json.Unmarshal(d, val.Interface())
		event.Payload = val.Elem().Interface()
		return nil
	}
	return errors.New("no match")
}

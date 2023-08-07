package gei

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type validate struct {
	tagName  string
	splitStr string
}

func validatestruct(c *Context) bool {
	value := reflect.ValueOf(c)
	tye := value.Type()
	for i := 0; i < value.NumField(); i++ {
		tag := tye.Field(i).Tag.Get("validate")
		var (
			tagK string
			tagV string
		)
		equalIndex := strings.Index(tag, "=")
		if equalIndex != -1 {
			tagK = tag[0:equalIndex]
			tagV = tag[equalIndex+1:]
		}
		field := value.Field(i)
		switch field.Kind() {
		case reflect.String:
			if tag == "requiered" {
				if len(field.Interface().(string)) < 1 {
					return false
				}
			}
			if tag == "email" {

			}
			if tag == "password" {
				r := regexp.MustCompile(`^1[3-9]\d{9}$`)
				if !r.MatchString(field.Interface().(string)) {
					return false
				}

			}
		case reflect.Uint8:
			if tagK == "gt" {
				target, _ := strconv.Atoi((tagV))
				if field.Uint() <= uint64(target) {
					return false
				}
			}
		}
	}
	return true
}

package gei

import (
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type validate struct {
	tagName  string
	splitStr string
}

func validatestruct(requestdata *RequestData) bool {
	value := reflect.ValueOf(requestdata)
	if value.Kind() == reflect.Ptr {
		value = value.Elem() // 获取指针指向的实际值
	}
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
		log.Printf("%s:%s", tag, requestdata.Email)
		log.Printf("%s:%s", tag, requestdata.Password)
		log.Printf("%s:%s", tag, requestdata.Phone)
		switch field.Kind() {
		case reflect.String:
			if tag == "requiered" {
				if len(field.Interface().(string)) < 1 {
					log.Printf("requiered error")
					return false
				}
			}
			if tag == "email" {
				log.Printf("email error")
			}
			if tag == "phone" {
				r := regexp.MustCompile(`^1[3-9]\d{9}$`)
				if !r.MatchString(field.Interface().(string)) {
					log.Printf("phone error")
					return false
				}

			}
		case reflect.Uint8:
			if tagK == "gt" {
				target, _ := strconv.Atoi((tagV))
				if field.Uint() <= uint64(target) {
					log.Printf("gt error")
					return false
				}
			}

		}
	}
	return true
}

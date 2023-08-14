package gei

import (
	"fmt"
	"log"
	"time"
)

func Quiksendemail() {
	fmt.Println(" email  success")
}

func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func Signup() HandlerFunc {
	return func(c *Context) {

		type RequestData struct {
			Email    string `validate:"email"`
			Password string `validate:"password"`
			Phone    string `validate:"phone"`
		}

		newdata := &RequestData{
			Email:    c.Req.FormValue("email"),
			Phone:    c.Req.FormValue("phone"),
			Password: c.Req.FormValue("password"),
		}

		if validatestruct(newdata) {
			Quiksendemail()
			log.Printf(" sign  up success")
		} else {
			log.Printf(" email false")
		}
		c.Next()
	}

}

func login() HandlerFunc {
	return func(c *Context) {

		type loginfo struct {
			Userid   string `validate:"userid"`
			Password string `validate:"password"`
		}
		newloginfo := &loginfo{
			Userid:   c.Req.FormValue("userid"),
			Password: c.Req.FormValue("password"),
		}

		if validatestruct(newloginfo) {
			log.Printf("log in success")
		}
		c.Next()
	}
}

func logout() HandlerFunc {
	return func(c *Context) {
		log.Printf("log out success")
		c.Next()
	}
}

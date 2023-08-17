package gei

import (
	"fmt"
	"log"
	"time"
	//"gopkg.in/gomail.v2"
)

func Quiksendemail() HandlerFunc {
	return func(c *Context) {
		/*
			m := gomail.NewMessage()
			m.SetHeader("From", "your_email@example.com")
			m.SetHeader("To", "recipient@example.com")
			m.SetHeader("Subject", "Hello")
			d := gomail.NewDialer("smtp.example.com", 587, "your_email@example.com", "your_email_password")
			if err := d.DialAndSend(m); err != nil {
				log.Printf(" email failed")
				return
			}
		*/
		fmt.Println(" email  success")

	}

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

func Login() HandlerFunc {
	return func(c *Context) {

		type loginfo struct {
			Username string `validate:"username"`
			Password string `validate:"password"`
		}
		newloginfo := &loginfo{
			Username: c.Req.FormValue("username"),
			Password: c.Req.FormValue("password"),
		}

		if validatestruct(newloginfo) {
			log.Printf("log in success")
		}
		c.Next()
	}
}

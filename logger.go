package gei

import (
	"fmt"
	"log"
	"net/http"
	"time"
	//"gopkg.in/gomail.v2"
)

type L map[string]interface{}

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
		cookie := http.Cookie{
			Name:  cookiename,
			Value: "user-authenticated",
		}
		log.Printf("login  success")

		http.SetCookie(c.Writer, &cookie)
		c.JSON(http.StatusOK, L{
			"message": " loggin in and cookie successfully",
		})
		c.Next()
	}
}

func AuthenticateMiddleware() HandlerFunc {
	return func(c *Context) {
		cookie, err := c.Req.Cookie(cookiename)
		if err != nil || cookie.Value != "user-authenticated" {
			c.JSON(http.StatusUnauthorized, L{
				"error": "Unauthorized",
			})
			return
		}
		c.JSON(http.StatusOK, L{
			"message": " info successfully",
		})
		c.Next()
	}
}

package gei

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	//"gopkg.in/gomail.v2"
	_ "github.com/go-sql-driver/mysql"
)

type L map[string]interface{}

var db *sql.DB

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

		username := c.Req.FormValue("username")
		password := c.Req.FormValue("password")

		db, err := sql.Open("mysql", "root:990726@tcp(localhost:3306)/geidb")
		if err != nil {
			fmt.Println("Failed to connect to the database:", err)
			return
		}
		defer db.Close()
		var dbPassword string
		query := "SELECT password FROM users WHERE username = ?"
		err = db.QueryRow(query, username).Scan(&dbPassword)
		if err != nil {
			log.Printf("error: db query")
			return
		}

		if password != dbPassword {
			c.JSON(http.StatusOK, L{
				"message": " password is wrong",
			})
		}

		cookie := http.Cookie{
			Name:  cookiename,
			Value: username + "-authenticated",
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
		if err != nil || !strings.HasSuffix(cookie.Value, "-authenticated") {
			c.JSON(http.StatusUnauthorized, L{
				"error": "Unauthorized",
			})
			return
		}
		db, err := sql.Open("mysql", "root:990726@tcp(localhost:3306)/geidb")
		if err != nil {
			fmt.Println("Failed to connect to the database:", err)
			return
		}
		username := strings.TrimSuffix(cookie.Value, "-authenticated")
		var dbUsername, dbPassword string
		err = db.QueryRow("SELECT username, password FROM users WHERE username = ?", username).Scan(&dbUsername, &dbPassword)
		if err != nil {
			c.JSON(http.StatusOK, L{
				"error": " is not authenticated",
			})
			return
		}

		c.JSON(http.StatusOK, L{
			"username": dbUsername,
			"password": dbPassword,
		})
		c.Next()
	}
}

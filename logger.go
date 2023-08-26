package gei

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gomail.v2"
)

type L map[string]interface{}

var db *sql.DB

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
			//send email
			m := gomail.NewMessage()
			m.SetHeader("From", "1097598746@qq.com")
			m.SetHeader("To", "wang_yuxuan2007@163.com")
			m.SetHeader("Subject", "Hello")
			d := gomail.NewDialer("smtp.qq.com", 587, "1097598746@qq.com", "tvqgixzzwsrhihdi")
			if err := d.DialAndSend(m); err != nil {
				log.Printf(" email failed")
				return
			}

			fmt.Println(" email  success")
			//写入数据库
			db, err := sql.Open("mysql", "root:990726@tcp(localhost:3306)/geidb")
			if err != nil {
				fmt.Println("Failed to connect to the database:", err)
				return
			}
			defer db.Close()
			var count int
			err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", newdata.Email).Scan(&count)
			if err != nil {
				c.JSON(http.StatusInternalServerError, H{
					"error": "Internal Server Error",
				})
				return
			}

			if count > 0 {
				c.JSON(http.StatusBadRequest, H{
					"error": "Username already exists",
				})
				return
			}

			_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", newdata.Email, newdata.Password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, L{
					"error": "Internal Server Error",
				})
				return
			}

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

func dbstart() HandlerFunc {
	return func(c *Context) {
		db, err := sql.Open("mysql", "root:990726@tcp(localhost:3306)/geidb")
		if err != nil {
			fmt.Println("Failed to connect to the database:", err)
			return
		}
		defer db.Close()
		createTableSQL := `
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL
		);
	`
		_, err = db.Exec(createTableSQL)
		if err != nil {
			log.Fatal(err)
		}
	}
}

package gei

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func Login() HandlerFunc {
	return func(c *Context) {
		newdata := c.Newrequestdata()
		if validatestruct(newdata) {
			newdata.Quiksendemail()
			log.Printf("success")
		} else {
			log.Printf(" email false")
		}
	}

}

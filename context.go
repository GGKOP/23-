package gei

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
)

type H map[string]interface{}

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	Params     map[string]string
	StatusCode int
	Email      string `validate:"email"`
	Password   string `validate:"password"`
	Phone      string `validate:"phone"`
}

type RequestData struct {
	Email string `json:"email"`
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	var requestData RequestData

	err := json.NewDecoder(req.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	email := requestData.Email
	phone := req.FormValue("phone")
	password := req.FormValue("password")

	return &Context{
		Writer:   w,
		Req:      req,
		Path:     req.URL.Path,
		Method:   req.Method,
		Email:    email,
		Phone:    phone,
		Password: password,
	}
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
	contentType := c.Req.Header.Get("Content-Type")
	fmt.Println(" 1Content-Type is:", contentType)
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
	contentType := c.Req.Header.Get("Content-Type")
	fmt.Println("2 Content-Type is:", contentType)
}

func (c *Context) ServeHTMLFile(filename string) {
	fmt.Println("Trying to serve:", filename) // 添加日志
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)  // 添加日志
		c.Status(http.StatusInternalServerError) // 或其他适当的错误代码
		c.Writer.Write([]byte("Internal Server Error"))
		return
	}
	fmt.Println("Sending HTML content") // 添加日志
	c.HTML(http.StatusOK, string(data))
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) Validate() bool {
	return validatestruct(c)
}

func (c *Context) Quiksendemail() bool {
	host := "smtp.example.com"
	port := "587"
	address := host + ":" + port

	// 发件人信息
	sender := "sender@example.com"
	password := "password" // 发件人的密码或者SMTP授权码

	// 收件人
	recipient := c.Email

	// 邮件内容
	subject := "Subject: Hello!\n"
	body := "Hello from Go!"

	// 邮件包括头部（Subject），空行，和实际内容
	msg := []byte(subject + "\n" + body)

	// 设置认证信息
	auth := smtp.PlainAuth("", sender, password, host)

	// 发送邮件
	err := smtp.SendMail(address, auth, sender, []string{recipient}, msg)
	if err != nil {
		return false
	}

	return true
}

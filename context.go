package gei

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type H map[string]interface{}
type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	Params     map[string]string
	StatusCode int
	handlers   []HandlerFunc
	index      int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
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

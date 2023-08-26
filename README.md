# 须知
本项目是模仿 golang的web开源框架Gin，主要复现了 前缀路由 分组控制 中间件 等逻辑代码。
其中对/login，/logout，/signin，/info 等接口进行了定义和实现

## 如何安装

安装 Go：从 官方网站 下载并安装 Go，

创建项目文件夹结构：在工作目录下，创建一个项目文件夹，比如 myproject，然后将你的代码文件夹 gei 移动到 myproject 下。

安装依赖项：在项目的根目录下，打开终端或命令提示符，运行以下命令来安装项目的依赖项

```bash
Copy code
go get -u github.com/go-sql-driver/mysql
go get -u gopkg.in/gomail.v2
```
这将会自动下载并安装这些依赖项到你的 Go 工作区。

编译项目：在终端中切换到 myproject 目录下，运行以下命令编译你的项目：

```bash
Copy code
go build gei
```
这将会生成一个可执行文件（如果你的代码包含 main 函数）。
运行项目：运行生成的可执行文件来启动你的项目：

```bash
Copy code
./gei
```
或者，如果你使用 Windows：

```bash
Copy code
gei.exe
```

## 快速开始
```Go
package main

import (
	"gei"
	"net/http"
)

const (
	cookiename = "user-cookie"
	secretKey  = "your-secret-key"
)

func main() {
	r := gei.New()
	r.RunMiddware(gei.Logger())
	r.GET("/", func(c *gei.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	r.Run(":9999")
}
```
能快速验证/login，/logout，/signin，/info  几个功能、


## 配置说明

### 邮箱配置
```GO
m.SetHeader("From", "邮件地址1")
m.SetHeader("To", "邮件地址2")
d := gomail.NewDialer("smtp.邮件地址1的邮箱类型（qq/163）.com", 587, "邮件地址1", "邮件地址1密钥)
 ```

### 数据库配置
安装好数据库驱动后 ，就进行对数据库的设置
```Go
			db, err := sql.Open("mysql", "用户名:数据库密码@tcp(主机名:端口号)/数据库名")
			if err != nil {
				fmt.Println("Failed to connect to the database:", err)
				return
			}
```


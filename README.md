# myGin
基于net/http编写的路由框架

## 介绍

模仿gin框架编写的一个极简路由框架,实现了路由管理、路由分组、中间件管理等基础功能。


## example
```
package main


import (
	"fmt"
	"myGin"
)


func main() {
	//初始化操作：路由注册
	r:= myGin.NewEngine()
	r.Use(func(c * myGin.Context){
		fmt.Println("begin middle1")
		c.Next()
		fmt.Println("end middle1")
	})

	defaultRouters := r.Group("/")
	{
		defaultRouters.GET("/", func(c *myGin.Context) {
			c.Writer.Write([]byte("首页"))
		})
		defaultRouters.GET("/news", func(c *myGin.Context) {
			c.JSON(200, myGin.H{
				"message": "这是新闻首页",
			})
		})
	}

	apiRouters := r.Group("/api")
	{

		apiRouters.Use(func(c * myGin.Context){
			fmt.Println("begin middle2")
			c.Next()
			fmt.Println("end middle2")
		})
		apiRouters.GET("/", func(c *myGin.Context) {
			c.Writer.Write([]byte("我是一个api接口"))
		})
		apiRouters.GET("/userlist", func(c *myGin.Context) {
			c.Writer.Write([]byte("我是一个api接口-userlist"))
		})
		play := apiRouters.Group("xx")
		play.Use(func(c * myGin.Context){
			fmt.Println("begin middle3")
			c.Next()
			fmt.Println("end middle3")
		})
		play.GET("/plist", func(c *myGin.Context) {
			c.Writer.Write([]byte("我是一个api接口-xx/plist"))
		})
	}

	//启动服务
	r.Run("127.0.0.1:8000",r)

}

```  

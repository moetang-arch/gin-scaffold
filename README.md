# gin-startup

### example:

```
	g := gin_startup.NewGinStartup()
	g.EnableFastCgi("tcp://127.0.0.1:15001")
	g.EnableHttp("tcp://127.0.0.1:15002")
	g.Custom(func(r *gin.Engine) {
		r.Use(gin.Recovery(), gin.Logger())
		...
	})
	g.Start()
```

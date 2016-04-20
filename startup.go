package gin_startup

import (
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"net/url"

	"github.com/gin-gonic/gin"
)

type ginStartup struct {
	engine *gin.Engine

	enableFastCgi bool
	fastCgiBind   string

	enableHttp bool
	httpBind   string
}

func NewGinStartup() *ginStartup {
	g := new(ginStartup)
	g.engine = gin.New()
	g.enableFastCgi = false
	g.enableHttp = false
}

func (this *ginStartup) EnableFastCgi(bind string) {
	this.enableFastCgi = true
	this.fastCgiBind = bind
}

func (this *ginStartup) EnableHttp(bind string) {
	this.enableHttp = true
	this.httpBind = bind
}

func (this *ginStartup) Custom(f func(*gin.Engine)) {
	f(this.engine)
}

func (this *ginStartup) Start() error {
	if this.enableFastCgi {
		go func() {
			u, err := url.Parse(this.fastCgiBind)
			if err != nil {
				return
			}
			addr, err := net.ResolveTCPAddr("tcp", u.Host)
			if err != nil {
				panic(err)
			}
			listener, err := net.ListenTCP("tcp", addr)
			if err != nil {
				panic(err)
			}
			log.Println("[Fastcgi] starting fastcgi on ", this.fastCgiBind)
			err = fcgi.Serve(listener, this.engine)
			if err != nil {
				log.Println("[Fastcgi] serve fastcgi error. " + err)
			}
		}()
	}
	if this.enableHttp {
		go func() {
			var err error = nil
			if gin.IsDebugging() {
				log.Printf("[GIN-debug] Listening and serving HTTP on %s\n", this.httpBind)
			}
			defer func() {
				if err != nil && gin.IsDebugging() {
					log.Printf("[GIN-debug] [ERROR] %v\n", err)
				}
			}()

			server := &http.Server{
				Addr:    this.httpBind,
				Handler: this.engine,
			}

			err = server.ListenAndServe()
		}()
	}

	return nil
}

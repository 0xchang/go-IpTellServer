package main

import (
	"IpTellServer/controller"
	"IpTellServer/logmiddleware"
	"flag"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(logmiddleware.MyLog())
	r.GET("/getaddr", controller.GetAddr)
	//r.GET("/test", controller.Test)
	r.GET("/ip", controller.GetIp)
	r.GET("/get", controller.MyGet)
	var (
		port int
		host string
	)

	flag.IntVar(&port, "port", 7788, "Specify port")
	flag.StringVar(&host, "host", "0.0.0.0", "Specify host")
	flag.Parse()

	if port > 65535 || port < 0 {
		panic("Port error, only allowed between 0 and 65535.")
	}

	r.Run(host + ":" + strconv.Itoa(port))
}

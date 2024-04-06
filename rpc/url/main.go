package main

import (
	"log"
	"net"
	"shorturl/kitex_gen/short/url/urlservice"
	"shorturl/model"

	"github.com/cloudwego/kitex/server"
)

func main() {

	model.ConnDB()
	defer model.DB.Close()

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8889")

	svr := urlservice.NewServer(new(UrlServiceImpl),server.WithServiceAddr(addr))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}

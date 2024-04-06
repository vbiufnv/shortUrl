package main

import (
	"log"
	user "shorturl/kitex_gen/short/user/userservice"
	"shorturl/model"
)


func main() {
	model.ConnDB()
	defer model.DB.Close()
	

	svr := user.NewServer(new(UserServiceImpl))

	err := svr.Run()
	

	if err != nil {
		log.Println(err.Error())
	}
}

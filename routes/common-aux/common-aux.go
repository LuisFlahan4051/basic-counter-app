package routes

import (
	"log"
	"net/http"
)

func Logcatch(writer http.ResponseWriter, status int, err error) {
	if err != nil {
		log.Println(err.Error())
		writer.WriteHeader(status)
		writer.Write([]byte(err.Error()))
	}
}

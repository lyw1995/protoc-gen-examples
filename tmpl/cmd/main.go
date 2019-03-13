
package main

import (
	"github/track/protoc-gen-examples/tmpl"
	"log"
	"net/http"
)
	

func Login(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	if method != "POST"  {
		writer.WriteHeader(http.StatusInternalServerError)
		_,_ = writer.Write([]byte("request method un support"))
	}
	var _ user.LoginReq
	var _ user.LoginResp

	writer.WriteHeader(http.StatusOK)
	_,_=writer.Write([]byte("did "+method+" Login successful\n"))
}

func UserInfo(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	if method != "GET"  {
		writer.WriteHeader(http.StatusInternalServerError)
		_,_ = writer.Write([]byte("request method un support"))
	}
	var _ user.IdReq
	var _ user.LoginResp

	writer.WriteHeader(http.StatusOK)
	_,_=writer.Write([]byte("did "+method+" UserInfo successful\n"))
}


	
func main() {
	
	http.HandleFunc("/login", Login)
	
	http.HandleFunc("/userinfo", UserInfo)
	
	log.Fatal(http.ListenAndServe(":8000", nil))
}
			
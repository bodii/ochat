package comm

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResType struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

func ResFailure(w http.ResponseWriter, code int, msg string) {
	Res(w, code, msg, nil)
}

func ResSuccess(w http.ResponseWriter, data any) {
	Res(w, 200, "success!", data)
}

func Res(w http.ResponseWriter, code int, msg string, data any) {
	result := ResType{
		Code: code,
		Msg:  msg,
		Data: data,
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resultJson)
}

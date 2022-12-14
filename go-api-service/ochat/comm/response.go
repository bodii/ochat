package comm

import (
	"encoding/json"
	"log"
	"net/http"
)

type R struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

func Res(w http.ResponseWriter, code int, msg string, data any) {
	result := R{
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

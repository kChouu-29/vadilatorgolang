package user

import (
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	Ctrl *UserController
}

func NewUserHandler(u *UserController) *UserHandler {
	return &UserHandler{Ctrl: u}
}

func (u *UserHandler) writeJson(w http.ResponseWriter, data any){
	w.Header().Set("content-type","application/json")
	json.NewEncoder(w).Encode(data)
}

func (u *UserHandler) errorJson(w http.ResponseWriter,  status int,message string){
	u.writeJson(w,map[string]string{"error":message})
}
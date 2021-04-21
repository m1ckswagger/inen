package handlers

import (
  "github.com/m1ckswagger/inenp/server/internal/models"
  "log"
)



type Requests struct {
	l *log.Logger
	db models.Repository
}

func (rh *Requests) ListRequests(w http.ResponseWriter, r *http.Request) {

}

func (rh *Requests) Create(vmr models.VMRequest) VMRequest {
  panic("implement me")
}




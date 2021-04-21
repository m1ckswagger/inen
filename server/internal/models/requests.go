package models

import (
  "encoding/json"
  "io"
)

type RequestStatus int

const (
  StatusRequested RequestStatus = iota
  StatusRejected
  StatusAccepted
  StatusRetired
  StatusInProgress
)

func (r RequestStatus) String() string {
  return [...]string{"Requested", "Rejected", "Accepted", "Retired", "In progress"}[r]
}

// VMRequest defines the structure of an API VM Request
type VMRequest struct {
  ID       int           `json:"id,omitempty"`
  Hostname string        `json:"hostname"`
  Email    string        `json:"requester"`
  OS       string        `json:"os"`
  HDD      int           `json:"hdd"`
  CPU      int           `json:"cpu"`
  RAM      int           `json:"ram"`
  Flavor   string        `json:"flavor"`
  Networks []string      `json:"networks"`
  Status   RequestStatus `json:"status"`
}

// FromJSON unmarshals JSON read from r into VMRequest
func (v *VMRequest) FromJSON(r io.Reader) error {
  d := json.NewDecoder(r)
  return d.Decode(v)
}

type VMRequests []VMRequest

// ToJSON writes serialized JSON to w
func (v *VMRequests) ToJSON(w io.Writer) error {
  e := json.NewEncoder(w)
  return e.Encode(v)
}

type Repository interface {
  List() (VMRequests, error)
  Create(request VMRequest) (VMRequest, error)
}

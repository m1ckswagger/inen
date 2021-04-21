package server

import "fmt"

type UserMemoryStore struct {
  DB map[string]string
}

func (u *UserMemoryStore) Find(username string) (*User, error) {
  pass, ok := u.DB[username]
  if !ok {
    return nil, fmt.Errorf("could not find user %q", username)
  }
  return &User{Login: username, Password: pass}, nil
}

func (u *UserMemoryStore) Auth(username, password string) bool {
  pass, ok := u.DB[username]
  if !ok {
    return false
  }
  return pass == password
}

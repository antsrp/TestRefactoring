package db

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"refactoring/internal/user"

	"github.com/pkg/errors"
)

const (
	UserNotFound = "User not found"
)

type UserStorage struct {
	path string
	user.UserStore
}

var (
	ErrUserNotFound = errors.New(UserNotFound)

	_ user.Storage = &UserStorage{}
)

func CreateUserStorage(name string) (*UserStorage, error) {
	path := fmt.Sprintf("%s//%s", getPathToStorageFolder(), name)
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "can't read file of store")
	}
	s := user.UserStore{}
	if err = json.Unmarshal(f, &s); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal data of store")
	}
	return &UserStorage{
		path:      path,
		UserStore: s,
	}, nil
}

func (s *UserStorage) Add(u *user.User) int {
	s.Increment++
	s.List[s.Increment] = *u

	return s.Increment
}

func (s *UserStorage) GetUser(id int) (*user.User, error) {
	if val, found := s.List[id]; !found { // user not found
		return nil, ErrUserNotFound
	} else {
		return &val, nil
	}
}

func (s *UserStorage) GetUsers() *user.UserStore {
	return &s.UserStore
}

func (s *UserStorage) Delete(id int) error {
	if _, found := s.List[id]; !found { // user not found
		return ErrUserNotFound
	}
	delete(s.List, id)
	if id == s.Increment { // decrease increment value
		s.Increment--
	}
	return nil
}

func (s *UserStorage) Update(id int, name string) error {
	if val, found := s.List[id]; !found { // user not found
		return ErrUserNotFound
	} else {
		val.DisplayName = name
		s.List[id] = val
	}
	return nil
}

func (s *UserStorage) Save() error {
	b, err := json.MarshalIndent(&s.UserStore, "", "\t")
	if err != nil {
		return errors.Wrap(err, "can't marshal data of store")
	}
	if err = ioutil.WriteFile(s.path, b, fs.ModePerm); err != nil {
		return errors.Wrap(err, "can't save data of store into file")
	}
	return nil
}

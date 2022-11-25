package service

import (
	"refactoring/internal/db"
	"refactoring/internal/exchange/requests"
	"refactoring/internal/user"
)

type Service struct {
	userStorage *db.UserStorage
}

func CreateNewService(us *db.UserStorage) *Service {
	return &Service{
		userStorage: us,
	}
}

func (s *Service) GetUsers() *user.UserStore {
	return s.userStorage.GetUsers()
}

func (s *Service) GetUser(id int) (*user.User, error) {
	return s.userStorage.GetUser(id)
}

func (s *Service) AddUser(params requests.CreateUserRequest) int {
	u := user.CreateUser(params.DisplayName, params.Email)
	return s.userStorage.Add(u)
}

func (s *Service) DeleteUser(id int) error {
	return s.userStorage.Delete(id)
}

func (s *Service) UpdateUser(id int, params requests.UpdateUserRequest) error {
	return s.userStorage.Update(id, params.DisplayName)
}

func (s *Service) Save() error {
	return s.userStorage.Save()
}

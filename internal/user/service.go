package user

import (
	"context"
	"fmt"

	"m31/internal/apperror"
	"m31/pkg/logging"
)

type Service struct {
	Storage Storage
	Logger  *logging.Logger
}

func NewService(logger *logging.Logger, storage Storage) *Service {
	return &Service{
		Storage: storage,
		Logger:  logger,
	}
}

func (s *Service) Create(ctx context.Context, CreateUserDTO *CreateUserDTO) (id string, err error) {
	user := NewUser(CreateUserDTO)

	users, err := s.Storage.FindAll(ctx)
	if err != nil {
		return
	}

	for _, u := range users {
		if u.Name == user.Name {
			return "", apperror.NewAppError(nil, fmt.Sprintf("Пользователь с именем %s уже существует!", user.Name), "", "US-000003")
		}
	}

	return s.Storage.Create(ctx, user)
}

func (s *Service) FindAll(ctx context.Context) (u []*User, err error) {
	return s.Storage.FindAll(ctx)
}

func (s *Service) GetUserByID(ctx context.Context, id string) (u *User, err error) {
	return s.Storage.FindOne(ctx, id)
}

func (s *Service) UpdateUser(ctx context.Context, id string, UpdateUserAgeDTO *UpdateUserAgeDTO) (u *User, err error) {
	u, err = s.Storage.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}

	u.Age = UpdateUserAgeDTO.Age

	return u, s.Storage.Update(ctx, u)
}

func (s *Service) DeleteUser(ctx context.Context, id string) (u *User, err error) {
	u, err = s.Storage.FindOne(ctx, id)

	if err != nil {
		return nil, err
	}

	if err := s.Storage.DeleteFriend(ctx, id); err != nil {
		return nil, err
	}

	return u, s.Storage.Delete(ctx, id)
}

func (s *Service) MakeFriends(ctx context.Context, CreateFriendDTO *CreateFriendDTO) (message string, err error) {

	targetUser, err := s.Storage.FindOne(ctx, CreateFriendDTO.TargetId)
	if err != nil {
		return "", err
	}

	sourceUser, err := s.Storage.FindOne(ctx, CreateFriendDTO.SourceId)
	if err != nil {
		return "", err
	}

	for _, user := range sourceUser.Friends {
		if user == targetUser.ID {
			return "", apperror.NewAppError(nil, fmt.Sprintf("Пользователь %s уже является другом %s", sourceUser.Name, targetUser.Name), "", "US-000003")
		}
	}

	sourceUser.Friends = append(sourceUser.Friends, targetUser.ID)
	if err := s.Storage.Update(ctx, sourceUser); err != nil {
		return "", err
	}

	targetUser.Friends = append(targetUser.Friends, sourceUser.ID)
	if err := s.Storage.Update(ctx, targetUser); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s и %s теперь друзья", sourceUser.Name, targetUser.Name), nil

}

func (s *Service) GetUserFriends(ctx context.Context, id string) (u []*User, err error) {
	user, err := s.Storage.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}

	u = make([]*User, 0, len(user.Friends))
	for _, userID := range user.Friends {
		friend, err := s.Storage.FindOne(ctx, userID)
		if err != nil {
			return nil, err
		}
		u = append(u, friend)
	}

	return
}

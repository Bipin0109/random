// package service

// import (
// 	"context"
// 	"skyfox/bookings/model"
// )

// type UserRepository interface {
// 	FindByUsername(ctx context.Context, username string) (model.User, error)
// 	Create(ctx context.Context, user *model.User) error

// }

// type userService struct {
// 	userRepo UserRepository
// }

// func NewUserService(userRepository UserRepository) *userService {
// 	return &userService{
// 		userRepo: userRepository,
// 	}
// }

// func (s *userService) UserDetails(ctx context.Context, username string) (model.User, error) {

// 	user, err := s.userRepo.FindByUsername(ctx, username)
// 	if err != nil {
// 		return model.User{}, err
// 	}
// 	return user, nil
// }
package service

import (
    "context"
    "skyfox/bookings/model"
    ae "skyfox/error"
)

type UserRepository interface {
    FindByUsername(ctx context.Context, username string) (model.User, error)
    Create(ctx context.Context, user *model.User) error
    GetPasswordHistory(ctx context.Context, username string) (model.PasswordHistory, error)
    SavePasswordHistory(ctx context.Context, passwordHistory *model.PasswordHistory) error
}

type userService struct {
    userRepo UserRepository
}

func NewUserService(userRepository UserRepository) *userService {
    return &userService{
        userRepo: userRepository,
    }
}

func (s *userService) UserDetails(ctx context.Context, username string) (model.User, error) {
    user, err := s.userRepo.FindByUsername(ctx, username)
    if err != nil {
        return model.User{}, err
    }
    return user, nil
}

func (s *userService) ChangePassword(ctx context.Context, username, newPassword string) error {
    user, err := s.userRepo.FindByUsername(ctx, username)
    if err != nil {
        return err
    }

    passwordHistory, err := s.userRepo.GetPasswordHistory(ctx, username)
    if err != nil {
        return err
    }

    // Check if the new password matches any of the previous passwords
    if (passwordHistory.PreviousPassword1 != nil && newPassword == *passwordHistory.PreviousPassword1) ||
        (passwordHistory.PreviousPassword2 != nil && newPassword == *passwordHistory.PreviousPassword2) ||
        (passwordHistory.PreviousPassword3 != nil && newPassword == *passwordHistory.PreviousPassword3) {
        return ae.UnProcessableError("PasswordReuseError", "New password cannot be the same as any of the previous passwords", nil)
    }

    // Update the password history
    passwordHistory.PreviousPassword3 = passwordHistory.PreviousPassword2
    passwordHistory.PreviousPassword2 = passwordHistory.PreviousPassword1
    passwordHistory.PreviousPassword1 = &user.Password

    // Save the updated password history
    err = s.userRepo.SavePasswordHistory(ctx, &passwordHistory)
    if err != nil {
        return err
    }

    // Update the user's password
    user.Password = newPassword
    err = s.userRepo.Create(ctx, &user)
    if err != nil {
        return err
    }

    return nil
}

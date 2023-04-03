package usecase

import (
	"dvigus_task/internal/pkg/logger"
)

type IUserUsecase interface {
	// GetUserByCookie(ctx context.Context, cookie string) (models.User, int, error)
}

type userUsecase struct {
	logger *logger.Logger
}

func NewUserUsecase(logger *logger.Logger) IUserUsecase {
	return &userUsecase{
		logger: logger,
	}
}

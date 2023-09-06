package users

import (
	"context"
	"sync"
	"time"

	"github.com/bjarke-xyz/uber-clone-backend/internal/core"
)

type UserService struct {
	userRepo UserRepository
}

func NewService(userRepo UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, userID string) (User, error) {
	user, err := s.userRepo.GetByUserID(ctx, userID)
	return user, core.WrapErr(err)
}

func (s *UserService) GetSimulatedUsers(ctx context.Context) ([]User, error) {
	users, err := s.userRepo.GetSimulatedUsers(ctx)
	return users, core.WrapErr(err)
}

type PostUserLogInput struct {
	Tag     string `json:"tag"`
	Message string `json:"message"`
}
type UserLogEvent struct {
	UserID    int64     `json:"userId"`
	Message   string    `json:"message"`
	Tag       string    `json:"tag"`
	Timestamp time.Time `json:"timestamp"`
}

func (s *UserService) AddUserLog(ctx context.Context, userID string, input *PostUserLogInput) (UserLogEvent, error) {
	user, err := s.userRepo.GetByUserID(ctx, userID)
	if err != nil {
		return UserLogEvent{}, core.Errorw(core.EINTERNAL, err)
	}
	if !user.Simulated {
		return UserLogEvent{}, core.Errorf(core.EINVALID, "only simulated users can POST logs")
	}

	userLogEvent := UserLogEvent{
		UserID:    user.ID,
		Tag:       input.Tag,
		Message:   input.Message,
		Timestamp: time.Now().UTC(),
	}
	go storeUserLog(userLogEvent)
	return userLogEvent, nil
}

func (s *UserService) GetRecentLogs(ctx context.Context) ([]UserLogEvent, error) {
	return recentUserLogs, nil
}

var recentUserLogs = make([]UserLogEvent, 0)
var recentUserLogsLock sync.RWMutex
var maxRecentLogs = 100

func storeUserLog(event UserLogEvent) {
	recentUserLogsLock.Lock()
	defer recentUserLogsLock.Unlock()
	// prepend
	recentUserLogs = append([]UserLogEvent{event}, recentUserLogs...)
	if len(recentUserLogs) > maxRecentLogs {
		recentUserLogs = recentUserLogs[:maxRecentLogs]
	}
}

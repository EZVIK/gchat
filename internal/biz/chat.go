package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type ChatRoom struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	CreateUserId int64  `json:"create_user_id"`
}

type ChatUserList struct {
	ID      int64   `json:"id"`
	UserIds []int64 `json:"user_ids"`
}

type UserChat struct {
	ID     int64 `json:"id"`
	UserId int64 `json:"user_id"`
	ChatId int64 `json:"chat_id"`
}

type ChatUsecase struct {
	repo ChatRepo
	log  *log.Helper
}

type ChatRepo interface {
	CreateChat(context.Context, *ChatRoom) (*int64, error)
	RemoveChat(context.Context, *ChatRoom) (*int64, error)
	JoinChatRoom(context.Context, *[]UserChat) error
	LeftChatRoom(context.Context, *ChatUserList) error
	QueryChatUserInfo(context.Context, *int64) (*ChatUserList, error)
}

// NewChatUsecase .
func NewChatUsecase(repo ChatRepo, logger log.Logger) *ChatUsecase {
	return &ChatUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *ChatUsecase) CreateChatRoom(ctx context.Context, u *ChatRoom) (*int64, error) {

	i, err := uc.repo.CreateChat(ctx, u)
	if err != nil {
		return nil, err
	}

	err = uc.repo.JoinChatRoom(ctx, &[]UserChat{{ChatId: *i, UserId: u.CreateUserId}})
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (uc *ChatUsecase) RemoveChatRoom(ctx context.Context, u *ChatRoom) (*int64, error) {
	return nil, nil
}

func (uc *ChatUsecase) JoinChatRoom(ctx context.Context, uclist *[]UserChat) error {
	return uc.repo.JoinChatRoom(ctx, uclist)
}

func (uc *ChatUsecase) LeaveChatRoom(ctx context.Context, lcr *ChatUserList) error {
	return uc.repo.LeftChatRoom(ctx, lcr)
}

func (uc *ChatUsecase) QueryChatUserIds(ctx context.Context, chatId int64) (*ChatUserList, error) {
	return uc.repo.QueryChatUserInfo(ctx, &chatId)
}

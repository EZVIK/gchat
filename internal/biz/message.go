package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type Message struct {
	ID        int64  `json:"id"`
	UserId    int64  `json:"user_id"`
	ChatId    int64  `json:"chat_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type MessagePackages struct {
	Total int                 `json:"total"`
	Data  map[int64][]Message `json:"data"`
}

type ReceiveMsg struct {
	MsgId   int64   `json:"msg_id"`
	ChatId  int64   `json:"chat_id"`
	UserIds []int64 `json:"user_ids"`
}

type MessageRepo interface {
	SendMsg(context.Context, *Message) (*int64, error)

	SavingGroupMsg(context.Context, *ReceiveMsg) (int, error)

	QueryUnreadChatIds(ctx context.Context, UserId int64) ([]int64, error)

	QueryUnreadMsg(ctx context.Context, userId int64) ([]int64, error)

	QueryMsgByIds(ctx context.Context, ids []int64) (*[]Message, error)
}

type MessageUsecase struct {
	repo MessageRepo
	log  *log.Helper
}

func NewMessageUsecase(repo MessageRepo, logger log.Logger) *MessageUsecase {
	return &MessageUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (m *MessageUsecase) SendMsg(ctx context.Context, msg *Message, ids []int64) (int, error) {

	// 1. save original msg
	i, err := m.repo.SendMsg(ctx, msg)
	if err != nil {
		return 0, err
	}

	// 2. Query related users by chat room id
	// ids

	// 3. Create msg and user id bind
	r := &ReceiveMsg{
		MsgId:   *i,
		ChatId:  msg.ChatId,
		UserIds: ids,
	}

	count, err := m.repo.SavingGroupMsg(ctx, r)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *MessageUsecase) ReceiveMsg(ctx context.Context, userId int64) (*[]Message, error) {

	msgIds, err := m.repo.QueryUnreadMsg(ctx, userId)
	if err != nil {
		return nil, err
	}

	msgList, err := m.repo.QueryMsgByIds(ctx, msgIds)
	if err != nil {
		return nil, err
	}

	return msgList, nil
}

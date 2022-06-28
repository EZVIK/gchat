package data

import (
	"context"
	"gchat/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type Chat struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	CreateUserId int64  `json:"create_user_id"`
}

type UserChatRela struct {
	ID     int64 `json:"id"`
	ChatId int64 `json:"chat_id"`
	UserId int64 `json:"user_id"`
	Status int64 `json:"status"`
}

type chatRepo struct {
	data *Data
	log  *log.Helper
}

// NewChatRepo .
func NewChatRepo(data *Data, logger log.Logger) biz.ChatRepo {
	return &chatRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (g *chatRepo) CreateChat(ctx context.Context, req *biz.ChatRoom) (*int64, error) {

	c := &Chat{Name: req.Name, CreateUserId: req.CreateUserId}
	tx := g.data.db.Create(c)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &c.ID, nil
}

func (g *chatRepo) RemoveChat(ctx context.Context, req *biz.ChatRoom) (*int64, error) {

	c := &Chat{ID: req.ID, Name: req.Name, CreateUserId: req.CreateUserId}
	tx := g.data.db.Delete(c)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &c.ID, nil
}

func (g *chatRepo) JoinChatRoom(ctx context.Context, uc *[]biz.UserChat) error {

	userChatRelaList := make([]UserChatRela, len(*uc))

	for i, v := range *uc {
		userChatRelaList[i] = UserChatRela{
			ChatId: v.ChatId,
			UserId: v.UserId,
			Status: 0,
		}
	}

	tx := g.data.db.Create(&userChatRelaList)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (g *chatRepo) LeftChatRoom(ctx context.Context, lcr *biz.ChatUserList) error {

	tx := g.data.db.Model(&UserChatRela{}).Where("chat_id = ? and user_id in ?", lcr.ID, lcr.UserIds).Update("status", 1)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (g *chatRepo) QueryChatUserInfo(ctx context.Context, chatId *int64) (*biz.ChatUserList, error) {

	ul := make([]UserChatRela, 0)
	tx := g.data.db.WithContext(ctx).Model(&UserChatRela{}).Where("chat_id = ?", *chatId).Find(&ul)

	if tx.Error != nil {
		return nil, tx.Error
	}

	ids := make([]int64, len(ul))
	for i, ucr := range ul {
		ids[i] = ucr.UserId
	}

	return &biz.ChatUserList{
		ID:      *chatId,
		UserIds: ids,
	}, nil

}

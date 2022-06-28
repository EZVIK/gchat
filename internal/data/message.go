package data

import (
	"context"
	"gchat/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type Message struct {
	ID      int64  `json:"id"`
	UserId  int64  `json:"user_id"`
	ChatId  int64  `json:"chat_id"`
	Content string `json:"content"`
	gorm.Model
}

type MessageGrouping struct {
	ID     int64 `json:"id"`
	ChatId int64 `json:"chat_id"`
	MsgId  int64 `json:"msg_id"`
	UserId int64 `json:"user_id"`
	Status int   `json:"status"`
}

type messageRepo struct {
	data *Data
	log  *log.Helper
}

func NewMessageRepo(data *Data, logger log.Logger) biz.MessageRepo {
	return &messageRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (g *messageRepo) SendMsg(ctx context.Context, req *biz.Message) (*int64, error) {

	m := &Message{UserId: req.UserId, ChatId: req.ChatId, Content: req.Content}

	tx := g.data.db.Model(&Message{}).Create(m)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &m.ID, nil
}

func (g *messageRepo) SavingGroupMsg(ctx context.Context, req *biz.ReceiveMsg) (int, error) {

	sms := make([]MessageGrouping, len(req.UserIds))
	for i, v := range req.UserIds {
		sms[i] = MessageGrouping{MsgId: req.MsgId, UserId: v, Status: 0}
	}

	tx := g.data.db.Create(&sms)
	if tx.Error != nil {
		return -1, tx.Error
	}

	return len(sms), nil
}

func (g *messageRepo) QueryUnreadChatIds(ctx context.Context, UserId int64) ([]int64, error) {

	var ids []int64

	return ids, nil
}

func (g *messageRepo) QueryUnreadMsg(ctx context.Context, userId int64) ([]int64, error) {

	var ids []int64
	tx := g.data.db.Model(&MessageGrouping{}).Where("user_id = ? and status = 0", userId).Select("msg_id").Find(&ids)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return ids, nil
}

func (g *messageRepo) QueryMsgByIds(ctx context.Context, ids []int64) (*[]biz.Message, error) {

	var msgList []Message
	tx := g.data.db.Model(&Message{}).Where("id in ?", ids).Find(&msgList)
	if tx.Error != nil {
		return nil, tx.Error
	}

	var bm []biz.Message
	for _, v := range msgList {
		bm = append(bm, biz.Message{
			ID:        v.ID,
			UserId:    v.UserId,
			ChatId:    v.ChatId,
			Content:   v.Content,
			CreatedAt: v.Model.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &bm, nil
}

// queryUnreadChatIds?status=0&groupBy=chat_id

// queryUnreadMsgByUserAndChatId?user_id=1&chat_id=1&status=0
// func (g *messageRepo) ReceiveMsg(ctx context.Context, UserId int64) (*biz.MessagePackages, error) {

// 	g.data.db.Select("user_id = ?", userId).Model(&SendedMsg{}).Where("status = ?", 0)

// 	return nil, nil
// }

// func (g *messageRepo) SendMsg(context.Context, *Message) (*int64, error)
// func (g *messageRepo) SavingGroupMsg(context.Context, *ReceiveMsg) (int, error)
// func (g *messageRepo) QueryUnreadChatIds(ctx context.Context, UserId int64) ([]int64, error)
// func (g *messageRepo) QueryUnreadMsg(ctx context.Context, userId, ChatId int64) ([]int64, error)
// func (g *messageRepo) QueryMsgByIds(ctx context.Context, ids []int64) ([]Message, error)

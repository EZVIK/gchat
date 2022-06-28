package service

import (
	"context"
	v1 "gchat/api/gchat/v1"
	"gchat/internal/biz"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
)

func (s *GchatService) SendMsg(ctx context.Context, req *v1.SendMsgRequest) (*v1.SendMsgReply, error) {

	msg := biz.Message{
		ChatId:    req.ChatId,
		UserId:    req.UserId,
		Content:   req.Content,
		CreatedAt: time.Now().Format("2006-01-02 03:04:05"),
	}

	cul, err := s.cc.QueryChatUserIds(ctx, req.GetChatId())
	if err != nil {
		return nil, errors.Newf(500, "QueryChatUserIds failed", err.Error())
	}

	_, err = s.mc.SendMsg(ctx, &msg, cul.UserIds)
	if err != nil {
		return nil, errors.Newf(500, "SendMsg failed", err.Error())
	}
	return &v1.SendMsgReply{}, nil
}

func (s *GchatService) ReceiveMsg(ctx context.Context, req *v1.ReceiveMsgRequest) (*v1.ReceiveMsgReply, error) {

	msgList, err := s.mc.ReceiveMsg(ctx, req.UserId)
	if err != nil {
		return nil, errors.Newf(500, "ReceiveMsg failed", err.Error())
	}

	msgMap := make(map[int64]*v1.ReceiveMsgs, 0)
	result := make([]*v1.ReceiveMsgs, 0)
	for _, msg := range *msgList {

		_, ok := msgMap[msg.ChatId]

		if !ok {
			msgMap[msg.ChatId] = &v1.ReceiveMsgs{
				ChatId:   msg.ChatId,
				Messages: make([]*v1.MessageObject, 0),
			}
		}

		msgMap[msg.ChatId].Messages = append(msgMap[msg.ChatId].Messages, &v1.MessageObject{
			Id:         msg.ID,
			UserId:     msg.UserId,
			Content:    msg.Content,
			ChatId:     msg.ChatId,
			CreateTime: msg.CreatedAt,
		})

	}

	for _, msgList := range msgMap {
		result = append(result, msgList)
	}

	return &v1.ReceiveMsgReply{
		Data: result,
	}, nil
}

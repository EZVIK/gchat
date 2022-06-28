package service

import (
	"context"
	"errors"
	v1 "gchat/api/gchat/v1"
	"gchat/internal/biz"
)

func (s *GchatService) CreateChat(ctx context.Context, req *v1.CreateChatRequest) (*v1.CreateChatReply, error) {

	i, err := s.cc.CreateChatRoom(ctx, &biz.ChatRoom{
		Name:         req.Name,
		CreateUserId: req.CreateUserId,
	})

	if err != nil {
		return nil, errors.New("create chat room failed: " + err.Error())
	}

	return &v1.CreateChatReply{
		Id: *i,
	}, nil
}

func (s *GchatService) RemoveChat(ctx context.Context, req *v1.RemoveChatRequest) (*v1.RemoveChatReply, error) {
	return &v1.RemoveChatReply{}, nil
}

func (s *GchatService) JoinChat(ctx context.Context, req *v1.JoinChatRequest) (*v1.JoinChatReply, error) {

	ucList := make([]biz.UserChat, len(req.UserIds))
	for i, v := range req.UserIds {
		ucList[i] = biz.UserChat{
			UserId: v,
			ChatId: req.ChatId,
		}
	}

	err := s.cc.JoinChatRoom(ctx, &ucList)
	if err != nil {
		return nil, errors.New("join chat room failed: " + err.Error())
	}

	return &v1.JoinChatReply{}, nil
}

func (s *GchatService) LeaveChat(ctx context.Context, req *v1.LeaveChatRequest) (*v1.LeaveChatReply, error) {

	ucl := []biz.UserChat{{
		UserId: req.UserId,
		ChatId: req.ChatId,
	}}

	err := s.cc.JoinChatRoom(ctx, &ucl)
	if err != nil {
		return nil, errors.New("leave chat room failed: " + err.Error())
	}

	return &v1.LeaveChatReply{}, nil
}

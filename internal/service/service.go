package service

import (
	v1 "gchat/api/gchat/v1"
	"gchat/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
// 实现了 api 定义的服务层，类似 DDD 的 application 层
// 处理 DTO 到 biz 领域实体的转换(DTO -> DO)
// 同时协同各类 biz 交互，但是不应处理复杂逻辑
var ProviderSet = wire.NewSet(NewGchatService)

type GchatService struct {
	v1.UnimplementedGchatServer
	uc  *biz.UserUsecase
	cc  *biz.ChatUsecase
	mc  *biz.MessageUsecase
	log *log.Helper
}

func NewGchatService(uc *biz.UserUsecase, cc *biz.ChatUsecase, mc *biz.MessageUsecase, logger log.Logger) *GchatService {
	return &GchatService{uc: uc, cc: cc, mc: mc, log: log.NewHelper(logger)}
}

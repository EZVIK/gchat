package server

import (
	"context"
	v1 "gchat/api/gchat/v1"
	"gchat/internal/conf"
	"gchat/internal/service"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"

	jwt2 "github.com/golang-jwt/jwt/v4"
)

func NewWhiteListMatcher() selector.MatchFunc {

	whiteList := make(map[string]struct{})
	whiteList["/gchat.api.gchat.v1.Gchat/Login"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, auth *conf.Auth, greeter *service.GchatService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),

			logging.Server(logger),

			//jwt.Server(func(token *jwtv4.Token) (interface{}, error) {
			//	return []byte(auth.Key), nil
			//}),

			selector.Server(

				jwt.Server(
					func(token *jwt2.Token) (interface{}, error) {
						return []byte(auth.Key), nil
					},
					jwt.WithSigningMethod(jwt2.SigningMethodHS256)),
			).
				Match(NewWhiteListMatcher()).
				Build(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	h := openapiv2.NewHandler()
	srv.HandlePrefix("/q/", h)
	v1.RegisterGchatHTTPServer(srv, greeter)
	return srv
}

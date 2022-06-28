package service

import (
	"context"
	"errors"
	v1 "gchat/api/gchat/v1"
	"gchat/internal/biz"
	"github.com/golang-jwt/jwt/v4"
)

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// TODO: Validate Params
// TODO: Trace ID
// TODO: Authorization
// TODO: Middleware
func (s *GchatService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	i, err := s.uc.Login(ctx, &biz.LoginRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	})

	if err != nil {
		return nil, errors.New("shit man")
	}

	return &v1.LoginReply{
		Id:          i.ID,
		Username:    i.Username,
		AccessToken: i.AccessToken,
	}, nil
}

//func ()  {
//
//	// Declare the expiration time of the token
//	// here, we have kept it as 5 minutes
//	expirationTime := time.Now().Add(2 * time.Hour)
//	// Create the JWT claims, which includes the username and expiry time
//	claims := &Claims{
//		Username: req.Username,
//		RegisteredClaims: jwt.RegisteredClaims{
//			// In JWT, the expiry time is expressed as unix milliseconds
//			ExpiresAt: &jwt.NumericDate{
//				Time: expirationTime,
//			},
//		},
//	}
//
//	// Declare the token with the algorithm used for signing, and the claims
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	// Create the JWT string
//	tokenString, err := token.SignedString("1234")
//	if err != nil {
//
//	}
//}

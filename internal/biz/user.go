package biz

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	v1 "gchat/api/gchat/v1"
	"gchat/internal/conf"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

type AddUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserQuery struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRepo interface {
	AddUser(context.Context, *AddUser) (*UserInfo, error)
	QueryUser(context.Context, *UserQuery) (*UserInfo, error)
}

type GenerateBasic struct {
	Client        string
	UserID        string
	UserNo        int
	CreateAt      time.Time
	CustomData    string
	Scope         string
	State         string
	QRCodeStatus  int
	DeviceId      string
	DeviceName    string
	DeviceType    string
	DeviceModel   string
	DeviceVersion string
}

// UserUsecase is a Greeter usecase.
type UserUsecase struct {
	jwtGenerator func(data *GenerateBasic, isGenRefresh bool) (string, string, error)
	repo         UserRepo
	log          *log.Helper
}

// NewUserUsecase new a User usecase.
func NewUserUsecase(auth func(data *GenerateBasic, isGenRefresh bool) (string, string, error), repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{jwtGenerator: auth, repo: repo, log: log.NewHelper(logger)}
}

func NewJwtGenerator(auth *conf.Auth) func(data *GenerateBasic, isGenRefresh bool) (string, string, error) {
	return func(data *GenerateBasic, isGenRefresh bool) (string, string, error) {
		priKey, err := GetPrivateKey(auth.Key)
		if err != nil {
			return "", "", err
		}
		claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{

			Issuer:    "gchat",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
			Subject:   data.UserID,
			Audience: jwt.ClaimStrings{
				"client1",
			},
		})
		signedString, err := claims.SignedString(priKey)
		if err != nil {
			return "", "", err
		}

		// TODO REFRESH TOKEN
		return signedString, "", nil
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserInfo
	AccessToken string `json:"access_token"`
}

func (uc *UserUsecase) Login(ctx context.Context, u *LoginRequest) (ur *LoginResponse, err error) {
	uc.log.WithContext(ctx).Infof("CreateGreeter: %v", u.Username)

	var userInfo *UserInfo

	// 查询用户
	userInfo, err = uc.repo.QueryUser(ctx, &UserQuery{
		Username: u.Username,
		Password: u.Password,
	})

	// 未注册时自动添加用户
	if err == gorm.ErrRecordNotFound {
		uc.log.WithContext(ctx).Infof("UnRegister user: %s", u.Username)

		userInfo, err = uc.repo.AddUser(ctx, &AddUser{
			Username: u.Username,
			Password: u.Password,
		})
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, errors.Errorf(200001, v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
	}

	ac, _, err := uc.jwtGenerator(&GenerateBasic{UserID: strconv.Itoa(int(userInfo.ID))}, false)
	if err != nil {
		return nil, errors.Errorf(200002, v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
	}

	return &LoginResponse{
		UserInfo:    *userInfo,
		AccessToken: ac,
	}, nil
}

func GetPrivateKey(key string) (*rsa.PrivateKey, error) {
	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	priKey, err := jwt.ParseRSAPrivateKeyFromPEM(decoded)
	return priKey, err
}

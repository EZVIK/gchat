package data

import (
	"context"
	"gchat/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type User struct {
	ID       int64  `gorm:"primaryKey" json:"id"`
	Username string `gorm:"type:varchar(100);not null"`
	//Email    *string `gorm:"type:varchar(200);not null"`
	Password string `gorm:"type:varchar(200);not null" json:"password"`
}

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (g *userRepo) AddUser(ctx context.Context, req *biz.AddUser) (*biz.UserInfo, error) {

	du := User{
		Username: req.Username,
		Password: req.Password,
	}
	tx := g.data.db.Create(&du)
	if tx.Error != nil {
		return nil, tx.Error
	}
	uq := &biz.UserInfo{
		ID:       du.ID,
		Username: du.Username,
		Password: req.Password,
	}

	return uq, nil
}

func (g *userRepo) QueryUser(ctx context.Context, req *biz.UserQuery) (*biz.UserInfo, error) {

	u := User{}

	tx := g.data.db.Model(&User{}).Where("username = ? and password = ?", req.Username, req.Password).First(&u)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &biz.UserInfo{
		ID:       u.ID,
		Username: u.Username,
	}, nil
}

package user

import (
	"context"
	"fmt"

	"pdf-service/user/api/internal/svc"
	"pdf-service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) error {
	user, err := l.svcCtx.UserModel.FindByName(l.ctx, req.Phone)
	if err != nil {
		return err
	}

	if user.Pwd == req.Pwd {
		fmt.Println("登录成功")
	}

	return nil
}

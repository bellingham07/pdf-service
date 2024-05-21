package user

import (
	"context"
	"pdf-service/user/api/internal/model"

	"pdf-service/user/api/internal/svc"
	"pdf-service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) error {
	user := &model.User{
		Phone: req.Phone,
		Pwd:   req.Pwd,
	}
	if err := l.svcCtx.UserModel.Insert(l.ctx, l.svcCtx.DB, user); err != nil {
		return err
	}
	return nil
}

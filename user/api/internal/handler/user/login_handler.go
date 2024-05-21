package user

import (
	"github.com/HiBugEnterprise/gotools/errorx"
	"github.com/HiBugEnterprise/gotools/httpc"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"pdf-service/user/api/internal/logic/user"
	"pdf-service/user/api/internal/svc"
	"pdf-service/user/api/internal/types"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpc.RespError(w, r, errorx.BadRequest("%s:%s", errorx.CodeInvalidParams.Msg(), err.Error()).Show())
			return
		}
		l := user.NewLoginLogic(r.Context(), svcCtx)
		err := l.Login(&req)
		if err != nil {
			httpc.RespError(w, r, err)
		} else {
			httpc.RespSuccess(r.Context(), w, nil)
		}
	}
}

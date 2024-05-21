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

func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			httpc.RespError(w, r, errorx.BadRequest("%s:%s", errorx.CodeInvalidParams.Msg(), err.Error()).Show())
			return
		}
		l := user.NewRegisterLogic(r.Context(), svcCtx)
		err := l.Register(&req)
		if err != nil {
			httpc.RespError(w, r, err)
		} else {
			httpc.RespSuccess(r.Context(), w, nil)
		}
	}
}

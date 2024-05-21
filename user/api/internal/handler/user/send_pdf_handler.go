package user

import (
	"github.com/HiBugEnterprise/gotools/httpc"
	"net/http"

	"pdf-service/user/api/internal/logic/user"
	"pdf-service/user/api/internal/svc"
)

func SendPdfHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewSendPdfLogic(r.Context(), svcCtx, r)
		err := l.SendPdf()
		if err != nil {
			httpc.RespError(w, r, err)
		} else {
			httpc.RespSuccess(r.Context(), w, nil)
		}
	}
}

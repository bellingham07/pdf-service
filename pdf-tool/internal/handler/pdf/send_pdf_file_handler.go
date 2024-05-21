package pdf

import (
	"github.com/HiBugEnterprise/gotools/httpc"
	"net/http"

	"pdf-service/pdf-tool/internal/logic/pdf"
	"pdf-service/pdf-tool/internal/svc"
)

func SendPdfFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := pdf.NewSendPdfFileLogic(r.Context(), svcCtx, w, r)
		resp, err := l.SendPdfFile()
		if err != nil {
			httpc.RespError(w, r, err)
		} else {
			httpc.RespSuccess(r.Context(), w, resp)
		}
	}
}

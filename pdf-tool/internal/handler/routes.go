// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	pdf "pdf-service/pdf-tool/internal/handler/pdf"
	"pdf-service/pdf-tool/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/send",
				Handler: pdf.SendPdfFileHandler(serverCtx),
			},
		},
		rest.WithPrefix("/pdf"),
	)
}

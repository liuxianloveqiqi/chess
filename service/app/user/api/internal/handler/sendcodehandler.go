package handler

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"user/api/internal/logic"
	"user/api/internal/svc"
	"user/api/internal/types"
	"user/common/response"
)

func SendcodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterByPhoneRep
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewSendcodeLogic(r.Context(), svcCtx)
		resp, err := l.Sendcode(&req)
		response.Response(w, resp, err) //②

	}
}

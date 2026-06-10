// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package agent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gateway/internal/logic/agent"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AgentChatStreamHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AgentChatReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if strings.TrimSpace(req.Question) == "" {
			response.Response(w, nil, response.ErrorBadRequest("问题不能为空"))
			return
		}

		// Buffer size of 16 is chosen as a reasonable default to balance throughput and memory usage.
		// You can change this based on your application's needs.
		// if your go-zero version less than 1.8.1, you need to add 3 lines below.
		// w.Header().Set("Content-Type", "text/event-stream")
		// w.Header().Set("Cache-Control", "no-cache")
		// w.Header().Set("Connection", "keep-alive")
		// 构建channel，用于接收和处理sse流式数据
		client := make(chan *types.AgentChatStreamChunk, 16)

		l := agent.NewAgentChatStreamLogic(r.Context(), svcCtx)
		threading.GoSafeCtx(r.Context(), func() {
			defer close(client)
			err := l.AgentChatStream(&req, client)
			if err != nil {
				logc.Errorw(r.Context(), "AgentChatStreamHandler", logc.Field("error", err))
				return
			}
		})

		for {
			// 监听channel
			select {
			case data, ok := <-client:
				if !ok {
					return
				}
				output, err := json.Marshal(data)
				if err != nil {
					logc.Errorw(r.Context(), "AgentChatStreamHandler", logc.Field("error", err))
					continue
				}

				// 这里是因为需要封装SSE的固定格式。
				if _, err := fmt.Fprintf(w, "data: %s\n\n", string(output)); err != nil {
					logc.Errorw(r.Context(), "AgentChatStreamHandler", logc.Field("error", err))
					return
				}
				// 强制将缓冲区的数据发送。
				if flusher, ok := w.(http.Flusher); ok {
					flusher.Flush()
				}
			case <-r.Context().Done():
				return
			}
		}
	}
}

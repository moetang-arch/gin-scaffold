package handlers

import (
	"github.com/gin-gonic/gin"

	"import.moetang.info/go/nekoq-api/ctx"
	"import.moetang.info/go/nekoq-api/entrance"
	"time"
)

const (
	REQUEST_ID_HEADER_KEY = "X-Req-Id"
	RPC_ID_HEADER_KEY     = "X-Req-Rpc-Id"

	REQUEST_CONTEXT_KEY = "__gin_request_context__"
)

func EntranceInit(c *gin.Context) {
	requestId := c.Request.Header.Get(REQUEST_ID_HEADER_KEY)
	if requestId == "" {
		requestId = entrance.GenerateRequestId()
	}
	rpcId := c.Request.Header.Get(RPC_ID_HEADER_KEY)
	if rpcId == "" {
		rpcId = "0"
	}
	now := time.Now()
	cur := now.Unix()*1000 + now.UnixNano()/1000/1000%1000
	timeLimit := ctx.TimeLimit{
		CurrentTimeMillis: uint64(cur),
		DeadlineMillis:    uint64(cur + 10*1000), //10s
	}
	ttl := ctx.TTL{
		CurrentTTL: 0,
		MaxTTL:     16,
	}
	reqCtx := ctx.NewContext(requestId, rpcId, timeLimit, ttl)
	c.Set(REQUEST_CONTEXT_KEY, reqCtx)
}

func GetCtx(c *gin.Context) *ctx.Context {
	result, ok := c.Get(REQUEST_CONTEXT_KEY)
	if ok {
		return result.(*ctx.Context)
	} else {
		return nil
	}
}

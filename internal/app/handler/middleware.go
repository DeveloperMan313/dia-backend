package handler

import "github.com/gin-gonic/gin"

func SetCorsHeaders(gCtx *gin.Context) {
	gCtx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
}

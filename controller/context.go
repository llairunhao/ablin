package controller

import (
	"crypto/md5"
	"fmt"
	"ggfly/dao"
	"ggfly/logger"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var pool sync.Pool

func init() {
	fmt.Println("init")
	pool = sync.Pool{
		New: func() interface{} {
			return &dao.GGResponse{Code: dao.StatusOk, Message: "", Data: nil}
		},
	}
}

func setRequestId(c *gin.Context) {
	num := strconv.FormatInt(time.Now().UnixNano()+rand.Int63(), 16)
	requestId := fmt.Sprintf("%x", md5.Sum([]byte(num)))
	c.Set("requestId", requestId)
}

func getRequestId(c *gin.Context) string {
	return c.GetString("requestId")
}

func tryBind(c *gin.Context, obj interface{}) bool {
	err := c.ShouldBindJSON(obj)
	if err != nil {
		writeFailedResponse(c, dao.StatusParameterError, err)
		return false
	}
	logger.DebugLog("%v", obj)
	return true
}

func writeFailedResponse(c *gin.Context, code dao.GGStatusCode, err error) {
	response := getResponse()
	response.Code = code
	response.Message = err.Error()
	c.JSON(http.StatusOK, response)
	saveResponse(response)
}

func writeSuccessResponse(c *gin.Context, obj interface{}) {
	response := getResponse()
	response.Code = dao.StatusOk
	response.Message = "请求成功"
	response.Data = obj
	c.JSON(http.StatusOK, response)
	saveResponse(response)
}

func getResponse() *dao.GGResponse {
	return pool.Get().(*dao.GGResponse)
}

func saveResponse(r *dao.GGResponse) {
	r.Data = nil
	r.Message = ""
	pool.Put(r)
}

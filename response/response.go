package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/metafiliana/evermos-test/util"
)

func SendResponse(c *gin.Context, code int, message string, data ...interface{}) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func SendResponseWithError(c *gin.Context, err error, data ...interface{}) {
	ret, ok := err.(*util.ErrorWrapper)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": `unknown error`,
			"data":    data,
		})
		return
	}

	c.JSON(ret.ErrorCode, gin.H{
		"code":    ret.ErrorCode,
		"message": ret.ErrorMessage,
		"data":    ret.Err.Error(),
	})

}

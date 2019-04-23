package handlers

import (
	"core"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 申请端口
func EndpointBindLeftHandler(ctx *gin.Context) {
	var params struct {
		Callid   string `binding:"required"`
		Leftip   string `binding:"required"`
		Leftport int    `binding:"required"`
	}

	if err := ctx.BindJSON(&params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	// fmt.Println(params.Callid, params.Leftip, params.Leftport)

	sess, _ := core.GetSessionManage()
	ps, err := sess.BindLeft(params.Callid, params.Leftip, params.Leftport)
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"ip":   sess.GetLocalIp(),
			"port": ps.PortRtp,
		})
	} else {
		fmt.Println(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"test": 100,
		})
	}

}

// 开始工作
func EndpointBindRightHandler(ctx *gin.Context) {
	var params struct {
		Callid    string `binding:"required"`
		Rightip   string `binding:"required"`
		Rightport int    `binding:"required"`
		Media     string `binding:"required"`
	}

	if err := ctx.BindJSON(&params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	sess, _ := core.GetSessionManage()
	err := sess.BindRight(params.Callid, params.Rightip, params.Rightport, params.Media)
	if err != nil {
		fmt.Println(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	sess.Start(params.Callid)

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func EndpointStopHandler(ctx *gin.Context) {
	var params struct {
		Callid string `binding:"required"`
	}

	if err := ctx.BindJSON(&params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	sess, _ := core.GetSessionManage()
	sess.Stop(params.Callid)

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func EndpointUpdateHandler(ctx *gin.Context) {
	var params struct {
		Callid    string `binding:"required"`
		Leftip    string `binding:"required"`
		Leftport  int    `binding:"required"`
		Rightip   string `binding:"required"`
		Rightport int    `binding:"required"`
	}

	if err := ctx.BindJSON(&params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	sess, _ := core.GetSessionManage()
	sess.Stop(params.Callid)

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func Endpoint200Handler(ctx *gin.Context) {
	var params struct {
		Callid string `binding:"required"`
	}

	if err := ctx.BindJSON(&params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	fmt.Println(params.Callid)

	sess, _ := core.GetSessionManage()
	sess.SplitRecord(params.Callid)

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func EndpointStatisHandler(ctx *gin.Context) {
	sess, _ := core.GetSessionManage()
	ctx.JSON(http.StatusOK, gin.H{
		"session.size": sess.Size(),
	})
}

func EndpointPreviewHandler(ctx *gin.Context) {
	sess, _ := core.GetSessionManage()
	sess.Preview()

	//ctx.JSON(http.StatusOK, gin.H{
	//	"msg": "preview",
	//})

	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"msg": "preview",
	})

}

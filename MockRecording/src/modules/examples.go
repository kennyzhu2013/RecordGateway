/*
@Time : 2018/8/21 11:19 
@Author : kenny zhu
@File : examples.go
@Software: GoLand
@Others:
*/
package modules

import (
	"github.com/kennyzhu/go-os/log"
	example "github.com/kennyzhu/go-os/dbservice/proto/example"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
)

type examples struct{
	cl example.PreferencesService
}

// All are run in goroutine
func (s *examples) Preferences(ctx *gin.Context) {
	action := ctx.Param("action")

	switch action {
	case "/GetPreference":
		s.getPreferences( ctx )
	case "/GetPreferencesList":
		s.getPreferencesList( ctx )
	default:
		ctx.JSON(404, map[string]string {
			"message": "Unknown action:" + action,
		})
	}
	log.Debug("Preferences done!")
}

func (s *examples) getPreferences(ctx *gin.Context) {
	log.Info("Received getPreferences http request")

	user,_ := strconv.Atoi ( ctx.DefaultQuery("user", "1") )

	response, err := s.cl.GetPreference(context.TODO(), &example.PreferenceRequest{
		User: int32( user ),
	})

	if err != nil {
		ctx.JSON(500, map[string]string{
			"message": err.Error(),
		})
		log.Error(err)
		return
	}

	/*
	prefersJson,_ := json.Marshal( response.Prefer )
	ctx.JSON(int(response.ResultCode), map[string]string{
		"message": string(prefersJson[:]),
	})*/
	prefersJson,_ := json.Marshal( response.Prefer )
	ctx.JSON(int(response.ResultCode), gin.H{
		"message": string(prefersJson[:]),
	})
	log.Info("getPreferences End:")
}

func (s *examples) getPreferencesList(ctx *gin.Context) {
	log.Info("Received getPreferencesList http request")

	index,_ :=  strconv.Atoi ( ctx.Query("index") )
	limit,_ :=  strconv.Atoi ( ctx.Query("limit") )

	response, err := s.cl.GetPreferencesList(context.TODO(), &example.PreferencesListRequest{
		Index:  int32(index) ,
		Limit: int32(limit),
	})
	if err != nil {
		ctx.JSON(500, map[string]string{
			"message": err.Error(),
		})
		log.Error(err)
		return
	}

	prefersJson,_ := json.Marshal( response.Prefers )
	ctx.JSON(int(response.ResultCode), map[string]string{
		"message": string(prefersJson[:]),
	})
	log.Info("getPreferencesList End:")
}

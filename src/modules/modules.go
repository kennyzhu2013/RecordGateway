/*
@Time : 2018/8/21 11:05 
@Author : kenny zhu
@File : modules.go
@Software: GoLand
@Others:
*/
package modules

import (
	"github.com/gin-gonic/gin"

	"github.com/kennyzhu/go-os/log"
	"encoding/json"
	"io/ioutil"
	"fmt"

	proto "github.com/micro/go-micro/server/debug/proto"
)

var Modules struct{
	Router *gin.Engine
}

// self init...
func init() {
	// gin.SetMode(gin.ReleaseMode)
	Modules.Router = gin.Default()
}

// {"method":"Debug.Health","params":[null],"id":85}..
type serverRequest struct {
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
	ID     *json.RawMessage `json:"id"`
}

type serverResponse struct {
	ID     *json.RawMessage `json:"id"`
	Result interface{}      `json:"result"`
	Error  interface{}      `json:"error"`
}

// Todo: for health beat..
func NoModules(ctx *gin.Context) {
	log.Info("Received NoModules API request")

	b, _ := ioutil.ReadAll(ctx.Request.Body)
	var body serverRequest
	if err := json.Unmarshal(b, &body); err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println( string(b) )
	// fmt.Println( body )
	rsp := &serverResponse{}
	rsp.ID = body.ID

	result := &proto.HealthResponse{}
	result.Status = "ok"
	rsp.Result = result

	prefersJson,_ := json.Marshal( rsp )
	ctx.JSON(200, prefersJson)
	/*
	ctx.Write(prefersJson)
	ctx.JSON(200, map[string]string{
		"message": "No module defined!",
	})
	*/
}

// anything to do..
// run here not go routine...
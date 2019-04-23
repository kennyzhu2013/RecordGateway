/*
@Time : 2019/4/17 16:55 
@Author : kenny zhu
@File : handlers.go
@Software: GoLand
@Others:
*/
package handlers

import (
	"github.com/gin-gonic/gin"
	"encoding/json"
	"io/ioutil"
	"fmt"

	proto "github.com/micro/go-micro/server/debug/proto"
)

var Handlers struct{
	Router *gin.Engine
}

// self init...
func init() {
	// gin.SetMode(gin.ReleaseMode)
	Handlers.Router = gin.Default()
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

// for health beat..
func NoModules(ctx *gin.Context) {
	fmt.Println("Received NoModules API request")

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


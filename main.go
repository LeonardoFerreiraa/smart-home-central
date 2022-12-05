package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	LUZ_DO_ESCRITORIO        = "204"
	LUZ_DO_QUARTO            = "209"
	LUZ_DO_DEPOSITO          = "210"
	VENTILADOR_DO_ESCRITORIO = "202"
	VENTILADOR_DO_QUARTO     = "205"
	VENTILADOR_DA_SALA       = "208"
	ABAJUR_DA_DIREITA        = "207"
	ABAJUR_DA_ESQUERDA       = "206"
	CAFETEIRA                = "201"
)

var actions = map[string]ActionHandler{
	"tag1": firstTagHandler,
	"tag2": secondTagHandler,
	"tag3": thirdTagHandler,
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/tag-trigger", tagRequestHandler)

	return r
}

func tagRequestHandler(ctx *gin.Context) {
	tagId := ctx.Query("tagId")
	action, prs := actions[tagId]
	if prs {
		action()
		ctx.Data(http.StatusOK, `text/html; charset=utf-8`, []byte(`<html><body><script>window.close();</script></body></html>`))
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": "tag not found"})
	}
}

func firstTagHandler() {
	doRequest(buildPowerRequest(LUZ_DO_QUARTO, "TOGGLE"))
	doRequest(buildPowerRequest(ABAJUR_DA_DIREITA, "TOGGLE"))
	doRequest(buildPowerRequest(ABAJUR_DA_ESQUERDA, "TOGGLE"))
	doRequest(buildPowerRequest(VENTILADOR_DO_QUARTO, "TOGGLE"))
}

func secondTagHandler() {
	doRequest(buildPowerRequest(LUZ_DO_DEPOSITO, "TOGGLE"))
}

func thirdTagHandler() {
	doRequest(buildPowerRequest(LUZ_DO_QUARTO, "OFF"))
	doRequest(buildPowerRequest(LUZ_DO_DEPOSITO, "OFF"))
	doRequest(buildPowerRequest(LUZ_DO_ESCRITORIO, "OFF"))
	doRequest(buildPowerRequest(VENTILADOR_DO_ESCRITORIO, "OFF"))
	doRequest(buildPowerRequest(VENTILADOR_DO_QUARTO, "OFF"))
	doRequest(buildPowerRequest(VENTILADOR_DA_SALA, "OFF"))
	doRequest(buildPowerRequest(ABAJUR_DA_DIREITA, "OFF"))
	doRequest(buildPowerRequest(ABAJUR_DA_ESQUERDA, "OFF"))
	doRequest(buildPowerRequest(CAFETEIRA, "OFF"))
}

func doRequest(req *http.Request) {
	client := &http.Client{}
	go client.Do(req)
}

func buildPowerRequest(deviceId string, status string) *http.Request {
	req, _ := http.NewRequest("GET", fmt.Sprintf("http://192.168.126.%s/cm", deviceId), nil)

	q := req.URL.Query()
	q.Add("cmnd", fmt.Sprintf("Power %s", status))

	req.URL.RawQuery = q.Encode()

	return req
}

type ActionHandler func()

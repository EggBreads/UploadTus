package main

import (
	"catenoid-company/tus-backend/controlHandles"
	"catenoid-company/tus-backend/lib"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

/**
	TODO Debug 시 Variable 중 문제가 발생될수 있는 표기가 발견
	TODO 그 부분을 수정하되, 전체적으로 다 확인
 */
func main() {
	//Log 설정

	//logger := &lumberjack.Logger{
	//	Filename: lib.LOGPATH,
	//	MaxSize: 100,
	//	MaxBackups: 3,
	//	MaxAge: 28,
	//}
	//
	//defer logger.Close()
	//
	//gin.DefaultWriter = logger
	//
	//log.SetOutput(logger)

	g := gin.Default()
	g.RedirectFixedPath = true

	tusInitConfig := &lib.TusConfigAndHandle{}

	config,err := tusInitConfig.SetStorage(lib.UPLOADPATH)

	if err != nil{
		fmt.Print(err)
		os.Exit(1)
	}

	h := &controlHandles.Handlers{UnroutedHandler: tusInitConfig.TusHandles, Engine: g}

	go h.SendResFromCompleted()


	gr := g.Group(lib.ROOTPATH)
	 {

		 g.Use(h.SetHeader(g))

		 gr.POST("", h.TusCreateFile())
		 gr.PATCH(":id", h.TusResumeStartFile())
		 gr.HEAD(":id", h.TusResumeSetFile())
		 gr.GET(":id", h.TusDownloadFile())

		 if config.StoreComposer.UsesTerminater{
			 gr.DELETE(":id", h.TusDeleteFile())
		 }
	 }


	//g.POST("", h.TusCreateFile())
	//g.PATCH(":id", h.TusResumeStartFile())
	//g.HEAD(":id", h.TusResumeSetFile())

	g.Run(":8081")
}

package main

import (
	"catenoid-company/tus-client/controller"
	"catenoid-company/tus-client/lib"
	"github.com/gin-gonic/gin"
	"log"
	"runtime"
)

/**
2020.11.06
Writer: Deuksoo Moon
Content: Tus을 이용한 파일 어이받기 서비스
 */
/**
	TODO Code Refactoring 해야됨
	TODO 전체 다 뜯어서 해야됨
 */
func main() {
	// Log 설정

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

	app := gin.Default()

	redisClient := lib.SetRedisConn()

	defer redisClient.Close()

	handlers := &controller.Handlers{
		Client: redisClient,
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	tgr := app.Group("/tus")
	{
		tgr.POST("/parallelFile", handlers.ParallelHandle)
		tgr.POST("/continuousFile", handlers.UploadContinuousHandle)
		tgr.POST("/moveToUploadFile", handlers.GetContinueUploadFile)
		tgr.POST("/deleteToUploadFile", handlers.DeleteHandle)
	}

	err := app.Run(":8082")


	if err != nil {
		log.Print(err)
	}

}

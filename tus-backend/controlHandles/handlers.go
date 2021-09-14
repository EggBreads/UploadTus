package controlHandles

import (
	"catenoid-company/tus-backend/lib"
	"fmt"
	"github.com/gin-gonic/gin"
	tusd "github.com/tus/tusd/pkg/handler"
	"log"
	"net/http"
)

type Handlers struct {
	*tusd.UnroutedHandler
	*gin.Engine
}

func (h *Handlers) SetHeaderTest() {
	//tusInitConfig := &lib.TusConfigAndHandle{}
	//StoragePath := g.
	//config,err := tusInitConfig.SetStorage(lib.UPLOADPATH)
	//
	//if err != nil{
	//	fmt.Print(err)
	//	os.Exit(1)
	//}

	//return gin.WrapH(http.StripPrefix(lib.ROOTPATH, h.Middleware(g)))
	//return gin.WrapH(http.StripPrefix("/", h.Middleware(g)))
}

func (h *Handlers) SetHeader(g *gin.Engine) gin.HandlerFunc {
	return gin.WrapH(http.StripPrefix(lib.ROOTPATH, h.Middleware(g)))
	//return gin.WrapH(http.StripPrefix("/", h.Middleware(g)))
}

func (h *Handlers) TusCreateFile() func(c *gin.Context) {
	return func(c *gin.Context) {
		u := c.Request.URL
		log.Print(u)
		h.PostFile(c.Writer, c.Request)
	}
}

func (h *Handlers) TusResumeSetFile() func(c *gin.Context) {
	return func(c *gin.Context) {
		h.HeadFile(c.Writer, c.Request)
	}
}

func (h *Handlers) TusResumeStartFile() func(c *gin.Context) {
	return func(c *gin.Context) {
		h.PatchFile(c.Writer, c.Request)
	}
}

func (h *Handlers) TusDownloadFile() func(c *gin.Context) {
	return func(c *gin.Context) {
		h.GetFile(c.Writer, c.Request)
	}
}

func (h *Handlers) TusDeleteFile() func(c *gin.Context) {
	return func(c *gin.Context) {
		h.DelFile(c.Writer, c.Request)
	}
}

func (h *Handlers) TusOptions(c *gin.Context) {
	c.JSON(http.StatusOK, "")
}

func (h *Handlers) SendResFromCompleted() {
	for {
		event := <-h.CompleteUploads
		fmt.Printf("Upload %s finished\n", event.Upload.ID)
	}
}

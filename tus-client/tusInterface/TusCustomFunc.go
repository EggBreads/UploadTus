package tusInterface

import (
	"github.com/eventials/go-tus"
	"net/http"
)

type CustomTusClient interface {
	NewUploadFromFile() (*tus.Upload, error)
}

type CustomTusInterface interface {
	//NewUploadFromFile() (*tus.Upload, error)
	//UploadFileCopy()
	SendToHttp(h map[string]string) (*http.Response, error)
	CreateAndCopyFromResFile(resp *http.Response) (map[string]interface{}, error)
}

type CustomTusModule interface {
	CreateUpload(u *tus.Upload) (*tus.Uploader, error)
	CreateOrResumeUpload(u *tus.Upload) (*tus.Uploader, error)
}

type CustomTusFileProcess interface {
	GetTusUpload() (*tus.Upload, error)
	TusFileCopy()
	DeleteContinuousFile() (*http.Response, error)
}

//type ExecuteTusUploads interface {
//	FileUpload(tusUtils uploads.TusUploads,tusStore *lib.TusStore)
//}
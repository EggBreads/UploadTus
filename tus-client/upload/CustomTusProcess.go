package uploads

import (
	"bytes"
	"catenoid-company/uploadTus/tus-client/lib"
	"catenoid-company/uploadTus/tus-client/tusInterface"
	"github.com/eventials/go-tus"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type CustomUtils struct {
	tFn tusInterface.CustomTusClient
	fn tusInterface.CustomTusInterface
}

type CustomTusUtils struct {
	C            *gin.Context
	SendHttpInfo map[string]string
}

func (ct *CustomTusUtils) NewUploadFromFile() (*tus.Upload, error) {
	_, header, errs := ct.C.Request.FormFile(lib.FILEFILEDNAME)

	if errs != nil {
		log.Print(errs)
	}

	metadata := map[string]string{
		"filename": header.Filename,
	}

	//fingerprint := fmt.Sprintf("%s-%d-%s", header.Filename, header.Size, time.Now())
	file, _ := header.Open()

	log.Print("[DEBUG] [NewUploadFromFile]")
	return tus.NewUpload(file, header.Size, metadata, ct.C.PostForm(lib.UPLOADQUERYKEYFILED)), nil
}

func (ct *CustomTusUtils) SendToHttp(h map[string]string) (*http.Response, error) {
	var req *http.Request
	var err error
	method := ct.SendHttpInfo[lib.METHOD]
	uri := ct.SendHttpInfo[lib.URI]
	params := ct.SendHttpInfo[lib.PARAMS]

	client := &http.Client{}

	if (params != "" && len(params) > 0) || method != "GET" {
		req, err = http.NewRequest(method, uri, bytes.NewReader([]byte(params)))
		if err != nil {
			return nil, err
		}

		if h != nil {
			lib.SetHeaders(req, h)
		}

		return client.Do(req)
	}

	req, err = http.NewRequest(method, uri, nil)
	if err != nil {
		return nil, err
	}

	if h != nil {
		lib.SetHeaders(req, h)
	}

	return client.Do(req)
}

func (ct *CustomTusUtils) CreateAndCopyFromResFile(resp *http.Response) (map[string]interface{}, error) {
	var out *os.File
	var filepath string
	var writeSize int64
	var err error
	var msg map[string]interface{}

	filepath = resp.Header.Get(lib.CONTENTDISPOSITION)
	pathArr := strings.Split(filepath, "=")
	filepath = strings.Replace(pathArr[len(pathArr)-1], "\"", "", -1)
	filepath = lib.STOREDIRPATH + filepath

	out, err = os.Create(filepath)

	if err != nil {
		log.Println(err)
		msg["status"] = "Fail"
		msg["msg"] = "There was a problem creating the file"
		return msg, err
	}
	defer out.Close()

	writeSize, err = io.Copy(out, resp.Body)
	if err != nil {
		msg["status"] = "Fail"
		msg["msg"] = "A problem occurred while writing the file"
		log.Println(err)
		return msg, err
	}

	log.Println(writeSize)

	return msg, nil
}

type ParallelUtils struct {
	C        *gin.Context
	Header   *multipart.FileHeader
}

func (p *ParallelUtils) NewUploadFromFile() (*tus.Upload, error) {
	//multi, header, errs := p.C.Request.FormFile(lib.FILEFILEDNAME)

	file, errs := p.Header.Open()

	if errs != nil {
		log.Print(errs)
	}

	byt := make([]byte, 1024)
	_, rErr := file.Read(byt)

	if rErr != nil {
		log.Println(rErr)
	}

	fileType := http.DetectContentType(byt)

	metadata := map[string]string{
		"filename": p.Header.Filename,
		"filetype": fileType,
	}

	//fingerprint := fmt.Sprintf("%s-%d-%s", header.Filename, header.Size, time.Now())


	log.Print("[DEBUG] [NewUploadFromFile]")
	return tus.NewUpload(file, p.Header.Size, metadata, p.C.PostForm(lib.UPLOADQUERYKEYFILED)), nil
}

//func (ct *CustomTusUtils) UploadFileCopy()  {
//	var resp *http.Response
//	var err error
//	var out *os.File
//	var filepath string
//	var writeSize int64
//
//	msg := map[string]string{
//		"status":"Success",
//		"msg" : "Move to upload file",
//	}
//
//	url := ct.c.PostForm(lib.FILEDOWNLOADNAME)
//
//	if url != ""{
//		msg["status"] = "Fail"
//		msg["msg"] = "The parameters are not correct"
//	}
//
//	resp, err = http.Get(url)
//	if err != nil {
//		log.Print(err)
//		msg["status"] = "Fail"
//		msg["msg"] = "Make sure the URL you requested is correct"
//	}
//
//	defer resp.Body.Close()
//
//	filepath = resp.Header.Get(lib.CONTENTDISPOSITION)
//	pathArr := strings.Split(filepath,"=")
//	filepath = strings.Replace(pathArr[len(pathArr)-1],"\"","",-1)
//	filepath  = lib.STOREDIRPATH + filepath
//
//	out, err = os.Create(filepath)
//
//	if err != nil {
//		log.Println(err)
//		msg["status"] = "Fail"
//		msg["msg"] = "There was a problem creating the file"
//	}
//	defer out.Close()
//
//	writeSize, err = io.Copy(out, resp.Body)
//	if err != nil {
//		msg["status"] = "Fail"
//		msg["msg"] = "A problem occurred while writing the file"
//		log.Println(err)
//	}
//
//
//	log.Println(writeSize)
//}

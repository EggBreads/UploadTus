package uploads

import (
	"catenoid-company/tus-client/lib"
	"catenoid-company/tus-client/tusInterface"
	"github.com/eventials/go-tus"
	"log"
	"sync"
)

type TusUploads struct {
	*TusUtils
	Fn tusInterface.CustomTusFileProcess
	ModuleFn tusInterface.CustomTusModule
}

func (t *TusUploads) RunTus(store tus.Store) *tus.Uploader {
	url := lib.HOST+lib.PATH
	tusConfig := &tus.Config{}

	uploadKey := t.Ctx.PostForm(lib.UPLOADQUERYKEYFILED)

	if  path, isPath := store.Get(uploadKey); path != "" && isPath{
		tusConfig.Store = store
		tusConfig.Resume = true
		tusConfig.OverridePatchMethod = false
		tusConfig.ChunkSize = lib.CHUNKSIZE
		url= path
	}else {
		tusConfig = nil
	}

	client, _ := tus.NewClient(url, tusConfig)
	// create an upload from a file.

	upload, err := t.Fn.GetTusUpload()
	if err != nil {
		log.Println(err)
	}

	// create the uploader.

	uploader, _ := client.CreateOrResumeUpload(upload)

	//uploadProcess := make(chan tus.Upload)
	//
	//uploader.NotifyUploadProgress(uploadProcess)

	if tusConfig == nil {
		store.Set(uploadKey, uploader.Url())
	}

	// start the uploading process.
	return uploader
}

type ParallelTusUploads struct {
	TFn        tusInterface.CustomTusClient
	ModuleFn   tusInterface.CustomTusModule
	Cc  	   *CustomClient
	PartialCnt int
	DownloadUrls []string
	FinishUrls []string
}

func (p *ParallelTusUploads) ParallelRun() bool {
	//defer mWc.Done()

	upload, e := p.TFn.NewUploadFromFile()

	if e != nil {
		log.Println(e)
	}

	var wg sync.WaitGroup

	for i := 0; i < p.PartialCnt; i++{
		wg.Add(1)
		go func() {
			defer wg.Done()
			uploader, _  := p.ModuleFn.CreateOrResumeUpload(upload)
			err := uploader.Upload()
			if err != nil {
				log.Println(err)
			}
			p.Cc.DownloadUrls = append(p.Cc.DownloadUrls, uploader.Url())
		}()
	}

	wg.Wait()

	p.Cc.IsFinish = true

	finishUploader,err := p.ModuleFn.CreateUpload(upload)

	if err != nil {
		log.Println("finish Error")
		return false
	}

	p.FinishUrls = append(p.FinishUrls, finishUploader.Url())
	log.Println(finishUploader.Url())

	return true
}

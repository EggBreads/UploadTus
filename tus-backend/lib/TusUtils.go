package lib

import (
	"github.com/tus/tusd/pkg/filestore"
	tusd "github.com/tus/tusd/pkg/handler"
	"golang.org/x/exp/errors"
)

type TusConfigAndHandle struct {
	TusHandles *tusd.UnroutedHandler
}

func (t *TusConfigAndHandle) SetStorage(UploadTmpPath string) (tusd.Config, error) {
	if UploadTmpPath == "" {
		return tusd.Config{},errors.New("Empty UploadPath")
	}

	store := &filestore.FileStore{
		Path: UploadTmpPath,
	}

	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	config := tusd.Config{
		BasePath:              ROOTPATH,
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
	}

	h, e := tusd.NewUnroutedHandler(config)
	if e != nil {
		return tusd.Config{},e
	}

	t.TusHandles = h
	return config,nil
}

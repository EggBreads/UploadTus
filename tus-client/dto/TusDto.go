package dto

type FileInfoDto struct {
	ID             string
	Size           int
	SizeIsDeferred bool
	Offset         int
	MetaData       MetaDataInFileInfoDto
	IsPartial      bool
	IsFinal        bool
	PartialUploads string
	Storage        StorageInFileInfoDto
}

type MetaDataInFileInfoDto struct {
	filename string
}

type StorageInFileInfoDto struct {
	Path string
	Type string
}

type ResponseDto struct {
	Status interface{}
	ResultMessage    interface{}
	ProcessStatus    interface{}
	SuccessUrls		 interface{}
	FailFileName	 []string
}

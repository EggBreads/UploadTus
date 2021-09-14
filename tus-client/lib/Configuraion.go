package lib

const (
	HOST = "http://localhost:8081/"
	PATH = "tus/"
	//PATH = "files/"
	//PATH             = ""
	FILEFILEDNAME    = "upload-file"
	FILEDOWNLOADNAME = "download-url"
	FILEDELETENAME   = "delete-name"
	UPLOADQUERYKEYFILED = "upload_key"
	METHOD           = "METHOD"
	URI              = "URI"
	PARAMS           = "PARAMS"
)

const (
	REDISSENTINELHOST     = "182.252.140.165"
	REDISSENTINELPORT     = "26379"
	REDISSENTINELMASTER   = "kollus"
	REDISSENTINELPOOLSIZE = "100"
	REDISTUSHASHKEY = "Tus"
)

const (
	CONTENTDISPOSITION = "Content-Disposition"

	TUSRESUMEABLE        = "Tus-Resumable"
	TUSRESUMEALBEVERSION = "1.0.0"

	TUSCONTENTLENGTH = "Content-Length"

	TUSXREQESTID    = "X-Request-ID"
	TUSUPLOADOFFSET = "Upload-Offset"
	TUSCONTENTTYPE  = "Content-Type"

	TUSUPLADMETADAT = "Upload-Metadata"

	CHUNKSIZE = 2 * 1024 * 1024
)

const (
	STOREDIRPATH = "/Users/kollus/http_upload/"
	LOGPATH      = "logs/tus-clinet.log"
)

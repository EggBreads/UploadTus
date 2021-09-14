package uploads

import (
	"github.com/eventials/go-tus"
	"io/ioutil"
	"net/http"
	netUrl "net/url"
	"strconv"
	"strings"
)

type CustomClient struct {
	Client *tus.Client
	PartialCnt int
	PartSuccessCnt chan int
	IsFinalCnt int
	DownloadUrls []string
	IsFinish bool
}

func (cc *CustomClient)CreateOrResumeUpload(u *tus.Upload) (*tus.Uploader, error)  {
	if u == nil {
		return nil, tus.ErrNilUpload
	}

	uploader, err := cc.Client.ResumeUpload(u)

	if err == nil {
		return uploader, err
	} else if (err == tus.ErrResumeNotEnabled) || (err == tus.ErrUploadNotFound) {
		return cc.CreateUpload(u)
	}

	return nil, err
}

func (cc *CustomClient)CreateUpload(u *tus.Upload) (*tus.Uploader, error)  {
	if u == nil {
		return nil, tus.ErrNilUpload
	}

	if cc.Client.Config.Resume && len(u.Fingerprint) == 0 {
		return nil, tus.ErrFingerprintNotSet
	}

	req, err := http.NewRequest("POST", cc.Client.Url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Length", "0")
	req.Header.Set("Upload-Length", strconv.FormatInt(u.Size(), 10))

	//req.Header.Set("Upload-Metadata", u.EncodedMetadata())
	// 병렬 처리하기하기위해 해당 파일 생성 header 처리
	if cc.IsFinish{
		concatVal := "final;"
		concatVal += strings.Join(cc.DownloadUrls," ")
		req.Header.Set("Upload-Concat",concatVal)
		req.Header.Set("Upload-Metadata", u.EncodedMetadata())
	}else{
		req.Header.Set("Upload-Concat","partial")
	}

	res, err := cc.Client.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 201:
		url := res.Header.Get("Location")

		baseUrl, err := netUrl.Parse(cc.Client.Url)
		if err != nil {
			return nil, tus.ErrUrlNotRecognized
		}

		newUrl, err := netUrl.Parse(url)
		if err != nil {
			return nil, tus.ErrUrlNotRecognized
		}
		if newUrl.Scheme == "" {
			newUrl.Scheme = baseUrl.Scheme
			url = newUrl.String()
		}

		if cc.Client.Config.Resume {
			cc.Client.Config.Store.Set(u.Fingerprint, url)
		}

		return tus.NewUploader(cc.Client, url, u, 0), nil
	case 412:
		return nil, tus.ErrVersionMismatch
	case 413:
		return nil, tus.ErrLargeUpload
	default:
		return nil, newClientError(res)
	}
}

func (cc *CustomClient) getUploadOffset(url string)  (int64, error) {
	req, err := http.NewRequest("HEAD", url, nil)

	if err != nil {
		return -1, err
	}

	if cc.Client.Config.Resume {
		req.Header.Set("Upload-Concat","partial")
	}

	res, err := cc.Client.Do(req)

	if err != nil {
		return -1, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		i, err := strconv.ParseInt(res.Header.Get("Upload-Offset"), 10, 64)

		if err == nil {
			return i, nil
		} else {
			return -1, err
		}
	case 403, 404, 410:
		// file doesn't exists.
		return -1, tus.ErrUploadNotFound
	case 412:
		return -1, tus.ErrVersionMismatch
	default:
		return -1, newClientError(res)
	}
}

func newClientError(res *http.Response) tus.ClientError {
	body, _ := ioutil.ReadAll(res.Body)
	return tus.ClientError{
		Code: res.StatusCode,
		Body: body,
	}
}



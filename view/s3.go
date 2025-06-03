package view

import "os"

// swagger:parameters UploadFile
type UploadFileReqWrap struct {
	// in:header
	Token string `json:"Authorization"`
	// in:formData
	// swagger:file
	File *os.File `json:"file"`
	// 上传到s3的key，不传值则为'dev/原文件名'
	// in:formData
	FileKey string `json:"file_key"`
	// 上传到s3的content- type，只支持（image/jpeg：jpg图片； image/png：png图片； audio/mpeg：视频；不传：binary/octet-stream，s3默认二进制格式；）
	// in:formData
	ContentType string `json:"content_type"`
}

// swagger:parameters UploadFileForOMC
type UploadFileReqWrapForOMC struct {
	UploadFileReqWrap
}

type UploadFileReq struct {
	// 文件key, 默认不填，服务端按规则生成，可以主动填写强制指定云端key
	FileKey       string
	LocalFilePath string
	ContentType   string
}

// swagger:parameters UploadFileContent
type UploadFileContentReqWrap struct {
	// in:header
	Token string `json:"Authorization"`
	// 文件内容base64
	// in:formData
	FileContent string `json:"f_base64_data"`
	// 上传到s3的key，不传值则为'dev/原文件名'
	// in:formData
	FileKey string `json:"file_key"`
	// 上传到s3的content- type，只支持（image/jpeg：jpg图片； image/png：png图片； audio/mpeg：视频；不传：binary/octet-stream，s3默认二进制格式；）
	// in:formData
	ContentType string `json:"content_type"`
}

// swagger:model
type UploadFileResp struct {
	FileUrl string `json:"file_url"`
	FileKey string `json:"file_key"`
}

// swagger:parameters DeleteFile
type DeleteFileReqWrap struct {
	// in:header
	Token string `json:"Authorization"`
	// s3 桶名
	// in:formData
	Bucket string `json:"bucket"`
	// 上传到s3的key
	// in:formData
	FileKey string `json:"file_key"`
}

// swagger:parameters DeleteFileForOMC
type DeleteFileReqWrapForOMC struct {
	DeleteFileReqWrap
}

// swagger:parameters GetFileForOMC
type GetFileReqWrapForOMC struct {
	GetFileReqWrap
}

// swagger:parameters GetFile
type GetFileReqWrap struct {
	// in:header
	Token string `json:"Authorization"`
	// s3 桶名
	// in:formData
	Bucket string `json:"bucket"`
	// 上传到s3的key
	// in:formData
	FileKey string `json:"file_key"`
}

// swagger:model
type GetFileResp struct {
	FileUrl string `json:"file_url"`
	FileKey string `json:"file_key"`
}

// swagger:parameters GetSignFile
type GetSignFileReqWrap struct {
	GetFileReqWrap
	// 过期时间，整数：表示多少秒后链接过期
	// in:formData
	Expires string `json:"expires"`
}

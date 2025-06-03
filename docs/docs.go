// Package classification go-zrbc-api
//
// Documentation of our go-zrbc-api API.
//
//	 Schemes: https
//	 BasePath: /
//	 Version: 1.0.0
//	 Title: QA API
//	 Host: wys.dev.zhanggao223.com
//
//	 Consumes:
//	 - application/json
//
//	 Produces:
//	 - application/json
//
//	 Security:
//	 - basic
//
//	SecurityDefinitions:
//	basic:
//	  type: basic
//
// swagger:meta
package docs

// Means the client request has error
// swagger:response BadRequestError
type BadRequestError struct {
	// in:string
	Error string `json:"err"`
	Msg   string `json:"msg"`
	Code  int    `json:"code"`
}

// Means the request is successful
// swagger:response ok
type OKResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Code   int    `json:"code"`
}

// Means the server has error
// swagger:response CommonError
type CommonErrorResponse struct {
	// in:body
	Body CommonError
}

type CommonError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Err  string `json:"err"`
}

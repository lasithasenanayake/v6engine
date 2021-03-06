package repositories

import (
	"duov6.com/objectstore/messaging"
	"encoding/json"
)

type RepositoryResponse struct {
	ResponseJson  string
	IsSuccess     bool
	IsImplemented bool
	Message       string
	Body          []byte
	Data          []map[string]interface{}
	Transaction   messaging.TransactionResponse
}

func (r *RepositoryResponse) GetErrorResponse(errorMessage string) {
	r.IsSuccess = false
	r.IsImplemented = true
	r.Message = errorMessage
}

func (r *RepositoryResponse) GetResponseWithBody(body []byte) {
	r.IsSuccess = true
	r.IsImplemented = true
	r.Message = "Operation Success!!!"
	r.Body = body
}

func (r *RepositoryResponse) GetSuccessResByObject(body interface{}) {
	r.IsSuccess = true
	r.IsImplemented = true
	r.Message = "Operation Success!!!"
	bytes, _ := json.Marshal(&body)
	r.Body = bytes //[:len(bytes)]
}

func (r *RepositoryResponse) GetSuccessResByString(body string) {
	r.IsSuccess = true
	r.IsImplemented = true
	r.Message = "Operation Success!!!"
	r.Body = ([]byte)(body)
}

package handler

import (
	"encoding/json"
	"fmt"
	"github.com/ironstone95/FlashQudoV2/handler/response"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type paging struct {
	limit int
	page  int
}

func newPaging(r *http.Request) *paging {
	p := new(paging)
	var limitQ, pageQ string
	var limit, page int
	var err error
	limitQ = r.URL.Query().Get("_limit")
	if len(limitQ) == 0 {
		limit = 10
	} else {
		limit, err = strconv.Atoi(limitQ)
		if err != nil {
			limit = 10
		}
	}
	err = nil
	pageQ = r.URL.Query().Get("_page")
	if len(pageQ) == 0 {
		page = 1
	} else {
		page, err = strconv.Atoi(pageQ)
		if err != nil {
			page = 1
		}
	}
	p.limit = limit
	p.page = page
	return p
}

func getParam(paramKey string, r *http.Request) (string, error) {
	vars := mux.Vars(r)
	if param, ok := vars[paramKey]; ok {
		return param, nil
	}
	return "", fmt.Errorf("parameter does not exist")
}

func SendError(rw http.ResponseWriter, message string, code int) {
	err := response.Error{Code: code, Message: message}
	enc := json.NewEncoder(rw)
	rw.WriteHeader(err.Code)
	_ = enc.Encode(err)
}

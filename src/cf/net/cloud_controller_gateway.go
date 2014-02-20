package net

import (
	"clocks"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"
)

func NewCloudControllerGateway(clock clocks.Clock) Gateway {
	invalidTokenCode := "1000"

	type ccErrorResponse struct {
		Code        int
		Description string
	}

	errorHandler := func(response *http.Response) errorResponse {
		headerBytes, _ := httputil.DumpResponse(response, false)

		jsonBytes, _ := ioutil.ReadAll(response.Body)
		response.Body.Close()

		ccResp := ccErrorResponse{}
		json.Unmarshal(jsonBytes, &ccResp)

		code := strconv.Itoa(ccResp.Code)
		if code == invalidTokenCode {
			code = INVALID_TOKEN_CODE
		}

		return errorResponse{
			Code:           code,
			Description:    ccResp.Description,
			ResponseBody:   string(jsonBytes),
			ResponseHeader: string(headerBytes),
		}
	}

	gateway := newGateway(clock, errorHandler)
	gateway.PollingEnabled = true
	return gateway
}

package sdk

import (
	"encoding/json"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"strings"
)

const (
	InvalidParamErrorCode                 = "InvalidParam"
	InvalidParameterErrorCode             = "InvalidParameter"
	UnauthorizedErrorCode                 = "Unauthorized"
	InvalidParamDateErrorMessage          = "The Param Date is invalid."
	InvalidParamAuthorizationErrorMessage = "The Param Authorization is invalid."
)

var errorCodeMap = map[string]string{
	"Forbidden.KeyNotFound":  "The specified Key is not found.",
	"Forbidden.NoPermission": "This operation is forbidden by permission system.",
	"InternalFailure":        "Internal Failure",
	"Rejected.Throttling":    "QPS Limit Exceeded",
}

type ErrorContent struct {
	HttpCode  int
	RequestId string
	HostId    string
	Code      string
	Message   string
}

func TransferTeaErrorServerError(err error) error {
	switch e := err.(type) {
	case *tea.SDKError:
		errData := &ErrorContent{}
		if e.Data != nil {
			err := json.Unmarshal([]byte(tea.StringValue(e.Data)), errData)
			if err != nil {
				return err
			}
		}
		errCode := tea.StringValue(e.Code)
		errMessage := tea.StringValue(e.Message)
		switch errCode {
		case InvalidParamErrorCode:
			if errMessage == InvalidParamDateErrorMessage {
				errData.Code = "IllegalTimestamp"
				errData.Message = `The input parameter "Timestamp" that is mandatory for processing this request is not supplied.`
			} else if errMessage == InvalidParamAuthorizationErrorMessage {
				errData.Code = "IncompleteSignature"
				errData.Message = "The request signature does not conform to Aliyun standards."
			}
		case UnauthorizedErrorCode:
			errData.Code = "InvalidAccessKeyId.NotFound"
			errData.Message = "The Access Key ID provided does not exist in our records."
		default:
			errData.Code = errCode
			errData.Message = errMessage
			msg, ok := errorCodeMap[errCode]
			if ok {
				errData.Message = msg
			}
		}
		responseContent, err := json.Marshal(errData)
		if err != nil {
			return err
		}
		return errors.NewServerError(errData.HttpCode, string(responseContent), "")
	}
	if strings.Contains(err.Error(), "Client.Timeout") {
		return errors.NewClientError(errors.TimeoutErrorCode, err.Error(), err)
	}
	return err
}

func TransferTeaErrorClientError(err error) error {
	switch e := err.(type) {
	case *tea.SDKError:
		return errors.NewClientError(tea.StringValue(e.Code), tea.StringValue(e.Message), err)
	}
	return err
}

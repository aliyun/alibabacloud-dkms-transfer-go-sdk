package sdk

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	dedicatedkmsopenapiutil "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi-util"
	dedicatedkmssdk "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/sdk"
	"io/ioutil"
	"net/http"
	"strings"
)

func (client *KmsTransferClient) GenerateDataKeyWithoutPlaintext(request *kms.GenerateDataKeyWithoutPlaintextRequest) (*kms.GenerateDataKeyWithoutPlaintextResponse, error) {
	if client.isUseKmsShareGateway {
		return client.Client.GenerateDataKeyWithoutPlaintext(request)
	}
	var aad []byte
	if request.EncryptionContext != "" {
		var err error
		aad, err = EncodeUserEncryptionContext(request.EncryptionContext)
		if err != nil {
			return nil, err
		}
	}
	numberOfBytesInteger := request.NumberOfBytes
	if !request.NumberOfBytes.HasValue() {
		if request.KeySpec == "" {
			numberOfBytesInteger = "32"
		} else if request.KeySpec == KmsKeySpecAes256 {
			numberOfBytesInteger = "32"
		} else if request.KeySpec == KmsKeySpecAes128 {
			numberOfBytesInteger = "16"
		} else {
			return nil, TransferTeaErrorServerError(tea.NewSDKError(map[string]interface{}{
				"code":    InvalidParameterErrorCode,
				"message": "The specified parameter KeySpec is not valid.",
				"data":    nil,
			}))
		}
	}
	numberOfBytes, err := numberOfBytesInteger.GetValue()
	if err != nil {
		return nil, err
	}
	dkmsRequest := &dedicatedkmssdk.GenerateDataKeyRequest{
		KeyId:         tea.String(request.KeyId),
		NumberOfBytes: tea.Int32(int32(numberOfBytes)),
		Aad:           aad,
	}
	ignoreSSL := client.GetHTTPSInsecure()
	runtimeOptions := &dedicatedkmsopenapiutil.RuntimeOptions{
		Verify:    tea.String(client.Verify),
		IgnoreSSL: tea.Bool(ignoreSSL),
	}
	runtimeOptions.Headers = append(runtimeOptions.Headers, tea.String(MigrationKeyVersionIdKey))

	dkmsResponse, err := client.dkmsClient.GenerateDataKeyWithOptions(dkmsRequest, runtimeOptions)
	if err != nil {
		return nil, TransferTeaErrorServerError(err)
	}
	dkmsEncryptRequest := &dedicatedkmssdk.EncryptRequest{
		KeyId:     tea.String(request.KeyId),
		Plaintext: []byte(base64.StdEncoding.EncodeToString(dkmsResponse.Plaintext)),
		Aad:       aad,
	}
	encryptResponse, err := client.dkmsClient.EncryptWithOptions(dkmsEncryptRequest, runtimeOptions)
	if err != nil {
		return nil, TransferTeaErrorServerError(err)
	}
	dkmsResponse.Iv = encryptResponse.Iv
	dkmsResponse.CiphertextBlob = encryptResponse.CiphertextBlob

	keyVersionId, ok := dkmsResponse.Headers[MigrationKeyVersionIdKey]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Can not found response headers parameter[%s]", MigrationKeyVersionIdKey))
	}
	mkvId := []byte(tea.StringValue(keyVersionId))

	var ciphertextBlob []byte
	ciphertextBlob = append(ciphertextBlob, mkvId...)
	ciphertextBlob = append(ciphertextBlob, dkmsResponse.Iv...)
	ciphertextBlob = append(ciphertextBlob, dkmsResponse.CiphertextBlob...)

	kmsResponse := kms.CreateGenerateDataKeyWithoutPlaintextResponse()
	kmsResponse.KeyId = tea.StringValue(dkmsResponse.KeyId)
	kmsResponse.KeyVersionId = tea.StringValue(keyVersionId)
	kmsResponse.CiphertextBlob = base64.StdEncoding.EncodeToString(ciphertextBlob)
	kmsResponse.RequestId = tea.StringValue(dkmsResponse.RequestId)
	var body []byte
	if strings.ToUpper(request.AcceptFormat) == "JSON" {
		body, err = json.Marshal(kmsResponse)
		if err != nil {
			return nil, err
		}
	} else if strings.ToUpper(request.AcceptFormat) == "XML" {
		body, err = xml.Marshal(kmsResponse)
		if err != nil {
			return nil, err
		}
	}
	httpResponse := &http.Response{}
	httpResponse.StatusCode = http.StatusOK
	httpResponse.Body = ioutil.NopCloser(bytes.NewReader(body))

	err = responses.Unmarshal(kmsResponse, httpResponse, request.AcceptFormat)
	if err != nil {
		return nil, err
	}
	return kmsResponse, nil
}

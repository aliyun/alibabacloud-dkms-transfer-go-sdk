package sdk

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	dedicatedkmsopenapiutil "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi-util"
	dedicatedkmssdk "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/sdk"
	"io/ioutil"
	"net/http"
)

func (client *KmsTransferClient) Encrypt(request *kms.EncryptRequest) (*kms.EncryptResponse, error) {
	dkmsRequest := &dedicatedkmssdk.EncryptRequest{
		KeyId:     tea.String(request.KeyId),
		Plaintext: []byte(request.Plaintext),
		Aad:       []byte(request.EncryptionContext),
	}
	ignoreSSL := client.GetHTTPSInsecure()
	runtimeOptions := &dedicatedkmsopenapiutil.RuntimeOptions{
		Verify:    tea.String(client.Verify),
		IgnoreSSL: tea.Bool(ignoreSSL),
	}
	runtimeOptions.Headers = append(runtimeOptions.Headers, tea.String(MigrationKeyVersionIdKey))

	dkmsResponse, err := client.dkmsClient.EncryptWithOptions(dkmsRequest, runtimeOptions)
	if err != nil {
		return nil, TransferTeaErrorServerError(err)
	}

	versionId, ok := dkmsResponse.Headers[MigrationKeyVersionIdKey]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Can not found response headers parameter[%s]", MigrationKeyVersionIdKey))
	}
	mkvId := []byte(tea.StringValue(versionId))

	var ciphertextBlob []byte
	ciphertextBlob = append(ciphertextBlob, mkvId...)
	ciphertextBlob = append(ciphertextBlob, dkmsResponse.Iv...)
	ciphertextBlob = append(ciphertextBlob, dkmsResponse.CiphertextBlob...)

	kmsResponse := kms.CreateEncryptResponse()
	kmsResponse.KeyId = tea.StringValue(dkmsResponse.KeyId)
	kmsResponse.CiphertextBlob = base64.StdEncoding.EncodeToString(ciphertextBlob)
	kmsResponse.RequestId = tea.StringValue(dkmsResponse.RequestId)
	body, err := json.Marshal(kmsResponse)
	if err != nil {
		return nil, err
	}
	httpResponse := &http.Response{}
	httpResponse.StatusCode = http.StatusOK
	httpResponse.Body = ioutil.NopCloser(bytes.NewReader(body))

	err = responses.Unmarshal(kmsResponse, httpResponse, "JSON")
	if err != nil {
		return nil, err
	}
	return kmsResponse, nil
}

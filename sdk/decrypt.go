package sdk

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	dedicatedkmsopenapiutil "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi-util"
	dedicatedkmssdk "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/sdk"
	"io/ioutil"
	"net/http"
)

const (
	EktIdLength              = 36
	GcmIvLength              = 12
	MigrationKeyVersionIdKey = "x-kms-migrationkeyversionid"
)

func (client *KmsTransferClient) Decrypt(request *kms.DecryptRequest) (*kms.DecryptResponse, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(request.CiphertextBlob)
	if err != nil {
		return nil, err
	}
	ektId := ciphertext[:EktIdLength]
	iv := ciphertext[EktIdLength : EktIdLength+GcmIvLength]
	ciphertextBlob := ciphertext[EktIdLength+GcmIvLength:]
	dkmsRequest := &dedicatedkmssdk.DecryptRequest{
		Headers:        make(map[string]*string),
		CiphertextBlob: ciphertextBlob,
		Iv:             iv,
		Aad:            []byte(request.EncryptionContext),
	}
	dkmsRequest.Headers[MigrationKeyVersionIdKey] = tea.String(string(ektId))
	ignoreSSL := client.GetHTTPSInsecure()
	runtimeOptions := &dedicatedkmsopenapiutil.RuntimeOptions{
		Verify:    tea.String(client.Verify),
		IgnoreSSL: tea.Bool(ignoreSSL),
	}

	dkmsResponse, err := client.dkmsClient.DecryptWithOptions(dkmsRequest, runtimeOptions)
	if err != nil {
		return nil, TransferTeaErrorServerError(err)
	}

	kmsResponse := kms.CreateDecryptResponse()
	kmsResponse.KeyId = tea.StringValue(dkmsResponse.KeyId)
	kmsResponse.Plaintext = string(dkmsResponse.Plaintext)
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

package sdk

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	dedicatedkmsopenapiutil "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi-util"
	dedicatedkmssdk "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/sdk"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	EktIdLength              = 36
	GcmIvLength              = 12
	MigrationKeyVersionIdKey = "x-kms-migrationkeyversionid"
)

func (client *KmsTransferClient) Decrypt(request *kms.DecryptRequest) (*kms.DecryptResponse, error) {
	if client.isUseKmsShareGateway {
		return client.Client.Decrypt(request)
	}
	var aad []byte
	if request.EncryptionContext != "" {
		var err error
		aad, err = EncodeUserEncryptionContext(request.EncryptionContext)
		if err != nil {
			return nil, err
		}
	}
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
		Aad:            aad,
	}
	dkmsRequest.Headers[MigrationKeyVersionIdKey] = tea.String(string(ektId))
	ignoreSSL := client.GetHTTPSInsecure()
	runtimeOptions := &dedicatedkmsopenapiutil.RuntimeOptions{
		Verify:    tea.String(client.Verify),
		IgnoreSSL: tea.Bool(ignoreSSL),
	}
	runtimeOptions.Headers = append(runtimeOptions.Headers, tea.String(MigrationKeyVersionIdKey))

	dkmsResponse, err := client.dkmsClient.DecryptWithOptions(dkmsRequest, runtimeOptions)
	if err != nil {
		return nil, TransferTeaErrorServerError(err)
	}
	keyVersionId, _ := dkmsResponse.Headers[MigrationKeyVersionIdKey]
	kmsResponse := kms.CreateDecryptResponse()
	kmsResponse.KeyId = tea.StringValue(dkmsResponse.KeyId)
	kmsResponse.KeyVersionId = tea.StringValue(keyVersionId)
	kmsResponse.Plaintext = string(dkmsResponse.Plaintext)
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

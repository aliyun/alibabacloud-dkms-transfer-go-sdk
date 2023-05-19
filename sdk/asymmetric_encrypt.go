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

func (client *KmsTransferClient) AsymmetricEncrypt(request *kms.AsymmetricEncryptRequest) (*kms.AsymmetricEncryptResponse, error) {
	if client.isUseKmsShareGateway {
		return client.Client.AsymmetricEncrypt(request)
	}
	plaintext, err := base64.StdEncoding.DecodeString(request.Plaintext)
	if err != nil {
		return nil, err
	}
	dkmsRequest := &dedicatedkmssdk.EncryptRequest{
		Headers:   make(map[string]*string),
		KeyId:     tea.String(request.KeyId),
		Plaintext: plaintext,
		Algorithm: tea.String(request.Algorithm),
	}
	if request.KeyVersionId != "" {
		dkmsRequest.Headers[MigrationKeyVersionIdKey] = tea.String(request.KeyVersionId)
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

	keyVersionId, _ := dkmsResponse.Headers[MigrationKeyVersionIdKey]
	kmsResponse := kms.CreateAsymmetricEncryptResponse()
	kmsResponse.KeyId = tea.StringValue(dkmsResponse.KeyId)
	kmsResponse.KeyVersionId = tea.StringValue(keyVersionId)
	kmsResponse.CiphertextBlob = base64.StdEncoding.EncodeToString(dkmsResponse.CiphertextBlob)
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

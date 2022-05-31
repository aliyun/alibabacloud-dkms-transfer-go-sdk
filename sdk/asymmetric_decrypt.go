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

func (client *KmsTransferClient) AsymmetricDecrypt(request *kms.AsymmetricDecryptRequest) (*kms.AsymmetricDecryptResponse, error) {
	ciphertextBlob, err := base64.StdEncoding.DecodeString(request.CiphertextBlob)
	if err != nil {
		return nil, err
	}
	dkmsRequest := &dedicatedkmssdk.DecryptRequest{
		KeyId:          tea.String(request.KeyId),
		CiphertextBlob: ciphertextBlob,
		Algorithm:      tea.String(request.Algorithm),
	}
	ignoreSSL := client.GetHTTPSInsecure()
	runtimeOptions := &dedicatedkmsopenapiutil.RuntimeOptions{
		Verify:    tea.String(client.Verify),
		IgnoreSSL: tea.Bool(ignoreSSL),
	}

	dkmsResponse, err := client.dkmsClient.DecryptWithOptions(dkmsRequest, runtimeOptions)
	if err != nil {
		return nil, TransferTeaErrorServerError(err)
	}

	kmsResponse := kms.CreateAsymmetricDecryptResponse()
	kmsResponse.KeyId = tea.StringValue(dkmsResponse.KeyId)
	kmsResponse.Plaintext = base64.StdEncoding.EncodeToString(dkmsResponse.Plaintext)
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

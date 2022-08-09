package sdk

import (
	"bytes"
	"encoding/json"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	dedicatedkmsopenapiutil "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi-util"
	dedicatedkmssdk "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/sdk"
	"io/ioutil"
	"net/http"
)

func (client *KmsTransferClient) GetSecretValue(request *kms.GetSecretValueRequest) (*kms.GetSecretValueResponse, error) {
	fetchExtendedConfig, _ := request.FetchExtendedConfig.GetValue()
	dkmsRequest := &dedicatedkmssdk.GetSecretValueRequest{
		SecretName:          tea.String(request.SecretName),
		VersionStage:        tea.String(request.VersionStage),
		VersionId:           tea.String(request.VersionId),
		FetchExtendedConfig: tea.Bool(fetchExtendedConfig),
	}
	ignoreSSL := client.GetHTTPSInsecure()
	runtimeOptions := &dedicatedkmsopenapiutil.RuntimeOptions{
		Verify:    tea.String(client.Verify),
		IgnoreSSL: tea.Bool(ignoreSSL),
	}

	dkmsResponse, err := client.dkmsClient.GetSecretValueWithOptions(dkmsRequest, runtimeOptions)
	if err != nil {
		return nil, TransferTeaErrorServerError(err)
	}

	kmsResponse := kms.CreateGetSecretValueResponse()
	kmsResponse.SecretName = tea.StringValue(dkmsResponse.SecretName)
	kmsResponse.VersionId = tea.StringValue(dkmsResponse.VersionId)
	kmsResponse.CreateTime = tea.StringValue(dkmsResponse.CreateTime)
	kmsResponse.SecretData = tea.StringValue(dkmsResponse.SecretData)
	kmsResponse.SecretDataType = tea.StringValue(dkmsResponse.SecretDataType)
	kmsResponse.AutomaticRotation = tea.StringValue(dkmsResponse.AutomaticRotation)
	kmsResponse.RotationInterval = tea.StringValue(dkmsResponse.RotationInterval)
	kmsResponse.NextRotationDate = tea.StringValue(dkmsResponse.NextRotationDate)
	kmsResponse.ExtendedConfig = tea.StringValue(dkmsResponse.ExtendedConfig)
	kmsResponse.LastRotationDate = tea.StringValue(dkmsResponse.LastRotationDate)
	kmsResponse.SecretType = tea.StringValue(dkmsResponse.SecretType)
	for _, state := range dkmsResponse.VersionStages {
		kmsResponse.VersionStages.VersionStage = append(kmsResponse.VersionStages.VersionStage, tea.StringValue(state))
	}
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

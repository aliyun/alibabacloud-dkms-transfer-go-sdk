package sdk

import (
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	credentialsprovider "github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials/provider"
	dedicatedkmsopenapi "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	regionId        = "regionId"
	accessKeyId     = "accessKeyId"
	accessKeySecret = "accessKeySecret"

	clientKeyFile = "path/to/ClientKey.json"
	password      = "1234"
	endpoint      = "endpoint"
)

var config = &dedicatedkmsopenapi.Config{
	Protocol:      tea.String("https"),
	ClientKeyFile: tea.String(clientKeyFile),
	Password:      tea.String(password),
	Endpoint:      tea.String(endpoint),
}

func TestNewClientWithAccessKey(t *testing.T) {
	client, err := NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret, config)

	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestNewClientWithProvider(t *testing.T) {
	provider := credentialsprovider.NewInstanceCredentialsProvider()
	client, err := NewClientWithProvider(regionId, config, provider)

	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestNewClientWithOptions(t *testing.T) {
	kmsConfig := sdk.NewConfig()
	credential := credentials.NewAccessKeyCredential(accessKeyId, accessKeySecret)
	client, err := NewClientWithOptions(regionId, kmsConfig, credential, config)

	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestNewClientWithStsToken(t *testing.T) {
	stsAccessKeyId := "stsAccessKeyId"
	stsAccessKeySecret := "stsAccessKeySecret"
	stsToken := "stsToken"
	client, err := NewClientWithStsToken(regionId, stsAccessKeyId, stsAccessKeySecret, stsToken, config)

	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestNewClientWithRamRoleArn(t *testing.T) {
	roleArn := "roleArn"
	roleSessionName := "roleSessionName"
	client, err := NewClientWithRamRoleArn(regionId, accessKeyId, accessKeySecret, roleArn, roleSessionName, config)

	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestNewClientWithEcsRamRole(t *testing.T) {
	roleName := "roleName"
	client, err := NewClientWithEcsRamRole(regionId, roleName, config)

	assert.Nil(t, err)
	assert.NotNil(t, client)
}

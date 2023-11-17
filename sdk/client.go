package sdk

import (
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials/provider"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	dedicatedkmsopenapi "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi"
	dedicatedkmssdk "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/sdk"
)

const (
	SDK_NSAME      = "alibabacloud-dkms-transfer-go-sdk"
	SDK_VERSION    = "0.1.8"
	SDK_USER_AGENT = SDK_NSAME + "/" + SDK_VERSION
)

type KmsTransferClient struct {
	*kms.Client
	dkmsClient           *dedicatedkmssdk.Client
	isUseKmsShareGateway bool
	Verify               string
}

func (client *KmsTransferClient) SetVerify(v string) {
	client.Verify = v
}

func NewClientWithProvider(regionId string, dkmsConfig *dedicatedkmsopenapi.Config, providers ...provider.Provider) (*KmsTransferClient, error) {
	kmsClient, err := kms.NewClientWithProvider(regionId, providers...)
	if err != nil {
		return nil, err
	}
	if dkmsConfig != nil {
		setUserAgent(dkmsConfig)
		dkmsClient, err := dedicatedkmssdk.NewClient(dkmsConfig)
		if err != nil {
			return nil, TransferTeaErrorClientError(err)
		}
		return &KmsTransferClient{Client: kmsClient, dkmsClient: dkmsClient}, nil
	} else {
		return &KmsTransferClient{Client: kmsClient, isUseKmsShareGateway: true}, nil
	}

}

func NewClientWithOptions(regionId string, kmsConfig *sdk.Config, credential auth.Credential, dkmsConfig *dedicatedkmsopenapi.Config) (*KmsTransferClient, error) {
	kmsClient, err := kms.NewClientWithOptions(regionId, kmsConfig, credential)
	if err != nil {
		return nil, err
	}
	if dkmsConfig != nil {
		setUserAgent(dkmsConfig)
		dkmsClient, err := dedicatedkmssdk.NewClient(dkmsConfig)
		if err != nil {
			return nil, TransferTeaErrorClientError(err)
		}
		return &KmsTransferClient{Client: kmsClient, dkmsClient: dkmsClient}, nil
	} else {
		return &KmsTransferClient{Client: kmsClient, isUseKmsShareGateway: true}, nil
	}
}

func NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret string, dkmsConfig *dedicatedkmsopenapi.Config) (*KmsTransferClient, error) {
	kmsClient, err := kms.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret)
	if err != nil {
		return nil, err
	}
	if dkmsConfig != nil {
		setUserAgent(dkmsConfig)
		dkmsClient, err := dedicatedkmssdk.NewClient(dkmsConfig)
		if err != nil {
			return nil, TransferTeaErrorClientError(err)
		}
		return &KmsTransferClient{Client: kmsClient, dkmsClient: dkmsClient}, nil
	} else {
		return &KmsTransferClient{Client: kmsClient, isUseKmsShareGateway: true}, nil
	}
}

func NewClientWithStsToken(regionId, stsAccessKeyId, stsAccessKeySecret, stsToken string, dkmsConfig *dedicatedkmsopenapi.Config) (*KmsTransferClient, error) {
	kmsClient, err := kms.NewClientWithStsToken(regionId, stsAccessKeyId, stsAccessKeySecret, stsToken)
	if err != nil {
		return nil, err
	}
	if dkmsConfig != nil {
		setUserAgent(dkmsConfig)
		dkmsClient, err := dedicatedkmssdk.NewClient(dkmsConfig)
		if err != nil {
			return nil, TransferTeaErrorClientError(err)
		}
		return &KmsTransferClient{Client: kmsClient, dkmsClient: dkmsClient}, nil
	} else {
		return &KmsTransferClient{Client: kmsClient, isUseKmsShareGateway: true}, nil
	}
}

func NewClientWithRamRoleArn(regionId string, accessKeyId, accessKeySecret, roleArn, roleSessionName string, dkmsConfig *dedicatedkmsopenapi.Config) (*KmsTransferClient, error) {
	kmsClient, err := kms.NewClientWithRamRoleArn(regionId, accessKeyId, accessKeySecret, roleArn, roleSessionName)
	if err != nil {
		return nil, err
	}
	if dkmsConfig != nil {
		setUserAgent(dkmsConfig)
		dkmsClient, err := dedicatedkmssdk.NewClient(dkmsConfig)
		if err != nil {
			return nil, TransferTeaErrorClientError(err)
		}
		return &KmsTransferClient{Client: kmsClient, dkmsClient: dkmsClient}, nil
	} else {
		return &KmsTransferClient{Client: kmsClient, isUseKmsShareGateway: true}, nil
	}
}

func NewClientWithRamRoleArnAndPolicy(regionId string, accessKeyId, accessKeySecret, roleArn, roleSessionName, policy string, dkmsConfig *dedicatedkmsopenapi.Config) (*KmsTransferClient, error) {
	kmsClient, err := kms.NewClientWithRamRoleArnAndPolicy(regionId, accessKeyId, accessKeySecret, roleArn, roleSessionName, policy)
	if err != nil {
		return nil, err
	}
	if dkmsConfig != nil {
		setUserAgent(dkmsConfig)
		dkmsClient, err := dedicatedkmssdk.NewClient(dkmsConfig)
		if err != nil {
			return nil, TransferTeaErrorClientError(err)
		}
		return &KmsTransferClient{Client: kmsClient, dkmsClient: dkmsClient}, nil
	} else {
		return &KmsTransferClient{Client: kmsClient, isUseKmsShareGateway: true}, nil
	}
}

func NewClientWithEcsRamRole(regionId string, roleName string, dkmsConfig *dedicatedkmsopenapi.Config) (*KmsTransferClient, error) {
	kmsClient, err := kms.NewClientWithEcsRamRole(regionId, roleName)
	if err != nil {
		return nil, err
	}
	if dkmsConfig != nil {
		setUserAgent(dkmsConfig)
		dkmsClient, err := dedicatedkmssdk.NewClient(dkmsConfig)
		if err != nil {
			return nil, TransferTeaErrorClientError(err)
		}
		return &KmsTransferClient{Client: kmsClient, dkmsClient: dkmsClient}, nil
	} else {
		return &KmsTransferClient{Client: kmsClient, isUseKmsShareGateway: true}, nil
	}
}

func NewClientWithRsaKeyPair(regionId string, publicKeyId, privateKey string, sessionExpiration int, dkmsConfig *dedicatedkmsopenapi.Config) (*KmsTransferClient, error) {
	kmsClient, err := kms.NewClientWithRsaKeyPair(regionId, publicKeyId, privateKey, sessionExpiration)
	if err != nil {
		return nil, err
	}
	if dkmsConfig != nil {
		setUserAgent(dkmsConfig)
		dkmsClient, err := dedicatedkmssdk.NewClient(dkmsConfig)
		if err != nil {
			return nil, TransferTeaErrorClientError(err)
		}
		return &KmsTransferClient{Client: kmsClient, dkmsClient: dkmsClient}, nil
	} else {
		return &KmsTransferClient{Client: kmsClient, isUseKmsShareGateway: true}, nil
	}
}

func setUserAgent(config *dedicatedkmsopenapi.Config) {
	if config.UserAgent != nil {
		config.UserAgent = tea.String(tea.StringValue(config.UserAgent) + " " + SDK_USER_AGENT)
	} else {
		config.UserAgent = tea.String(SDK_USER_AGENT)
	}

}

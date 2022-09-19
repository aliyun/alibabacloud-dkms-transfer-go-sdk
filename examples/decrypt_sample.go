package main

import (
	"fmt"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	dedicatedkmsopenapi "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi"
	"github.com/aliyun/alibabacloud-dkms-transfer-go-sdk/sdk"
	"io/ioutil"
)

func main() {
	regionId := "<your region id>"
	accessKeyId := "<your access key id>"
	accessKeySecret := "<your access key secret>"

	clientKeyContent := `<your client key content>`
	password := "<your client key password>"
	endpoint := "<your dkms instance service endpoint>"

	ciphertextBlob := "<your ciphertextBlob>"

	config := &dedicatedkmsopenapi.Config{
		Protocol:         tea.String("https"),
		ClientKeyContent: tea.String(clientKeyContent),
		Password:         tea.String(password),
		Endpoint:         tea.String(endpoint),
	}
	client, err := sdk.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret, config)
	if err != nil {
		panic(err)
	}
	// 验证服务端证书
	ca, err := ioutil.ReadFile("path/to/caCert.pem")
	if err != nil {
		panic(err)
	}
	client.SetVerify(string(ca))
	// 如需忽略服务端证书验证,可打开此处注释代码
	//client.SetHTTPSInsecure(true)

	request := kms.CreateDecryptRequest()
	request.CiphertextBlob = ciphertextBlob

	result, err := client.Decrypt(request)
	if err != nil {
		panic(err)
	}

	fmt.Println("KeyId:", result.KeyId)
	fmt.Println("KeyVersionId:", result.KeyVersionId)
	fmt.Println("Plaintext:", result.Plaintext)
	fmt.Println("RequestId:", result.RequestId)

}

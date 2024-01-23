English | [简体中文](README-CN.md)

# Alibaba Cloud Dedicated KMS Transfer SDK for Go

![](https://aliyunsdk-pages.alicdn.com/icons/AlibabaCloud.svg)

Alibaba Cloud Dedicated KMS Transfer SDK for Go can help Golang developers to migrate from the KMS keys to the Dedicated KMS keys.

- [Alibaba Cloud Dedicated KMS Homepage](https://www.alibabacloud.com/help/zh/doc-detail/311016.htm)
- [Sample Code](/examples)
- [Issues](https://github.com/aliyun/alibabacloud-dkms-transfer-go-sdk/issues)
- [Release](https://github.com/aliyun/alibabacloud-dkms-transfer-go-sdk/releases)

## Features
* Dedicated KMS provides a tenant-specific instance that is deployed in the VPC of a tenant to allow access over an internal network.
* Dedicated KMS uses a tenant-specific cryptographic resource pool to implement resource isolation and cryptographic isolation. This improves security.
* Dedicated KMS simplifies the management of HSMs. You can use the stable, easy-to-use upper-layer key management features and cryptographic operations provided by Dedicated KMS to manage your HSMs.
* Dedicated KMS allows you to integrate your HSMs with Alibaba Cloud services in a seamless manner. This delivers secure and controllable encryption capabilities for Alibaba Cloud services. For more information, see [Alibaba Cloud services that can be integrated with KMS](https://www.alibabacloud.com/help/en/key-management-service/latest/alibaba-cloud-services-that-can-be-integrated-with-kms#concept-2318937).
* Reduce the cost of migrating the Shared KMS keys to Dedicated KMS keys. 

## Requirements

- Golang 1.12 or later.

## Installation

If you use `go mod` to manage your dependence, You can declare the dependency on AlibabaCloud DKMS SDK for Go in the
go.mod file:

```text
require (
	github.com/aliyun/alibabacloud-dkms-transfer-go-sdk v0.1.9
)
```

Or, Run the following command to get the remote code package:

```shell
$ go get -u github.com/aliyun/alibabacloud-dkms-transfer-go-sdk
```

## Client Mechanism
Alibaba Cloud Dedicated KMS Transfer SDK for Go transfers the the following method of request to dedicated KMS vpc gateway by default.

* Encrypt
* Decrypt
* GenerateDataKey
* GenerateDataKeyWithoutPlaintext
* GetPublicKey
* AsymmetricEncrypt
* AsymmetricDecrypt
* AsymmetricSign
* AsymmetricVerify
* GetSecretValue

## Quick Examples

```go
package example

import (
	"fmt"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	dedicatedkmsopenapi "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi"
	"github.com/aliyun/alibabacloud-dkms-transfer-go-sdk/sdk"
	"io/ioutil"
)

func main() {
	config := &dedicatedkmsopenapi.Config{
		Protocol:         tea.String("https"),
		ClientKeyContent: tea.String("<your client key content>"),
		Password:         tea.String("<your client key password>"),
		Endpoint:         tea.String("<your dkms instance service endpoint>"),
	}
	client, err := sdk.NewClientWithAccessKey("<your region id>", "<your access key id>", "<your access key secret>", config)
	if err != nil {
		panic(err)
	}
	
	// verify CA cert
	ca, err := ioutil.ReadFile("path/to/caCert.pem")
	if err != nil {
		panic(err)
	}
	client.SetVerify(string(ca))
	// or, ignore CA cert
	//client.SetHTTPSInsecure(true)

	request := kms.CreateEncryptRequest()
	request.KeyId = "<your key id>"
	request.Plaintext = "<your plaintext>"

	result, err := client.Encrypt(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

```

## License

[Apache-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Copyright (c) 2009-present, Alibaba Cloud All rights reserved.
 

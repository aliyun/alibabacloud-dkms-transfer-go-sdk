[English](README.md) | 简体中文

# 阿里云专属KMS适配SDK for Go

![](https://aliyunsdk-pages.alicdn.com/icons/AlibabaCloud.svg)

阿里云专属KMS适配SDK for Go可以帮助Golang开发者快速完成由KMS密钥向专属KMS密钥迁移适配工作。

- [阿里云专属KMS主页](https://help.aliyun.com/document_detail/311016.html)
- [代码示例](/examples)
- [Issues](https://github.com/aliyun/alibabacloud-dkms-transfer-go-sdk/issues)
- [Release](https://github.com/aliyun/alibabacloud-dkms-transfer-go-sdk/releases)

## 优势

* 专属KMS提供租户独享的服务实例，并部署到租户的VPC内，满足私有网络接入需求。
* 专属KMS使用租户独享的密码资源池（HSM集群），实现资源隔离和密码学隔离，以获得更高的安全性。
* 专属KMS可以降低使用HSM的复杂度，为您的HSM提供稳定、易用的上层密钥管理途径和密码计算服务。
*
专属KMS可以将您的HSM与云服务无缝集成，为云服务加密提供更高的安全性和可控制性。更多信息，请参见[支持服务端集成加密的云服务](https://help.aliyun.com/document_detail/141499.htm?#concept-2318937)
。
* 减低用户从共享KMS密钥移专属KMS密钥的成本

## 软件要求

- Golang 1.12及以上。

## 安装

您可以使用`go mod`管理您的依赖：

```
require (
	github.com/aliyun/alibabacloud-dkms-transfer-go-sdk v0.1.0
)
```

或者，通过`go get`命令获取远程代码包：

```
$ go get -u github.com/aliyun/alibabacloud-dkms-transfer-go-sdk
```

## 快速使用

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
	
	// 验证服务端证书
	ca, err := ioutil.ReadFile("path/to/caCert.pem")
	if err != nil {
		panic(err)
	}
	client.SetVerify(string(ca))
	// 或者，忽略服务端证书验证
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

## 许可证

[Apache-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Copyright (c) 2009-present, Alibaba Cloud All rights reserved.

package utils

import "github.com/aliyun/alibaba-cloud-sdk-go/sdk/log"

func CheckError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

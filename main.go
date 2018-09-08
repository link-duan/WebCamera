package main

import (
	"github.com/yaphper/WebCamera/app"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/log"
)

func main() {

	log.Fatal(app.StartServer())
}

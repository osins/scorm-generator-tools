package main

import (
	"flag"

	"github.com/wangsying/scorm-generator-tools/service"
)

var configYML = flag.String("c", "", "配置文件路径")
var outXML = flag.String("o", "", "配置文件输出路径")

func main() {
	flag.Parse()

	scormService := service.NewService()
	scormService.GenScorm2004(*configYML, *outXML)
}

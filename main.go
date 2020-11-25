package main

import (
	"flag"
	"regexp"
	"strings"

	"github.com/wangsying/scorm-generator-tools/service"
)

var confYML = flag.String("c", "", "配置文件路径")
var outXML = flag.String("o", "", "配置文件输出路径")
var dst = flag.String("d", "", "打包位置")

func main() {
	flag.Parse()

	if *dst == "" {
		re3, _ := regexp.Compile("(\\\\[^\\\\]+\\\\{0,}$)")
		*dst = re3.ReplaceAllString(*outXML, "") + "\\package.zip"
	}

	if strings.HasSuffix(*outXML, "\\") {
		*outXML = strings.TrimSuffix(*outXML, "\\")
	}

	scormService := service.NewService()
	scormService.GenScorm2004(*confYML, *outXML, *dst)
}

package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/wangsying/scorm-generator-tools/schema"
	"github.com/wangsying/scorm-generator-tools/schema/config"
	"github.com/wangsying/scorm-generator-tools/schema/scorm"
	"gopkg.in/yaml.v3"
)

var configYML = flag.String("c", "", "配置文件路径")
var outXML = flag.String("o", "", "配置文件输出路径")

func main() {
	flag.Parse()

	content := readYMLToScorm()

	scormToXML(content)
}

func readYMLToScorm() scorm.Scorm {
	//应该是 绝对地址
	yamlFile, err := ioutil.ReadFile(*configYML)
	if err != nil {
		fmt.Println("error 1: ", err.Error())
	}

	courses := config.Courses{}

	err = yaml.Unmarshal(yamlFile, &courses)
	if err != nil {
		fmt.Println("error 2:", err.Error())
	}

	scormService := scorm.NewService()
	scormContent := scormService.NewScorm2004()

	for _, value := range courses.Courses {

		re3, _ := regexp.Compile("[^a-zA-Z\\d]")
		value.Organization = "scorm_" + re3.ReplaceAllString(value.Organization, "_")

		fmt.Println("title: ", value.Title)
		fmt.Println("organization: ", value.Organization)

		organization := scorm.Organization{Title: value.Title, Identifier: value.Organization, Items: []scorm.Item{}}

		for wIdx, w := range value.Coursewares {
			fmt.Println("	idx: ", wIdx)
			fmt.Println("	title: ", w.Title)
			fmt.Println("	href: ", w.Href)

			itemID := value.Organization + "_item_" + strconv.Itoa(wIdx)
			scoID := itemID + "_sco"
			resID := itemID + "_res"
			organization.Items = append(organization.Items,
				scorm.Item{
					Title:         w.Title,
					Identifier:    itemID,
					Identifierref: scoID})

			homeResource := scorm.Resource{Identifier: scoID, Href: "index.html"}
			homeResource.Type = "sco"
			homeResource.Files = []scorm.File{scorm.File{Href: "index.html"}}
			homeResource.Dependency = scorm.Dependency{Identifierref: resID}
			scormContent.Resources = append(scormContent.Resources, homeResource)

			resource := scorm.Resource{Identifier: resID, Href: w.Href}
			resource.Type = "asset"
			resource.Files = loadResources(w.Href)
			scormContent.Resources = append(scormContent.Resources, resource)

			for i, f := range resource.Files {
				fmt.Printf("	resources[%d]: %s\n", i, f.Href)
			}
		}

		scormContent.Organizations = append(scormContent.Organizations, organization)
	}

	return scormContent
}

func loadResources(dir string) []scorm.File {
	resources := []scorm.File{}
	dir = strings.TrimRight(dir, "/")
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		sub := dir + "/" + f.Name()
		if f.IsDir() {
			result := loadResources(sub)
			for _, r := range result {
				resources = append(resources, r)
			}
		} else {
			resources = append(resources, scorm.File{Href: sub})
		}
	}

	return resources
}

func scormToXML(content scorm.Scorm) {
	v := new(schema.XMLManifestNode)
	v.Identifier = xml.Attr{Name: xml.Name{Local: "identifier"}, Value: "easymind_scorm_2004_course_generator"}
	v.Version = xml.Attr{Name: xml.Name{Local: "version"}, Value: "1"}
	v.Xmlns = xml.Attr{Name: xml.Name{Local: "xmlns"}, Value: "http://www.imsglobal.org/xsd/imscp_v1p1"}
	v.XmlnsXsi = xml.Attr{Name: xml.Name{Local: "xmlns:xsi"}, Value: "http://www.w3.org/2001/XMLSchema-instance"}
	v.XmlnsAdlcp = xml.Attr{Name: xml.Name{Local: "xmlns:adlcp"}, Value: "http://www.adlnet.org/xsd/adlcp_v1p3"}
	v.XmlnsAdlseq = xml.Attr{Name: xml.Name{Local: "xmlns:adlseq"}, Value: "http://www.adlnet.org/xsd/adlseq_v1p3"}
	v.XmlnsAdlnav = xml.Attr{Name: xml.Name{Local: "xmlns:adlnav"}, Value: "http://www.adlnet.org/xsd/adlnav_v1p3"}
	v.XmlnsImsss = xml.Attr{Name: xml.Name{Local: "xmlns:imsss"}, Value: "easymind_scorm_2004_course_generator"}
	v.XsiSchemaLocation = xml.Attr{Name: xml.Name{Local: "xsi:schemaLocation"}, Value: "http://www.imsglobal.org/xsd/imscp_v1p1 imscp_v1p1.xsd http://www.adlnet.org/xsd/adlcp_v1p3 adlcp_v1p3.xsd http://www.adlnet.org/xsd/adlseq_v1p3 adlseq_v1p3.xsd http://www.adlnet.org/xsd/adlnav_v1p3 adlnav_v1p3.xsd http://www.imsglobal.org/xsd/imsss imsss_v1p0.xsd"}

	v.MetadataNode.Schema = content.Metadata.Schema
	v.MetadataNode.Schemaversion = content.Metadata.Version

	v.OrganizationNode = schema.XMLOrganizationsNode{
		Default: xml.Attr{Name: xml.Name{Local: "default"},
			Value: string(content.Organizations[0].Identifier)}}

	for _, o := range content.Organizations {
		organization := schema.XMLOrganization{Title: o.Title, Items: []schema.XMLItemType{}}
		organization.Identifier = xml.Attr{
			Name:  xml.Name{Local: "identifier"},
			Value: string(o.Identifier)}

		for _, i := range o.Items {
			item := schema.XMLItemType{
				Title:         i.Title,
				Identifier:    xml.Attr{Name: xml.Name{Local: "identifier"}, Value: i.Identifier},
				Identifierref: xml.Attr{Name: xml.Name{Local: "identifierref"}, Value: i.Identifierref},
			}

			organization.Items = append(organization.Items, item)

			fmt.Println(i.Title)
			fmt.Println(item.Title)
		}

		v.OrganizationNode.Organizations = append(v.OrganizationNode.Organizations, organization)
	}

	v.ResourceNode = schema.XMLResourcesNode{Resource: []schema.XMLResourceType{}}
	for _, r := range content.Resources {
		resource := schema.XMLResourceType{
			Identifier: xml.Attr{Name: xml.Name{Local: "identifier"}, Value: r.Identifier},
			Type:       xml.Attr{Name: xml.Name{Local: "type"}, Value: "webcontent"},
			ScormType:  xml.Attr{Name: xml.Name{Local: "adlcp:scormType"}, Value: r.Type},
		}

		if r.Type == "sco" {
			resource.Href = xml.Attr{Name: xml.Name{Local: "href"}, Value: "index.html"}
		}

		for _, f := range r.Files {
			xmlFile := schema.XMLFileType{
				Href: xml.Attr{Name: xml.Name{Local: "href"}, Value: f.Href}}

			resource.Files = append(resource.Files, xmlFile)
		}

		v.ResourceNode.Resource = append(v.ResourceNode.Resource, resource)
	}

	f, err := os.Create(*outXML + "\\imsmanifest.xml")
	if err != nil {
		fmt.Println("文件创建失败", err.Error())
		return
	}
	defer f.Close()

	//序列化到文件中
	encoder := xml.NewEncoder(f)
	err = encoder.Encode(v)
	if err != nil {
		fmt.Println("编码错误：", err.Error())
		return
	}

	fmt.Println("编码成功:", *outXML+"\\imsmanifest.xml")
}

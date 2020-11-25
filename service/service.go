package service

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/wangsying/scorm-generator-tools/schema/config"
	"github.com/wangsying/scorm-generator-tools/schema/scorm"
	"github.com/wangsying/scorm-generator-tools/schema/sxml"
	"gopkg.in/yaml.v2"
)

// Service scorm service
type Service interface {
	NewScorm2004()
	GenScorm2004(confYML, outXML string, dst string)
}

type service struct {
	content scorm.Scorm
}

// NewService scorm service new
func NewService() Service {
	return &service{}
}

// GenScorm2004 生成scorm配置文件
func (s *service) GenScorm2004(confYML, outXML string, dst string) {
	s.readConfigAndGenScorm2004(confYML)
	s.scormToXML(outXML)
	s.zip(dst, outXML)
}

// New2004 创建scorm 2004标准课件
func (s *service) NewScorm2004() {
	s.content = scorm.Scorm{
		Metadata:      scorm.Metadata{Schema: "ADL SCORM", Version: "2004 3rd Edition"},
		Organizations: []scorm.Organization{},
		Resources:     []scorm.Resource{}}
}

// readConfigAndGenScorm2004 读取配置文件并生成scorm配置
func (s *service) readConfigAndGenScorm2004(configYML string) {
	//应该是 绝对地址
	yamlFile, err := ioutil.ReadFile(configYML)
	if err != nil {
		fmt.Println("error 1: ", err.Error())
	}

	courses := config.Courses{}

	err = yaml.Unmarshal(yamlFile, &courses)
	if err != nil {
		fmt.Println("error 2:", err.Error())
	}

	s.NewScorm2004()

	for _, value := range courses.Courses {

		re3, _ := regexp.Compile("[^a-zA-Z\\d]")
		value.Organization = "scorm_" + re3.ReplaceAllString(value.Organization, "_")

		organization := scorm.Organization{Title: value.Title, Identifier: value.Organization, Items: []scorm.Item{}}
		items, res := s.loadCoursewares(organization.Identifier+"_item", value.Coursewares)
		for _, r := range res {
			s.content.Resources = append(s.content.Resources, r)
		}

		organization.Items = items
		s.content.Organizations = append(s.content.Organizations, organization)
	}
}

// scormToXML 由scorm结构生成scorm的xml配置文件
func (s *service) scormToXML(outXML string) {
	v := new(sxml.XMLManifestNode)
	v.Identifier = xml.Attr{Name: xml.Name{Local: "identifier"}, Value: "easymind_scorm_2004_course_generator"}
	v.Version = xml.Attr{Name: xml.Name{Local: "version"}, Value: "1"}
	v.Xmlns = xml.Attr{Name: xml.Name{Local: "xmlns"}, Value: "http://www.imsglobal.org/xsd/imscp_v1p1"}
	v.XmlnsXsi = xml.Attr{Name: xml.Name{Local: "xmlns:xsi"}, Value: "http://www.w3.org/2001/XMLSchema-instance"}
	v.XmlnsAdlcp = xml.Attr{Name: xml.Name{Local: "xmlns:adlcp"}, Value: "http://www.adlnet.org/xsd/adlcp_v1p3"}
	v.XmlnsAdlseq = xml.Attr{Name: xml.Name{Local: "xmlns:adlseq"}, Value: "http://www.adlnet.org/xsd/adlseq_v1p3"}
	v.XmlnsAdlnav = xml.Attr{Name: xml.Name{Local: "xmlns:adlnav"}, Value: "http://www.adlnet.org/xsd/adlnav_v1p3"}
	v.XmlnsImsss = xml.Attr{Name: xml.Name{Local: "xmlns:imsss"}, Value: "easymind_scorm_2004_course_generator"}
	v.XsiSchemaLocation = xml.Attr{Name: xml.Name{Local: "xsi:schemaLocation"}, Value: "http://www.imsglobal.org/xsd/imscp_v1p1 imscp_v1p1.xsd http://www.adlnet.org/xsd/adlcp_v1p3 adlcp_v1p3.xsd http://www.adlnet.org/xsd/adlseq_v1p3 adlseq_v1p3.xsd http://www.adlnet.org/xsd/adlnav_v1p3 adlnav_v1p3.xsd http://www.imsglobal.org/xsd/imsss imsss_v1p0.xsd"}

	v.MetadataNode.Schema = s.content.Metadata.Schema
	v.MetadataNode.Schemaversion = s.content.Metadata.Version

	v.OrganizationNode = sxml.XMLOrganizationsNode{
		Default: xml.Attr{Name: xml.Name{Local: "default"},
			Value: string(s.content.Organizations[0].Identifier)}}

	for _, o := range s.content.Organizations {
		organization := sxml.XMLOrganization{Title: o.Title, Items: []sxml.XMLItemNode{}}
		organization.Identifier = xml.Attr{
			Name:  xml.Name{Local: "identifier"},
			Value: string(o.Identifier)}

		organization.Items = s.loadItems(o.Items)

		v.OrganizationNode.Organizations = append(v.OrganizationNode.Organizations, organization)
	}

	v.ResourceNode = sxml.XMLResourcesNode{Resource: []sxml.XMLResourceType{}}
	for _, r := range s.content.Resources {
		resource := sxml.XMLResourceType{
			Identifier: xml.Attr{Name: xml.Name{Local: "identifier"}, Value: r.Identifier},
			Type:       xml.Attr{Name: xml.Name{Local: "type"}, Value: "webcontent"},
			ScormType:  xml.Attr{Name: xml.Name{Local: "adlcp:scormType"}, Value: r.Type},
		}

		if r.Type == "sco" {
			resource.Href = xml.Attr{Name: xml.Name{Local: "href"}, Value: r.Href}
		}

		if r.Dependency != (scorm.Dependency{}) {
			resource.Dependency = []sxml.XMLDependencyType{{
				Identifierref: xml.Attr{
					Name:  xml.Name{Local: "identifierref"},
					Value: r.Dependency.Identifierref}}}
		}

		for _, f := range r.Files {
			xmlFile := sxml.XMLFileType{
				Href: xml.Attr{Name: xml.Name{Local: "href"}, Value: f.Href}}

			resource.Files = append(resource.Files, xmlFile)
		}

		v.ResourceNode.Resource = append(v.ResourceNode.Resource, resource)
	}

	f, err := os.Create(outXML + "\\imsmanifest.xml")
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
}

func (s *service) loadCoursewares(prefix string, coursewares []config.Courseware) ([]scorm.Item, []scorm.Resource) {
	items := []scorm.Item{}
	resources := []scorm.Resource{}

	for wIdx, w := range coursewares {
		itemID := prefix + "_" + strconv.Itoa(wIdx)
		scoID := itemID + "_sco"
		resID := itemID + "_res"

		item := scorm.Item{
			Title:         w.Title,
			Identifier:    itemID,
			Identifierref: scoID}
		// 获取目录前缀
		re3, _ := regexp.Compile("[^\\\\]*$")
		indexHomeHref := re3.FindString(w.Href) + "/index.html"

		if len(w.Coursewares) > 0 {
			itemItems, itemChildren := s.loadCoursewares(itemID, w.Coursewares)
			for _, c := range itemChildren {
				item.Items = itemItems
				resources = append(resources, c)
			}
		} else {
			homeResource := scorm.Resource{Identifier: scoID, Href: indexHomeHref}
			homeResource.Type = "sco"
			homeResource.Files = []scorm.File{scorm.File{Href: indexHomeHref}}
			homeResource.Dependency = scorm.Dependency{Identifierref: resID}

			resource := scorm.Resource{Identifier: resID, Href: w.Href}
			resource.Type = "asset"
			resource.Files = s.loadResources(w.Href)
			resources = append(resources, homeResource)
			resources = append(resources, resource)
		}

		items = append(items, item)
	}

	return items, resources
}

// loadItems 载入课件项目
func (s *service) loadItems(items []scorm.Item) []sxml.XMLItemNode {
	result := []sxml.XMLItemNode{}

	for _, i := range items {
		item := sxml.XMLItemNode{
			Title:      i.Title,
			Identifier: xml.Attr{Name: xml.Name{Local: "identifier"}, Value: i.Identifier},
		}

		if len(i.Items) > 0 {
			item.Items = s.loadItems(i.Items)
		} else {
			item.Identifierref = xml.Attr{Name: xml.Name{Local: "identifierref"}, Value: i.Identifierref}
		}

		result = append(result, item)
	}

	return result
}

// loadResources 读取课件资源(目录和文件列表)
func (s *service) loadResources(dir string) []scorm.File {
	resources := []scorm.File{}
	dir = strings.TrimRight(dir, "/")
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		sub := dir + "/" + f.Name()
		if f.IsDir() {
			result := s.loadResources(sub)
			for _, r := range result {
				resources = append(resources, r)
			}
		} else if sub != "" {
			re3, err := regexp.Compile("(.*\\\\)")
			if err != nil {
				fmt.Println("regexp error: ", err)
			} else {
				if re3 != nil {
					href := re3.ReplaceAllString(sub, "")
					resources = append(resources, scorm.File{Href: href})
				}
			}
		}
	}

	return resources
}

func (s *service) zip(dst, src string) (err error) {
	// 创建准备写入的文件
	fw, err := os.Create(dst)
	defer fw.Close()
	if err != nil {
		return err
	}

	// 通过 fw 来创建 zip.Write
	zw := zip.NewWriter(fw)
	defer func() {
		// 检测一下是否成功关闭
		if err := zw.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// 下面来将文件写入 zw ，因为有可能会有很多个目录及文件，所以递归处理
	return filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) (err error) {
		if path == src {
			return
		}

		if errBack != nil {
			return errBack
		}

		// 通过文件信息，创建 zip 的文件信息
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return
		}

		// 替换文件信息中的文件名
		fh.Name = strings.ReplaceAll(strings.TrimPrefix(path, src+"\\"), "\\", "/")

		// 这步开始没有加，会发现解压的时候说它不是个目录
		if fi.IsDir() {
			fh.Name += "/"
		}

		// 写入文件信息，并返回一个 Write 结构
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return
		}

		// 检测，如果不是标准文件就只写入头信息，不写入文件数据到 w
		// 如目录，也没有数据需要写
		if !fh.Mode().IsRegular() {
			return nil
		}

		// 打开要压缩的文件
		fr, err := os.Open(path)
		defer fr.Close()
		if err != nil {
			return
		}

		// 将打开的文件 Copy 到 w
		io.Copy(w, fr)

		return
	})
}

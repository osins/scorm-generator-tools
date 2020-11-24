package scorm

// Resource 课件资源
type Resource struct {
	Identifier string
	Files      []File
	Href       string
	Type       string
	Dependency Dependency
}

// File 资源文件
type File struct {
	Href string
}

// Dependency 资源依赖
type Dependency struct {
	Identifierref string
}

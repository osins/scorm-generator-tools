package scorm

// Organization 生产课件的组织信息
type Organization struct {
	Identifier string
	Title      string
	Items      []Item
}

// Item 课件项目
type Item struct {
	Identifier    string
	Identifierref string
	Title         string
	Items         []Item
}

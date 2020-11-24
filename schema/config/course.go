package config

// Courses 课程列表
type Courses struct {
	Courses []Course `yaml:"courses"`
}

// Course 课程
type Course struct {
	Title        string       `yaml:"title"`
	Organization string       `yaml:"organization"`
	Coursewares  []Courseware `yaml:"coursewares"`
}

// Courseware 课件
type Courseware struct {
	Title       string       `yaml:"title"`
	Href        string       `yaml:"href"`
	Coursewares []Courseware `yaml:"coursewares"`
}

package scorm

// Scorm 课件标准
type Scorm struct {
	Metadata      Metadata
	Organizations []Organization
	Resources     []Resource
}

package scorm

// Service scorm service
type Service interface {
	NewScorm2004() Scorm
}

type service struct {
}

// Scorm 课件标准
type Scorm struct {
	Metadata      Metadata
	Organizations []Organization
	Resources     []Resource
}

// NewService scorm service new
func NewService() Service {
	return &service{}
}

// New2004 创建scorm 2004标准课件
func (*service) NewScorm2004() Scorm {
	metadata := Metadata{
		Schema:  "ADL SCORM",
		Version: "2004 3rd Edition"}
	organizations := []Organization{}
	resources := []Resource{}

	return Scorm{
		Metadata:      metadata,
		Organizations: organizations,
		Resources:     resources}
}

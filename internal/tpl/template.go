package tpl

const (
	TemplateTypeInstance = iota
	TemplateTypeJamulus
)

type CreateInstanceTemplate struct {
	Template struct {
		TemplateName        string
		TemplateDescription string
		TemplateType        uint64
		IsTemplate          bool `json:"-"`
	}
	Instance struct {
		InstanceType    string
		KeyPair         string
		KeyPairPath     string
		SecurityGroupID string
	}
	Jamulus struct {
		CentralServer        string
		ServerInfo           string
		FastUpdate           bool
		LogPath              string
		HTMLStatusFile       string
		RecordingPath        string
		NoRecordOnStart      bool
		EnableMultiThreading bool
		MaxUsers             int
		WelcomeMessage       string
	}
}

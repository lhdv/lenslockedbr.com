package views

const (
	AlertLvlError = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo = "info"
	AlertLvlSuccess = "success"
)

// Data is the top level structure that views expect data to come in.
type Data struct {
	Alert *Alert
	Yield interface{}
}

// Alert is used to render Bootstrap Alert messages in templates
type Alert struct {
	Level string
	Message string
}

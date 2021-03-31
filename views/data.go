package views

var (
	AlertLevelError = "danger"
	AlertLevelWarning = "warning"
	AlertLevelInfo = "info"
	AlertLevelSuccess = "success"

	AlertMessageGeneric = "Something went wrong!"
)

type Alert struct {
	Level string
	Message string
}

type Data struct {
	Alert *Alert
	Content interface{}
}

func (d *Data) SetAlert(err error) {
	if pErr, ok := err.(PublicError); ok {
		d.Alert = &Alert{
			Level: AlertLevelError,
			Message: pErr.Public(),
		}
	} else {
		d.Alert = &Alert{
			Level: AlertLevelError,
			Message: AlertMessageGeneric,
		}
	}
}

func (d *Data) AlertError(message string) {
	d.Alert = &Alert{
		Level: AlertLevelError,
		Message: message,
	}
}

type PublicError interface {
	error
	Public() string
}

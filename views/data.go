package views

import (
	"log"
	"net/http"
	"time"
)

var (
	AlertLevelError   = "danger"
	AlertLevelWarning = "warning"
	AlertLevelInfo    = "info"
	AlertLevelSuccess = "success"
	
	AlertMessageGeneric = "Something went wrong!"
)

type Alert struct {
	Level   string
	Message string
}

type Data struct {
	Alert   *Alert
	User    interface{}
	Content interface{}
}

func (d *Data) SetAlert(err error) {
	if pErr, ok := err.(PublicError); ok {
		d.Alert = &Alert{
			Level:   AlertLevelError,
			Message: pErr.Public(),
		}
	} else {
		log.Println(err)
		d.Alert = &Alert{
			Level:   AlertLevelError,
			Message: AlertMessageGeneric,
		}
	}
}

func (d *Data) AlertError(message string) {
	d.Alert = &Alert{
		Level:   AlertLevelError,
		Message: message,
	}
}

func (d *Data) AlertSuccess(message string) {
	d.Alert = &Alert{
		Level:   AlertLevelSuccess,
		Message: message,
	}
}

func (d *Data) AlertInfo(message string) {
	d.Alert = &Alert{
		Level:   AlertLevelInfo,
		Message: message,
	}
}

func (d *Data) AlertWarning(message string) {
	d.Alert = &Alert{
		Level:   AlertLevelWarning,
		Message: message,
	}
}

type PublicError interface {
	error
	Public() string
}

func persistAlert(w http.ResponseWriter, alert Alert) {
	lvl, msg := getCookiesForAlert(alert, time.Now().Add(5 * time.Minute))
	http.SetCookie(w, lvl)
	http.SetCookie(w, msg)
}

func clearAlert(w http.ResponseWriter) {
	alert := Alert{
		Level: "",
		Message: "",
	}
	lvl, msg := getCookiesForAlert(alert, time.Now())
	http.SetCookie(w, lvl)
	http.SetCookie(w, msg)
}

func getAlert(req *http.Request) *Alert {
	lvl, err := req.Cookie("alert_level")
	if err != nil {
		return nil
	}
	msg, err := req.Cookie("alert_message")
	if err != nil {
		return nil
	}

	alert := &Alert{
		Level: lvl.Value,
		Message: msg.Value,
	}
	return alert
}

func RedirectAlert(w http.ResponseWriter, req *http.Request, urlStr string, code int, alert Alert) {
	persistAlert(w, alert)
	http.Redirect(w, req, urlStr, code)
}

func getCookiesForAlert(alert Alert, expires time.Time) (*http.Cookie, *http.Cookie) {
	lvl := &http.Cookie{
		Name: "alert_level",
		Value: alert.Level,
		Expires: expires,
		HttpOnly: true,
	}
	msg := &http.Cookie{
		Name: "alert_message",
		Value: alert.Message,
		Expires: expires,
		HttpOnly: true,
	}
	return lvl, msg
}

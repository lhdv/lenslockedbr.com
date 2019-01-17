package views

import (
	"log"
	"net/http"
	"time"

	"lenslockedbr.com/models"
)

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"

	// AlertMsgGeneric is displayed when any random error is
	// encountered by our backend.
	AlertMsgGeneric = "Something went wrong. Please try again, " +
		"and contact us if the problem persists."
)

type PublicError interface {
	error
	Public() string
}

// Alert is used to render Bootstrap Alert messages in templates
type Alert struct {
	Level   string
	Message string
}

// Data is the top level structure that views expect data to come in.
type Data struct {
	Alert *Alert
	User  *models.User
	Yield interface{}
}

func (d *Data) SetAlert(err error) {
	var msg string

	if pErr, ok := err.(PublicError); ok {
		msg = pErr.Public()
	} else {
		log.Println(err)
		msg = AlertMsgGeneric
	}

	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

// RedirectAlert accepts all the normal params for an http.Redirect
// and performs a redirect, but only after persisting the provided
// alert in a cookie so that it can be displayed when the new page
// is loaded
func RedirectAlert(w http.ResponseWriter, r *http.Request, urlStr string,
	code int, alert Alert) {
	persistAlert(w, alert)
	http.Redirect(w, r, urlStr, code)
}

/////////////////////////////////////////////////////////////////////
//
// Helper methods
//
/////////////////////////////////////////////////////////////////////

func persistAlert(w http.ResponseWriter, alert Alert) {

	expiresAt := time.Now().Add(5 * time.Minute)
	lvl := http.Cookie{
		Name:     "alert_level",
		Value:    alert.Level,
		Expires:  expiresAt,
		HttpOnly: true,
	}

	msg := http.Cookie{
		Name:     "alert_message",
		Value:    alert.Message,
		Expires:  expiresAt,
		HttpOnly: true,
	}

	http.SetCookie(w, &lvl)
	http.SetCookie(w, &msg)
}

func clearAlert(w http.ResponseWriter) {

	lvl := http.Cookie{
		Name:     "alert_level",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}

	msg := http.Cookie{
		Name:     "alert_message",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}

	http.SetCookie(w, &lvl)
	http.SetCookie(w, &msg)
}

func getAlert(r *http.Request) *Alert {

	lvl, err := r.Cookie("alert_level")
	if err != nil {
		return nil
	}

	msg, err := r.Cookie("alert_message")
	if err != nil {
		return nil
	}

	alert := Alert{
		Level:   lvl.Value,
		Message: msg.Value,
	}

	return &alert
}

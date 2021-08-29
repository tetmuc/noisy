package alert

type IAlert interface {
	AlertText(keyWord string, title, msg string, nominees ...string) error
	AsyncAlertText(keyWord string, title, msg string, nominees ...string)
}

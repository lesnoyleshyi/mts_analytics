package interfaces

import "net/http"

// shrink to one method?
type Responder interface {
	RespondSuccess(w http.ResponseWriter, msg string, status int)
	RespondError(w http.ResponseWriter, msg string, status int)
}

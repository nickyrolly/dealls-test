package middleware

import (
	"net/http"
)

func credentialCheck(r *http.Request, email string, password string) (user UserSession, err error) {
	// user, err = authRepo.ValidateAuthenticationLogin(r.Context(), passcode, password)
	// if err != nil {
	// 	log.Println("[credential check] error", err)
	// }
	return user, err
}

func profileCheck(r *http.Request, xUserID string) (user UserSession, err error) {

	// participantID, err := strconv.ParseInt(xUserID, 10, 64)
	// if err != nil {
	// 	return user, err
	// }

	// user, err = authRepo.ValidateProfile(r.Context(), participantID)
	// if err != nil {
	// 	log.Println("[credential check] error", err)
	// }
	return user, err
}

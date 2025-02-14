package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/context"
	"github.com/nickyrolly/dealls-test/common"
)

const ALLOWED_HEADERS = "Cookie, Origin, Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, X-UserID, X-Userid"

func CorsMiddleware(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", ALLOWED_HEADERS)
		w.Header().Set("Access-Control-Expose-Headers", "*")

		nextHandler.ServeHTTP(w, r)
	})
}

func BasicMiddleware(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, password, ok := r.BasicAuth()
		if !ok {
			common.CustomResponseAPI(w, r, http.StatusUnauthorized, map[string]interface{}{"success": false, "error_message": "Unauthorized"})
			return
		}

		if email == "" {
			log.Println("[BasicMiddleware] Empty Email")
			common.CustomResponseAPI(w, r, http.StatusBadRequest, map[string]interface{}{"success": false, "error_message": "email cannot be empty"})
			return
		}

		if password == "" {
			log.Println("[BasicMiddleware] Empty Password")
			common.CustomResponseAPI(w, r, http.StatusBadRequest, map[string]interface{}{"success": false, "error_message": "email cannot be empty"})
			return
		}

		var user UserSession
		// userIdShort := strings.Replace(passcode, "-", "", -1)
		// userIdInt, err := strconv.ParseInt(userIdShort, 10, 64)
		// if err != nil {
		// 	log.Println("[BasicMiddleware] Passcode parse error")
		// 	common.CustomResponseAPI(w, r, http.StatusBadRequest, map[string]interface{}{"success": false, "error_message": "passcode parse error"})
		// 	return
		// }
		user.Email = email
		user.Password = password

		context.Set(r, "user", user)
		defer context.Clear(r)

		nextHandler.ServeHTTP(w, r)
	})
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("key_exam_x_user")

func AuthenticationMiddleware(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		// var isCredentialAuthentication bool
		user := UserSession{}

		// Get the JSON body and decode into credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			// If the structure of the body is wrong, return an HTTP error
			log.Println("[AuthenticationMiddleware] Cred Decoder err", err)
			common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{"success": false, "error_message": "Pastikan nomor pendaftaran / level pendidikan / password Anda benar."})
			return
		}

		// passcodeShort := strings.Replace(creds.PassCode, "-", "", -1)
		// passcodeInt, err := strconv.ParseInt(passcodeShort, 10, 64)
		// if err != nil {
		// 	log.Println("[AuthenticationMiddleware] Cred Decoder err creds.PassCode: "+creds.PassCode, err)

		// 	common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{"success": false, "error_message": "Pastikan nomor pendaftaran / level pendidikan / password Anda benar."})

		// 	return
		// }

		user, err = credentialCheck(r, creds.Email, creds.Password)
		if err != nil {
			log.Println("[AuthenticationMiddleware] user credential error, passcodeInt: "+creds.Email, err)

			common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{"success": false, "error_message": "Pastikan nomor pendaftaran / level pendidikan / password Anda benar."})
			return
		}

		examinationTime := time.Now()
		expirationTime := examinationTime.Add(360 * time.Minute)
		claims := &Claims{
			Email: creds.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			log.Println("[AuthenticationMiddleware] Token Signed err : ", err)
			common.CustomResponseAPI(w, r, http.StatusUnauthorized, map[string]interface{}{"success": false, "error_message": "Unauthorized"})
			return
		}

		// if staging
		// http.SetCookie(w, &http.Cookie{
		// 	Name:     "X_STP",
		// 	Value:    tokenString,
		// 	Expires:  expirationTime,
		// 	HttpOnly: true,
		// 	Secure:   true,
		// 	SameSite: http.SameSiteNoneMode,
		// })

		user.SessionToken = tokenString
		user.SessionExpirationTime = expirationTime
		context.Set(r, "user", user)
		defer context.Clear(r)

		// http.SetCookie(w, &http.Cookie{
		// 	Name:    "X_STP",
		// 	Value:   tokenString,
		// 	Expires: expirationTime,
		// })

		nextHandler.ServeHTTP(w, r)

	})
}

package middlewares

import (
	"strings"
	"net/http"
	"crypto/rsa"
	"encoding/json"

	"github.com/unrolled/render"
	"github.com/zenazn/goji/web"
	"github.com/dvsekhvalnov/jose2go"

	"github.com/YagoCarballo/kumquat-academy-api/tools"
	"github.com/YagoCarballo/kumquat-academy-api/database/models"
)

type Permissions struct {
	Read  bool
	Write bool
}

type checkToken struct {
	h           func(web.C, http.ResponseWriter, *http.Request)
	permissions *Permissions
	privateKey	*rsa.PrivateKey
	publicKey	*rsa.PublicKey
}

func (e checkToken) checkPermisions(c web.C, w http.ResponseWriter, r *http.Request) {
	// TODO: Check Permissions for Restricted Areas
	if e.permissions.Write {
		c.Env["permissions"] = &e.permissions
		e.h(c, w, r)
		return
	}

	renderer := render.New()
	renderer.JSON(w, http.StatusForbidden, map[string]interface{}{
		"error":   "AccessDenied",
		"message": "Not enough permissions to access this area.",
	})
}

func (e checkToken) ServeHTTPC(context web.C, writter http.ResponseWriter, request *http.Request) {
	var authorization string
	//	deviceId := r.Header.Get("Device") // TODO: Check Device ID matches token

	cookie, err := request.Cookie("token")
	if err != nil {
		TriggerUnautorizedError(writter); return
	} else if (err != http.ErrNoCookie) {
		cookieData, err := e.parseCookie(cookie)
		if err != nil {
			TriggerUnautorizedError(writter); return
		}

		authorization = cookieData.AccessToken
		context.Env["token"] = cookieData
	}

	if authorization == "" && strings.HasPrefix(authorization, "Bearer ") == false {
		TriggerUnautorizedError(writter); return

	} else {
		token := strings.Replace(authorization, "Bearer ", "", 1)
		session, err := models.DBSession.FindSession(token)

		if err != nil || session == nil {
			TriggerUnautorizedError(writter); return
		} else {
			context.Env["session"] = &session

			// If URL is restricted, check permissions as well
			if e.permissions != nil {
				e.checkPermisions(context, writter, request)

			} else {
				e.h(context, writter, request)
			}
		}
	}
}

func TriggerUnautorizedError(writter http.ResponseWriter) {
	renderer := render.New()
	renderer.JSON(writter, http.StatusUnauthorized, map[string]interface{}{
		"error": "Unautorized",
	})
	return;
}

func (e checkToken) parseCookie(cookie *http.Cookie) (*tools.JWTSession, error) {
	var cookieData tools.JWTSession

	// Decodes the Encripted Cookie JWT Token
	payload, _, err := jose.Decode(cookie.Value, e.privateKey)
	if err != nil {
		return &cookieData, err
	}

	// Parses the Decoded JSON
	err = json.Unmarshal([]byte(payload), &cookieData)
	if err != nil {
		return &cookieData, err
	}

	return &cookieData, nil
}

func CheckSession(h func(web.C, http.ResponseWriter, *http.Request), privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) web.Handler {
	return checkToken{h, nil, privateKey, publicKey}
}

func Restricted(permissions Permissions, h func(web.C, http.ResponseWriter, *http.Request), privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) web.Handler {
	return checkToken{h, &permissions, privateKey, publicKey}
}

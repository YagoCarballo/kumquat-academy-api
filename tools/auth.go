package tools

import (
	"strconv"
	"net/http"

	. "github.com/YagoCarballo/kumquat-academy-api/constants"
)

func ParseID(rawId string) (uint32, int, map[string]interface{}) {
	id, err := strconv.ParseUint(rawId, 10, 32)
	if err != nil {
		return 0, http.StatusConflict, map[string]interface{}{
			"error":   "InvalidData",
			"message": "The data provided is invalid",
		}
	}

	return uint32(id), http.StatusOK, nil
}

func VerifyAccess(id interface{}, userId uint32, action AccessRight, isActionPermitted func(uint32, interface{}, AccessRight) bool) (int, map[string]interface{}) {
	allowed := isActionPermitted(userId, id, action)

	if allowed {
		return http.StatusOK, nil
	}

	return http.StatusForbidden, map[string]interface{}{
		"error":   "AccessDenied",
		"message": "You don't have enough access rights",
	}
}
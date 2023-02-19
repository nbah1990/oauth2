package handlers

import (
	"encoding/base64"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"oauth-proxy/internal/oauth-proxy/enums"
	"oauth-proxy/internal/oauth-proxy/infrastructure/repositories"
	"strings"
)

type ClientAuthenticator struct {
	ClientRepository repositories.ClientRepositoryI
}

func (ca ClientAuthenticator) ExtractBasicAuthCredentials(c echo.Context) (clientID string, secret string, err error) {
	auth := c.Request().Header.Get("Authorization")
	if auth == "" {
		err = fmt.Errorf("authorization header is empty")
		return
	}

	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Basic" {
		err = fmt.Errorf("authorization header format is invalid")
		return
	}

	b, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		err = fmt.Errorf("failed to decode base64-encoded string")
		return
	}

	cred := strings.SplitN(string(b), ":", 2)
	if len(cred) != 2 {
		err = fmt.Errorf("credentials format is invalid")
		return
	}

	clientID, secret = cred[0], cred[1]
	return
}

func (ca ClientAuthenticator) AuthenticateClient(clientID string, secret string, scope enums.GlobalClientScope) bool {
	client, err := ca.ClientRepository.GetByID(clientID)
	if err != nil {
		logrus.Warning(`unsuccessful access attempt`)
		return false
	}

	if !client.CheckSecret(secret) {
		logrus.Warning(`unsuccessful access attempt`)
		return false
	}

	if !client.HasScope(string(scope)) {
		logrus.Warning(`unsuccessful access attempt`)
		return false
	}

	return true
}

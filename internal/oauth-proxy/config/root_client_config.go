package config

import (
	"github.com/google/uuid"
)

type RootClientConfig struct {
	Enabled  bool
	ClientID string
	Secret   string
}

func NewRootClientConfig() RootClientConfig {
	rcc := RootClientConfig{}

	rce := GetEnv(`ROOT_CLIENT_ENABLED`, ``) == TrueString
	rcc.Enabled = rce

	cID := GetEnv(`ROOT_CLIENT_ID`, ``)
	if len(cID) == 0 {
		panic(`ROOT_CLIENT_ID can not be an empty string' `)
	}

	if !isValidUUID(cID) {
		panic(`ROOT_CLIENT_ID should be a valid UUID`)
	}

	cSec := GetEnv(`ROOT_CLIENT_SECRET`, ``)
	if len(cID) == 0 {
		panic(`ROOT_CLIENT_SECRET can not be an empty string' `)
	}

	return RootClientConfig{
		Enabled:  GetEnv(`ROOT_CLIENT_ENABLED`, ``) == `true`,
		ClientID: cID,
		Secret:   cSec,
	}
}

func isValidUUID(str string) bool {
	_, err := uuid.Parse(str)
	return err == nil
}

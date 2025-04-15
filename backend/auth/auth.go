package auth

import (
	"fmt"
	"gotth/template/backend/configuration"
	"gotth/template/backend/store"
	"net/http"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

type Casdoor struct {
	Client *casdoorsdk.Client
}

var casdoor Casdoor

func GetCasdoor() Casdoor {
	return casdoor
}

func InitCasdoor(cfg configuration.Configutration) {
	casdoorClient := casdoorsdk.NewClient(
		cfg.CasdoorEndpoint,
		cfg.CasdoorClientID,
		cfg.CasdoorClientSecret,
		cfg.CasdoorCertificate,
		cfg.CasdoorOrganization,
		cfg.CasdoorApplication,
	)

	casdoor = Casdoor{
		Client: casdoorClient,
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) (*casdoorsdk.Claims, error) {
	store := store.GetStore()
	session, err := store.GetToken(r)
	if err != nil {
		return nil, err
	}

	token, ok := session.Values["token"].(string)
	if !ok || token == "" {
		return nil, fmt.Errorf("Token does not exist")
	}
	c := GetCasdoor()
	return c.Client.ParseJwtToken(token)
}

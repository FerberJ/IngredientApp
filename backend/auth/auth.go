package auth

import (
	"fmt"
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

func InitCasdoor() {
	casdoorClient := casdoorsdk.NewClient(
		"http://localhost:8000",
		"a26be6b9115770948966",
		"8182163c5cfdbc9a641e6bb547d469304123ad76",
		`-----BEGIN CERTIFICATE-----
MIIE3TCCAsWgAwIBAgIDAeJAMA0GCSqGSIb3DQEBCwUAMCgxDjAMBgNVBAoTBWFk
bWluMRYwFAYDVQQDEw1jZXJ0LWJ1aWx0LWluMB4XDTI1MDIyMTE3NTA0NFoXDTQ1
MDIyMTE3NTA0NFowKDEOMAwGA1UEChMFYWRtaW4xFjAUBgNVBAMTDWNlcnQtYnVp
bHQtaW4wggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQC/BzJkbT5tK1ag
obMRVm7my6+m1O0P6UXT5Vn+zvWvnWND3/hEseD2SDTv628pARXqvJ+v7nirxbi/
quqxHAEno57zxGX5j+AT7cSNvaj9MuFbib47vmjhaeOo7lQDlNKsbwI6Ggp2XFQL
qWaZDJXei3PwTUw7gTEJ5XvaQS+h319al2E9gXyNuEt6Mu0YgJ1BeuzCdAQb/YlU
3coKdFE9CfNZDoerQj7wV7kK/O2OGBIiKg6JGaGeGXwOHGijNBgkjYMtnBEWSoEf
wKUYx/fk8hf+tQ8X4nbqAt1xYd7Z8yk8i1raguKmwknmo6oGFLHJRqVt9fJ5fJr4
QqEgzPn/QEsF/V4vliXb6Ss/6BPOSYi3JGR550u4J9H0LR5W0i4CfSkxg/viLDnc
QYqqSV9ELucZ8OpSeH0S0/OMwZUUekt9mJbKLJxXEWw1Z/y9a959De1nOtAdELqT
CuuHPzQd3f61b4T6GGqGZeX7SL5ZusGu0fJqtPRpSU4aC+yhiBzuQxiH5wDkB9Ai
1UZcrS8TwMo5R3NhxmCQtnkSduVYacV4loZvZOgBhe/SffJDIKPSVlKwFm3VGhAk
6GENHFEFZBqkWujAOzc6xnnT8dJtis3CZXD7uvutP/sfkTtX/jR/bAsbXfmwhOJH
/4zHipC952ftOMxblScoWOhWk80NCwIDAQABoxAwDjAMBgNVHRMBAf8EAjAAMA0G
CSqGSIb3DQEBCwUAA4ICAQAGNKIESa8/9IPI9AR11WaeIvrdwnEHjXziZq4mastb
uVkb+rlwATRv2VXgSn0SbbFLnaDjtui7yKXq9CPyG/nOg9n+hhdqzfL+BmwULN/h
nFXPud0yIXFRPFMTrSvK4yj7mN06iu07SoaE8AhTcWSy545mHBCl+77cSXf38mFc
jQYdrMUBvc3FtSR7ml4I/MJGbQ25TqIAgSzta0xwSTIGp22UsqQB91dpirm7t29/
sxXfG1XZ0Qah1VUenQ5QJ3LwfB1RU++M8YR9v4t8UdULIAfyvlfPk7cP03IAJoqF
ZGEf9+ElWKYaGTiFji5DWbxSKdC8ey6TxpQONh49uGeD43/Rxqbf98PBO9zwla9i
1GSCX4z9veqx6mb7iR/DXH1As3ybXcClJg9LDx1NLi9m0zQb2+TkDcMfHQausvXX
CpIhU10VF4Ep5mAX+Xt/VlZeChg+fy4j51SFNHFAxevPLdB6qYPmcRljlm8y4uyk
zMtttv6vl07pGzHrRiNnUn71cMj2KCSlEgBIPdTaW3YSGtpvOIia67BRYeWadl+u
HlUvpwL/aNXbz/xy9fnvvL/JY9fbkqM+QAC76QhoKN+8iyJgrHrFCnKIY/bSSso3
+8Nng1vfj11jeM2pzZBDKcIc6nwF6eXkZ62HabD3HKiamtTjDZE1Ee2tszWBEMoP
ow==
-----END CERTIFICATE-----`,
		"recipe",
		"recipe",
	)

	casdoor = Casdoor{
		Client: casdoorClient,
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) (*casdoorsdk.Claims, error) {
	store := store.GetStore()
	session, err := store.Get(r, "session")
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

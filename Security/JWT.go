package security

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/SermoDigital/jose"
	"github.com/SermoDigital/jose/jws"
	"github.com/dgrijalva/jwt-go"
	models "github.com/flow_server/Models"
	utils "github.com/flow_server/Utils"
)

/*ValidateJwt validates incoming jwt from client and server*/
func ValidateJwt(jwtToken string) (models.JwtBodyStruct, error) {
	splitJWT := strings.Split(jwtToken, ".")
	jwtHeader, _ := (jose.Base64Decode([]byte(splitJWT[0])))
	jwtBody, _ := (jose.Base64Decode([]byte(splitJWT[1])))
	// jwtSigniture, _ := (jose.Base64Decode([]byte(splitJWT[2])))

	headerStruct := models.JwtHeaderStruct{}
	jwtBodyStruct := models.JwtBodyStruct{}
	headerError := json.Unmarshal(jwtHeader, &headerStruct)
	bodyError := json.Unmarshal(jwtBody, &jwtBodyStruct)

	utils.ErrorNilCheck(headerError)
	utils.ErrorNilCheck(bodyError)

	bodyValidated := validateAud(jwtBodyStruct)

	if !bodyValidated {
		return jwtBodyStruct, errors.New("Invalid jwtBodyStruct")
	}

	return jwtBodyStruct, nil
}

func validateAud(jwtBodyStruct models.JwtBodyStruct) bool {
	valid := false
	switch jwtBodyStruct.Iss {
	case "https://accounts.google.com":
		{
			if jwtBodyStruct.Aud != "178079389303-jdtfifkob381duk64fuppqp8004gk4u7.apps.googleusercontent.com" {
				break
			}
			valid = true
		}
	default:
		break
	}

	return valid
}

/*//||
//jwtBodyStruct.Aud != "178079389303-jdtfifkob381duk64fuppqp8004gk4u7.apps.googleusercontent.com"*/
//"1004547079072-9hcr020dn0bnsvbfusgisad3iu28chk5.apps.googleusercontent.com"
var mySignedKey = []byte("mysupersecretphrease")

/*CreateJwtToken packages JSON payload to JWT format*/
func CreateJwtToken() string {
	expires := time.Now().Add(time.Duration(10) * time.Second)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": expires.String(),
		"iat": time.Now().String(),
		"iss": "flow_server",
		"aud": "7eecb1a9-bf11-4134-8808-11eb94125031",
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(mySignedKey)

	fmt.Println(tokenString, err)
	return tokenString
}

func getClaim() jws.Claims {
	expires := time.Now().Add(time.Duration(10) * time.Second)

	claims := jws.Claims{}
	claims.SetExpiration(expires)
	claims.SetIssuedAt(time.Now())
	claims.SetIssuer("flow_server")
	claims.SetAudience("7eecb1a9-bf11-4134-8808-11eb94125031")
	return claims
}
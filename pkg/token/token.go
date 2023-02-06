package token

import (
	"cmd/main.go/internal/config"
	"cmd/main.go/pkg/logger"
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

var (
	privateKey     *rsa.PrivateKey // из файла  PrivateKeyFile
	publicKey      *rsa.PublicKey  // из файла  PublicKeyFile
	expiresMinutes time.Duration
)

// Init загрузка закрытого и открытого ключа из файлов
func Init(configJWT *config.JWTConfig) (err error) {
	expiresMinutes = configJWT.ExpiresMinutes

	// загрузка закрытого ключа из файла
	if configJWT.PrivateKeyFile != "" {
		err = LoadPrivateRSAKey(configJWT.PrivateKeyFile)
		if err != nil {
			return
		}
	}

	// загрузка публичного ключа из файла
	err = LoadPublicRSAKey(configJWT.PublicKeyFile)
	return
}

// LoadPrivateRSAKey загрузка закрытого ключа из файла
func LoadPrivateRSAKey(privateKeyFile string) error {
	//pathApp := filepath.Dir(os.Args[0]) + string(filepath.Separator)
	//privateKeyFile = pathApp + privateKeyFile

	// загрузка приватного ключа из файла
	keyByte, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return logger.Wrap(&err)
	}
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(keyByte)
	if err != nil {
		return logger.Wrap(&err)
	}

	log.Debug().Str("File", privateKeyFile).Msg("LoadPrivateRSAKey+")
	return nil
}

// LoadPublicRSAKey загрузка публичного ключа из файла
func LoadPublicRSAKey(publicKeyFile string) error {
	//pathApp := filepath.Dir(os.Args[0]) + string(filepath.Separator)
	//publicKeyFile = pathApp + publicKeyFile

	keyByte, err := os.ReadFile(publicKeyFile)
	if err != nil {
		return logger.Wrap(&err, "token.LoadPublicRSAKey")
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(keyByte)
	if err != nil {
		return logger.Wrap(&err, "token.LoadPublicRSAKey")
	}
	log.Debug().Str("publicKeyFile", publicKeyFile).Msg("token.LoadPublicRSAKey+")
	return nil
}

// Validate проверка токена
func Validate(token *string) (jwt.MapClaims, error) {
	// проверка в токене Bearer
	if len(*token) < 8 || (*token)[:7] != "Bearer " {
		return nil, errors.New(HeaderAuthorizationNotBeginBearer)
	}

	*token = (*token)[7:]

	tok, err := jwt.Parse(*token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected Method: %s", jwtToken.Header["alg"])
		}

		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf(InvalidToken)
	}

	return claims, nil
}

// Create создания токена
// Пример использования:
// tokenString, err := token.Create(map[string]interface{}{"email":user.Email, "role": user.Role})
func Create(mapDataJWT map[string]interface{}) (tokenString string, err error) {
	claims := make(jwt.MapClaims)
	for key, value := range mapDataJWT {
		claims[key] = value
	}
	claims["exp"] = time.Now().UTC().Add(time.Minute * expiresMinutes).Unix()
	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
}

// func VerifyToken(tokenString string) error {
//
//		// parse and verify signature
//		token, err := jwt.Parse([]byte(tokenString), jwt.WithVerify(jwa.HS256, []byte(config.Config.JWT.SecretJWTKey)))
//		if err != nil {
//			return err
//		}
//
//		// validate the essential claims
//		err = jwt.Validate(token)
//		return err
//	}

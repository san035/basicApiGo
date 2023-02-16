package token

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
	"github.com/san035/basicApiGo/pkg/logger"
	"os"
	"time"
)

var (
	privateKey     *rsa.PrivateKey // из файла  PrivateKeyFile
	publicKey      *rsa.PublicKey  // из файла  PublicKeyFile
	expiresMinutes time.Duration
)

// Init загрузка закрытого и открытого ключа из файлов
func Init(configJWT *JWTConfig) (err error) {
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
		return nil, logger.New(HeaderAuthorizationNotBeginBearer, map[string]string{"Authorization": CoveredToken(token)})
	}

	*token = (*token)[7:]

	tok, err := jwt.Parse(*token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected Method: %s", jwtToken.Header["alg"])
		}

		return publicKey, nil
	})
	if err != nil {
		return nil, logger.Wrap(&err, map[string]string{"Authorization": CoveredToken(token)})
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, logger.New(InvalidToken, map[string]string{"Authorization": CoveredToken(token)})
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

// Использование:
//
//	var Config = struct {
//		JWT token.JWTConfig
//	}{}
type JWTConfig struct {
	// секретный ключ для формирования токена, алгоритм rs256
	PrivateKeyFile string        `env:"JWT_FILE_PRIVATE_KEY_RSA" default:"rsa_key/jwt_privat_key_rsa" yaml:"PrivateKeyFile"`
	PublicKeyFile  string        `env:"JWT_FILE_PUBLIC_KEY_RSA"  default:"rsa_key/jwt_public_key_rsa" yaml:"PublicKeyFile"`
	ExpiresMinutes time.Duration `default:"10800" env:"JWT_EXPIRES_MINUTES" yaml:"ExpiresMinutes"` // срок действия токена в минутах
}

// CoveredToken возвращает для публичного показа
func CoveredToken(token *string) string {
	len_token := len(*token)
	if len_token > 12 {
		return (*token)[0:10] + "***" + (*token)[len_token-4:]
	}
	return "Empty"
}

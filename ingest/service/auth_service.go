package service

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"ingest/config"
	"ingest/model"
	"ingest/utils"
	"time"
)

func RegisterAccount(user model.User) (string, error) {
	hash, err := HashPassword(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = hash
	user.ID = uuid.NewString()
	token, err := GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", err
	}
	CreateMinimumViableUser(user)
	return token, nil
}

func LoginAccount(user model.User) (string, error) {
	hash := GetPasswordForEmail(user.Email)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(user.Password))
	if err != nil {
		utils.SugarLogger.Errorln(err.Error())
		return "", err
	}
	id := GetIDForEmail(user.Email)
	token, err := GenerateJWT(id, user.Email)
	if err != nil {
		return "", err
	}
	return token, nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.SugarLogger.Errorln(err.Error())
		return "", err
	}
	return string(hash), nil
}

func GenerateJWT(id string, email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &model.AuthClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.AuthSigningKey))
	if err != nil {
		utils.SugarLogger.Errorln(err.Error())
		return "", err
	}
	return signedToken, nil
}

func ValidateJWT(token string) (*model.AuthClaims, error) {
	claims := &model.AuthClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return config.AuthSigningKey, nil
	})
	if err != nil {
		utils.SugarLogger.Errorln(err.Error())
		return nil, err
	}
	return claims, nil
}

func CreateMinimumViableUser(user model.User) interface{} {
	user.FirstName = ""
	user.LastName = ""
	user.Subteam = ""
	return DB.Create(&user)
}

func GetIDForEmail(email string) string {
	var id string
	DB.Table("user").Where("email = ?", email).Select("id").Scan(&id)
	return id
}

func GetPasswordForEmail(email string) string {
	var password string
	DB.Table("user").Where("email = ?", email).Select("password").Scan(&password)
	return password
}

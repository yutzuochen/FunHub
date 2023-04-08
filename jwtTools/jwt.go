package jwtTools

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func ValidateJWT(secretKey string, next func(w http.ResponseWriter, r *http.Request, id string)) http.Handler {
	fmt.Println("start to ValidateJWT")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if r.Header["Token"] != nil {
		if cookie, err := r.Cookie("token"); err == nil {
			token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("not authorized"))
					// return SecretKey, nil
				}
				return []byte(secretKey), nil
			})
			// error in pasing token
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("not authorized: " + err.Error()))
				return
			}
			// fmt.Println("token.Claims: ", token.Claims)
			// fmt.Println("token.Claims['id']: ", token.Claims.(jwt.MapClaims)["id"])
			if token.Valid {
				// fmt.Println("token is valid!")
				// http.Redirect(w, r, "/web", http.StatusMovedPermanently)
				next(w, r, token.Claims.(jwt.MapClaims)["id"].(string))
				//....
			} else {
				fmt.Println("token is not valid :(((((")
				http.Redirect(w, r, "/login", http.StatusMovedPermanently)
			}
		} else {
			fmt.Println("valid error is: ", err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized, some err"))
			return
		}
	})
}

func CreateJWT(secretKey, id string) (string, error) {
	fmt.Println("start to Create JWT")

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Hour).Unix()
	claims["id"] = id

	tokenStr, err := token.SignedString([]byte(secretKey))
	fmt.Println("tokenStr: ", tokenStr)

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return tokenStr, nil
}

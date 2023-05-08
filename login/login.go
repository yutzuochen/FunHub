package login

import (
	"Funhub/db"
	"Funhub/jwtTools"
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

const SecretKey = "mason boxing golang"

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start to login ")

	// post or get?
	s := `
	<html>
	<head>
	<meta charset="utf-8">
	<title>Auth</title>
	</head>
	<body>
	
	<h1>Login</h1>
	<p>
		<form name="input" method="post">
		<table>
			<tr>
				<td>account：</td>
				<td><input type="text" name="account"></td>
			</tr>
			<tr>
				<td>keyword：</td>
				<td><input type="password" name="password"></td>
			</tr>
		</table>
			<input type="submit" value="Submit">
		</form>
	</p> 
	</body>
	</html>`
	if r.Method == http.MethodGet {
		fmt.Fprint(w, s)
		return
	}

	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		return
	}
	account := r.Form.Get("account")
	password := r.Form.Get("password")
	fmt.Println("account: ", account)
	fmt.Println("password: ", password)
	id := checkPW(account, password)
	fmt.Println("member id: ", id)
	// if id != ""  => means no that account
	if id != "" {
		// set cookie
		tokenString, err := jwtTools.CreateJWT(SecretKey, id)
		if err != nil {
			fmt.Println("create error: ", err)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: tokenString,
			// Expires: expirationTime,
		})
		fmt.Fprintln(w, "Success to login! Welcome to Funhub :))")
		// http.Redirect(w, r, "/web", http.StatusMovedPermanently)
	} else {
		fmt.Fprintln(w, s)
		fmt.Fprintln(w, "loginng failed, get fucking out!")
	}
}

func checkPW(account, password string) string {
	fmt.Println("===== start to checkPW =====")
	sqlStr := `select account, id from login where account = ? and password =?`

	rows, err := db.DbHdr.Query(sqlStr, account, password)
	if err != nil {
		fmt.Println("db query failed, error: ", err)
		return ""
	}
	defer rows.Close()
	if rows.Next() {
		var ac string
		var id string

		if err := rows.Scan(&ac, &id); err != nil {
			log.Fatal("gg: ", err)
		}
		log.Printf("ac %s id is %s\n", ac, id)

		return id
	}
	return ""
}

func Validate(secretKey string, next func(w http.ResponseWriter, r *http.Request, id string)) http.Handler {
	fmt.Println("start to ValidateJWT")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie("token"); err != nil {
			logrus.Errorln("error from validating: ", err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("You haven't be authorized, please login fitst"))
			return
		} else {
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

			if token.Valid {
				next(w, r, token.Claims.(jwt.MapClaims)["id"].(string))
			} else {
				fmt.Println("token is not valid :(((((")
				http.Redirect(w, r, "/login", http.StatusMovedPermanently)
			}
		}
	})
}

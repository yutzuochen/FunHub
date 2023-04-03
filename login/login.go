package login

import (
	"Funhub/db"
	"Funhub/jwtTools"
	"fmt"
	"log"
	"net/http"
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
		fmt.Println("w.Header(): ", w.Header())

		// http.Redirect(w, r, "/web", http.StatusMovedPermanently)
		fmt.Println("finish??")
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
	fmt.Println("[checkPW] rows: ", rows)
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

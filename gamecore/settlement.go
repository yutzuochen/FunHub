package gamecore

import (
	"Funhub/db"
	"fmt"
)

func diceSettlement(userID string, ante int, netIncome int) {
	fmt.Println("start to dice settlement")

	// UPDATE the field with the new value
	_, err := db.DbHdr.Exec("UPDATE user.account SET balance = balance + ? WHERE id = ?", netIncome, userID)
	if err != nil {
		fmt.Println("db.Exec error: ", err)
	}
}

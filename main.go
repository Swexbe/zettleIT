package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/swexbe/zettleIT/api"
)

const timeFormat = "2006-01-02"

func main() {

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	token := api.GetAuthkey(username, password)

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Payout Date: ")
	date, _ := reader.ReadString('\n')

	date = strings.TrimSuffix(date, "\n")

	endDate, err := time.Parse(timeFormat, date)
	if err != nil {
		log.Fatalln(err)
	}

	startDate := endDate.AddDate(0, 0, -14)

	sdString := startDate.Format(timeFormat)
	edString := endDate.Format(timeFormat)

	transactions := api.GetTransactions(sdString, edString, token)
	purchases := api.GetPurchases(sdString, edString, token)

	purchasesMap := make(map[string]string)

	for _, v := range purchases {
		purchasesMap[v.Payments[0].UUID] = v.UserDisplayName
	}

	amountSold := make(map[string]int)

	numPayouts := 0

	fmt.Print("\nDISTRIBUTION OF PAYMENTS: \n")

	for _, v := range transactions {

		if v.Type == "PAYOUT" {
			numPayouts++

		}

		if numPayouts == 0 {
			continue
		}

		if numPayouts >= 2 {
			break
		}

		seller := purchasesMap[v.UUID]
		if seller == "" {
			seller = "Total"
		}
		amountSold[seller] = amountSold[seller] + v.Amount

	}

	for key, value := range amountSold {

		amount := float64(value) / 100

		fmt.Printf("%s : %.2f kr \n", key, amount)
	}

}

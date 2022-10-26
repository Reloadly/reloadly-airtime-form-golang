package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	reloadly "github.com/reloadly/reloadly-sdk-golang/airtime"
)

type TopupDetails struct {
	OperatorID  string
	Amount      string
	Reference   string
	Number      string
	CountryCode string
}

func main() {
	tmpl, err := template.ParseFiles("./templates/airtime.html")
	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}

	rAccessToken, err := reloadly.NewClient("rzcLzmMmxZ919IVPgQlr6MDxiJRkEyjA", "7qpcVxQZ44-DDe1PBoZkM4b7WGnUBY-600Tkaol88NOrPn8yoeojMgQALtFGuDC", true)
	if err != nil {
		fmt.Println(err)
	}

	http.Handle("/css/", http.StripPrefix("/css", http.FileServer((http.Dir("css")))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		rTopupDetails := reloadly.Topuprequest{
			Amount:     r.FormValue("amount"),
			OperatorID: r.FormValue(("operatorid")),
		}

		rCustomIdentifier := reloadly.AddCustomIdentifier(r.FormValue("reference"))

		rPhone := reloadly.Phone{
			Number:      r.FormValue("number"),
			CountryCode: r.FormValue(("countrycode")),
		}
		rTopUp, err := rAccessToken.Topup(rTopupDetails.Amount, rTopupDetails.OperatorID, true, rPhone, rCustomIdentifier)
		if err != nil {
			fmt.Println(err)
		}

		jsonResp, err := json.Marshal(rTopUp)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}

		fmt.Printf("json data: %s\n", jsonResp)

		tmpl.Execute(w, struct{ Success bool }{true})

	})

	http.ListenAndServe(":8000", nil)

}

package main 

import (
	"google.golang.org/api/sheets/v4"
	"github.com/joho/godotenv"
	"context"
	"fmt"
	"log"
	"os"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"slices"
)

func main () {
	godotenv.Load()
	
	// data := []map[string]interface{}{
	// 	{"a":"val1", "b":"val2"},
	// 	{"a":"val1", "c":"val3"},
	// 	{"a":"val1", "c":"val3"},
	// }
	//clearGoogleSheet(os.Getenv("SHEET_ID"), "Name")
	//writeToGoogleSheets(os.Getenv("SHEET_ID"), "Name", data)

}

func clearGoogleSheet (spreadsheetId string, sheet string) {
	srv := getGoogleSheetsService()
	clearRequest := sheets.ClearValuesRequest{}
	_, err := srv.Spreadsheets.Values.Clear(spreadsheetId, sheet, &clearRequest).Do()
    if err != nil {
        fmt.Print(err)
    }
}

func getHeadingsFromData(data []map[string]interface{}) []string {
	var keys []string
	for i := 0; i < len(data); i++ {
		for k := range data[i] {
			if !slices.Contains(keys, k) {
				keys = append(keys, k)
			}
		}
	}
	return keys
}

func getGoogleSheetsService () (*sheets.Service) {
	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "YOUR_REDIRECT_URL",
		Scopes: []string{
			"https://www.googleapis.com/auth/spreadsheets",
		},
		Endpoint: google.Endpoint,
	}

	tok := &oauth2.Token{}
	tok.RefreshToken = os.Getenv("GOOGLE_REFRESH_TOKEN")
	client := config.Client(context.Background(), tok)

	srv, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
	return srv

}

func writeToGoogleSheets (spreadsheetId string, sheet string, data []map[string]interface{}) {
	var vr sheets.ValueRange
	var rows [][]interface{}

	var headings = getHeadingsFromData(data)
	var row []interface{}
	for i := range headings {
		row = append(row, headings[i])
	}
	rows = append(rows, row)

	for i := 0; i < len(data); i++ {
		var row []interface{}
		for _, key := range headings {
			row = append(row, data[i][key])
		}
		rows = append(rows, row)
	}

	vr.Values = rows
	srv := getGoogleSheetsService()
	srv.Spreadsheets.Values.Update(spreadsheetId, sheet, &vr).ValueInputOption("RAW").Do()
}

func readFromGoogleSheets (spreadsheetId string, sheet string) {
	srv := getGoogleSheetsService()
	resp, _ := srv.Spreadsheets.Values.Get(spreadsheetId, sheet).Do()

	if len(resp.Values) > 0 {
		for _, row := range resp.Values {
			fmt.Printf("%s\n", row) 
			//fmt.Printf("%s, %s, %s, %s\n", row[0], row[1], row[2], row[3]) // Print columns A and E, which correspond to indices 0 and 4.
		}	
	}
}


func appendDataToGoogleSheets (spreadsheetId string, sheet string) {
	//srv.Spreadsheets.Values.Append(spreadsheetId, sheet, &vr).ValueInputOption("RAW").Do()
}
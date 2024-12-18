package controller

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/joaovrmoraes/http-go-cli/model"
	"github.com/joaovrmoraes/http-go-cli/view"
)

func HandleRequest(method, url, bearer, data string, save bool) {
	model.LoadHistoryFromFile()

	start := time.Now()
	resp, err := model.MakeRequest(method, url, bearer, data)
	if err != nil {
		fmt.Printf("Error making the request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	elapsed := time.Since(start).Round(time.Millisecond)
	timeColor := color.New(color.FgCyan)
	timeColor.Printf("Time elapsed: %v \n", elapsed)

	view.PrintStatus(resp.StatusCode)

	view.PrintHeaders(resp.Header)
	body, err := model.ReadResponseBody(resp.Body)
	if err != nil {
		fmt.Printf("Error reading the response body: %v\n", err)
		return
	}

	if len(body) == 0 {
		fmt.Println("Response body is empty, nothing to display.")
		return
	}

	coloredJSON, err := view.FormatJSON(body)
	if err != nil {
		fmt.Printf("Error formatting the JSON: %v\n", err)
		return
	}

	title := method + " : " + fmt.Sprintf("%d", resp.StatusCode) + " | " + url + " | " + elapsed.String()

	if save {
		view.SaveToFile(coloredJSON)
	} else {
		view.StartInterface(string(coloredJSON), title, resp.Header)
	}

	model.AddToHistory(method, url, bearer, data)
}

package config_collector

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
)

type SubscriptionLoader struct {
}

func (s SubscriptionLoader) GetSub() {
	// URL of the webpage
	url := "https://raw.githubusercontent.com/soroushmirzaei/telegram-configs-collector/main/splitted/mixed"

	// Send HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching webpage:", err)
		return
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: HTTP status %d\n", resp.StatusCode)
		return
	}

	// Parse the HTML
	doc, err := html.Parse(resp.Body)

	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	text := doc.FirstChild.LastChild.FirstChild.Data

	file, err := os.Create("/home/ali/Desktop/v2ray_client/configs/output.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(text)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Webpage text saved to output.txt")
}

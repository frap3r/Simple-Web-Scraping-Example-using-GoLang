package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Scrape() {
	fmt.Print("Type product name to search (as clear as possible): ")
	reader := bufio.NewReader(os.Stdin)
	productName, _ := reader.ReadString('\n')               // Read user's input untill '\n'
	productName = strings.Replace(productName, "\n", "", 1) // Remove '\n' from input
	originalProductName := productName                      // Store the original input
	productName = strings.Replace(productName, " ", "+", 9) // Replace 'spaces' with "+" to put it in URL
	productName = strings.ToLower(productName)

	url := "http://www.ebay.com/sch/i.html?_from=R40&_trksid=p2380057.m570.l1313&_nkw=" + productName

	res, _ := http.Get(url)

	// Check if there's a connection error
	if res.StatusCode != 200 {
		fmt.Println("Connection Error : ", res.StatusCode)
		return
	}

	// Load html document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	s := "#mainContent .srp-results .s-item__title"  // Selector for tittle
	s2 := "#mainContent .srp-results .s-item__price" // Selector for price

	noResult := true               // Checker
	prices := make([]string, 0, 1) // Store prices in this slice
	doc.Find(s2).Each(func(i int, selection *goquery.Selection) {

		price := selection.Text()
		prices = append(prices, price)

	})

	fmt.Printf("\n--- Results ---\n")
	// Get titles from htlm document
	doc.Find(s).Each(func(i int, selection *goquery.Selection) {

		title := selection.Text()
		productNameNoSpace := strings.Replace(originalProductName, " ", "", 99)
		titleNoSpace := MakeComperable(title)

		// Avoid from irrelevent results (as much as possible)
		if strings.Contains(titleNoSpace, productNameNoSpace) { // Title must contains the exact user input
			noResult = false // Result found

			fmt.Printf("Product Title : %s   ==> Price : %s\n", title, prices[i])

		}

	})

	// Result NOT found
	if noResult {
		fmt.Printf("No results ...\nMake sure you type everything correctly !\n")
		return

	}

}

func MakeComperable(title string) string {
	// Preparing the title proper to compare
	title = strings.ToLower(title)
	title = strings.TrimSpace(title)
	title = strings.Replace(title, " ", "", 99)
	return title

}

func main() {
	Scrape()
}

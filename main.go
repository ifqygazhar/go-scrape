package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

const logo = `

 ██████╗  ██████╗ ███████╗ ██████╗██████╗  █████╗ ██████╗ ███████╗
██╔════╝ ██╔═══██╗██╔════╝██╔════╝██╔══██╗██╔══██╗██╔══██╗██╔════╝
██║  ███╗██║   ██║███████╗██║     ██████╔╝███████║██████╔╝█████╗  
██║   ██║██║   ██║╚════██║██║     ██╔══██╗██╔══██║██╔═══╝ ██╔══╝  
╚██████╔╝╚██████╔╝███████║╚██████╗██║  ██║██║  ██║██║     ███████╗
 ╚═════╝  ╚═════╝ ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝     ╚══════╝
`

func main() {
	fmt.Printf("\033[32m%s\033[0m\n", logo)

	domain := getInput("Domain: ")
	fromPage := getInput("From page: ")
	toPage := getInput("To page: ")

	domainExt(domain, fromPage, toPage)
}

func getInput(prompt string) string {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	return input
}

func domainExt(extension string, from, to string) {
	fromPage, err := strconv.Atoi(from)
	if err != nil {
		fmt.Println("Invalid 'From page' input.")
		os.Exit(1)
	}
	toPage, err := strconv.Atoi(to)
	if err != nil {
		fmt.Println("Invalid 'To page' input.")
		os.Exit(1)
	}

	fileName := "result-" + extension + ".txt"
	for page := fromPage; page <= toPage; page++ {
		url := fmt.Sprintf("https://zoxh.com/tld/%s/%d", extension, page)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error fetching page %d: %v\n", page, err)
			continue
		}
		defer resp.Body.Close()

		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading page %d: %v\n", page, err)
			continue
		}

		regex := regexp.MustCompile(`href="/i/([a-zA-Z0-9.-]+)"`)
		domains := regex.FindAllSubmatch(content, -1)

		for _, domain := range domains {
			fmt.Printf("[+] Page %d => \033[32m%s\033[0m\n", page, domain[1])
			f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println("Error opening file:", err)
				continue
			}
			_, err = f.WriteString(fmt.Sprintf("%s\n", domain[1]))
			if err != nil {
				fmt.Println("Error writing to file:", err)
			}
			f.Close()
		}
	}
}


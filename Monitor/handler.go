package Monitor

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

func IdFound(string2 string, page playwright.Page) bool {
	locator := page.Locator(string2)
	exists, err := locator.Count()
	if err != nil {
		log.Panicf("could not count element: at log in here %v", err)
	}
	if exists > 0 {
		return true

	} else {
		return false
	}

}

func TextFound(string2 string, page playwright.Page) bool {
	bodyText, err := page.Locator("body").InnerText()
	if err != nil {
		log.Panicf("could not get inner text: %v", err)
	}
	if strings.Contains(bodyText, "The event has ended") {
		return true
	} else {
		return false
	}

}

func CheckForTickets(string2 string, page playwright.Page) bool {
	if _, err := page.Goto(string2); err != nil {
		log.Panicf("could not goto: %v", err)
	}
	ClickPW(".btn.btn-large.btn-find-seats", page)

	if TextFound("No Tickets found", page) == false {
		//send webhook tickets are available
		return true
	} else {
		return false
	}
}

func ClickPW(string string, page playwright.Page) {
	locator := page.Locator(string)
	AssertErrorToNil("Failed to click on button", locator.Click())
	time.Sleep(800)
}

func ProxyLoad() []ProxyStruct {
	var returnPS []ProxyStruct
	var path = "./ProxyList.csv"
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("couldn't open - err: %v", err)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	for i := 0; true; i++ {
		if i == 0 {
			log.Println("Loading proxies")
			_, err := csvReader.Read()
			if err != nil {
				log.Fatalf("failed to open csv - err: %v", err)
			}

		} else {
			rec, err := csvReader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatalf("CSV reader failed - err : %v", err)
			}
			split := strings.Split(rec[0], ":")
			srv := (split[0] + ":" + split[1])
			usr := split[2]
			pss := split[3]

			var accDataStrct = ProxyStruct{
				ip:  srv,
				usr: usr,
				pw:  pss,
			}
			returnPS = append(returnPS, accDataStrct)

		}

	}
	return returnPS
}

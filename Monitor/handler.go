package Monitor

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

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
			fmt.Println("Loading proxies")
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

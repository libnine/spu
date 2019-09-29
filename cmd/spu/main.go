package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"../../src/cal"
	"../../src/fin"
	"../../src/pretty"
	"../../src/rh"
	"../../src/tr"
	"../../src/vol"
)

func main() {
	args := os.Args[1:]

	var wg sync.WaitGroup
	var syms, data, trend []string

	switch v := len(args); {
	case v == 1:
		switch arg := os.Args[1]; {
		case len(arg) > 5:
			fmt.Printf("%s: invalid ticker format.\n", os.Args[1])

		case arg == "help", arg == "-h":
			fmt.Printf("\nUsage:\tspu [OPTIONS] COMMAND\n\nA tool for getting securities data.\n\nOptions:")
			fmt.Printf("\n\t-c\tDaily economic calendar.")
			fmt.Printf("\n\t-g\tCNN Fear & Greed indices.")
			fmt.Printf("\n\t-rh\tRobinhood trends. Use the following commands:")
			fmt.Printf("\n\t\t\tdec\tTop decreases in ownership.")
			fmt.Printf("\n\t\t\tinc\tTop increases in ownership.")
			fmt.Printf("\n\t\t\tmost\tMost popular robinhood stocks by ownership.")
			fmt.Printf("\n\t-t\tStockTwits & Yahoo trending tickers.")
			fmt.Println("\n")

		case arg == "-c":
			cal.Calendar()

		case arg == "-g":
			chG := make(chan []string)
			go tr.Greed(&chG)
			_greed := <-chG
			fmt.Printf("\nNow:\t\t%s\nPrevious Day:\t%s\nPrevious Week:\t%s\nPrevious Month:\t%s\nPrevious Year:\t%s\n\n",
				_greed[0], _greed[1], _greed[2], _greed[3], _greed[4])

		case arg == "-t":
			fmt.Println(trend)
			fallthrough

		default:
			return
		}

	case v > 1:
		if args[0] == "-rh" {
			if args[1] == "most" {
				rh.RHpop()
				return
			}
			// } else if args[2] == "dec" {
			// 	RHdec()
			// 	return
			// } else if args[2] == "inc" {
			// 	RHinc()
			// 	return
			// }
		}

		for _, v := range args {
			if len(v) > 5 {
				fmt.Printf("%s: invalid ticker format.\n", v)
				continue
			}

			syms = append(syms, v)
		}

		if len(syms) != 0 {
			for _, v := range syms {
				if regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(v) == false {
					fmt.Printf("Passing %s for invalid format.\n", v)
					continue
				}

				wg.Add(1)
				go func(v string, wg *sync.WaitGroup) {
					ch := make(chan []string)
					chFv := make(chan []string)

					go vol.Vol(v, &ch)
					go fin.Finviz(v, &chFv)

					_vol := <-ch
					_finviz := <-chFv

					dump := fmt.Sprintf("%s\t%s\t\t%s\t%s\t\t%s\t\t%s\t\t%s", strings.ToUpper(v), _finviz[0], _finviz[1], _vol[1], _vol[0], _vol[2], _vol[3])

					data = append(data, dump)
					wg.Done()
				}(v, &wg)
			}
		}

		wg.Wait()
		pretty.Pretty(data)
	}
}

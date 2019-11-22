package cli

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"../cal"
	"../etc"
	"../fin"
	"../pretty"
	"../quandl"
	"../rh"
	"../tr"
	"../vol"
)

func Run(args []string) {
	var wg sync.WaitGroup
	var syms, data []string

	switch v := len(args); {
	case v == 1:
		switch arg := os.Args[1]; {
		case len(arg) > 5:
			fmt.Printf("%s: invalid ticker format.\n", os.Args[1])

		case arg == "help", arg == "-h":
			fmt.Printf("\nUsage:\tspu [OPTIONS] COMMAND [ARG...]\n\nA tool for retrieving securities data.\n\nOptions:")
			fmt.Printf("\n\t-a\tAAII Investor Sentiment.")
			fmt.Printf("\n\t-c\tDaily economic calendar.")
			fmt.Printf("\n\t-g\tCNN Fear & Greed indices.")
			fmt.Printf("\n\t-rh\tRobinhood trends. Use the following commands:")
			fmt.Printf("\n\t\t\tdec\tTop decreases in ownership.")
			fmt.Printf("\n\t\t\tinc\tTop increases in ownership.")
			fmt.Printf("\n\t\t\tmost\tMost popular robinhood stocks by ownership.")
			fmt.Printf("\n\t\t\tpop\tLargest robinhood popularity changes.")
			fmt.Printf("\n\t-t\tStockTwits & Yahoo trending tickers.")
			fmt.Printf("\n\nCommands:")
			fmt.Printf("\n\tnews\tPrint recent headlines from news sources.")
			fmt.Printf("\n\t\t-bbg\tBloomberg\n\t\t-mw\tMarketWatch\n\t\t-rtrs\tReuters\n\t\t-wsj\tWall Street Journal")
			fmt.Println("\n")

		case arg == "-a":
			quandl.Q()

		case arg == "-c":
			cal.Calendar()

		case arg == "-g":
			chG := make(chan []string)
			go tr.Greed(&chG)
			_greed := <-chG
			fmt.Printf("\nNow:\t\t%s\nPrevious Day:\t%s\nPrevious Week:\t%s\nPrevious Month:\t%s\nPrevious Year:\t%s\n\n",
				_greed[0], _greed[1], _greed[2], _greed[3], _greed[4])

		case arg == "-t":
			tr.Trending()

		default:
			return
		}

	case v > 1:
		if args[0] == "news" {
			switch args[1] {
			case "-bbg":
				etc.News("bbg")
				return
			case "-mw":
				etc.News("mw")
				return
			case "-rtrs":
				etc.News("rtrs")
				return
			case "-wsj":
				etc.News("wsj")
				return
			default:
				fmt.Println("No valid news source selected.")
				return
			}
		}

		if args[0] == "-rh" {
			switch args[1] {
			case "most":
				rh.RHpop()
				return
			case "pop":
				rh.RHchg()
				return
			case "inc":
				rh.RHinc()
				return
			case "dec":
				rh.RHdec()
				return
			default:
				panic(fmt.Sprintf("No command for '%s' under -rh.\n", args[1]))
			}
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
		} else {
			fmt.Println("No command specified.")
		}

		wg.Wait()
		pretty.Pretty(data)
	}
}

package cli

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
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
			fmt.Printf("\n\t-a\tAAII Investor Sentiment")
			fmt.Printf("\n\t-g\tCNN Fear & Greed indices")
			fmt.Printf("\n\t-l\tUpcoming IPO lockup expirations")
			fmt.Printf("\n\t-t\tStockTwits & Yahoo trending tickers")
			fmt.Printf("\n\t-y\tYield curve")
			fmt.Printf("\n\nCommands:")
			fmt.Printf("\n\tcal\tDaily economic calendar")
			fmt.Printf("\n\t\tChoose weekday by number ('cal 1' returns Monday)\n")
			fmt.Printf("\n\tcef\tClosed end funds")
			fmt.Printf("\n\t\t-c\tCompile CEF data dump\n\t\t-d\tFunds down over 1%%\n\t\t-up\tFunds up over 1%%\n")
			fmt.Printf("\n\tnews\tRecent headlines from news sources")
			fmt.Printf("\n\t\t-bbg\tBloomberg\n\t\t-mw\tMarketWatch\n\t\t-rtrs\tReuters\n\t\t-wsj\tWall Street Journal\n")
			fmt.Printf("\n\trh\tRobinhood data")
			fmt.Printf("\n\t\t-dec\tTop decreases in ownership")
			fmt.Printf("\n\t\t-inc\tTop increases in ownership")
			fmt.Printf("\n\t\t-most\tMost popular robinhood stocks by ownership")
			fmt.Printf("\n\t\t-pop\tLargest robinhood popularity changes\n")
			fmt.Printf("\n\tpff\tPreferred stock data")
			fmt.Printf("\n\t\t-c\tCompile preferred data dump\n\t\t-d\tTickers down over 1%%\n\t\t-rel\tTickers with over 2x average volume\n\t\t-up\tTickers up over 1%%\n\t\t-y\tTickers with highest yield")
			fmt.Println("\n")

		case arg == "-a":
			quandl.Qaaii()

		case arg == "-g":
			chG := make(chan []string)
			go tr.Greed(&chG)
			_greed := <-chG
			fmt.Printf("\nNow:\t\t%s\nPrevious Day:\t%s\nPrevious Week:\t%s\nPrevious Month:\t%s\nPrevious Year:\t%s\n\n",
				_greed[0], _greed[1], _greed[2], _greed[3], _greed[4])

		case arg == "-l":
			fin.LUexp()

		case arg == "-t":
			tr.Trending()

		case arg == "-y":
			quandl.Qyc()

		default:
			return
		}

	case v > 1:
		if args[0] == "cal" {
			n, err := strconv.Atoi(args[1])
			if err != nil {
				log.Fatal(err)
			}

			if n == 6 || n == 7 {
				panic("No weekend data.")
			}

			if n == 0 || n >= 7 {
				panic(fmt.Sprintf("Invalid day number '%v' (Monday = 1, Sunday = 7)", args[1]))
			}

			cal.Calendar(n)
			return
		}

		if args[0] == "cef" {
			switch args[1] {
			case "-c":
				fin.CefCompile()
				return
			case "-d", "-down":
				fin.Cef("down")
				return
			case "-u", "-up":
				fin.Cef("up")
				return
			default:
				panic(fmt.Sprintf("No command for '%s' under cef.", args[1]))
			}
		}

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
				panic(fmt.Sprintf("No news source found for '%s'.", args[0]))
			}
		}

		if args[0] == "pff" {
			switch args[1] {
			case "-c":
				fin.PffCompile()
				return
			case "-d", "-down":
				fin.Pff("down")
				return
			case "-rel":
				fin.Pff("rel")
				return
			case "-u", "-up":
				fin.Pff("up")
				return
			case "-y":
				fin.Pff("yield")
				return
			default:
				panic(fmt.Sprintf("No command for '%s' under pff.", args[1]))
			}
		}

		if args[0] == "rh" {
			switch args[1] {
			case "-most":
				rh.RHpop()
				return
			case "-pop":
				rh.RHchg()
				return
			case "-inc":
				rh.RHinc()
				return
			case "-dec":
				rh.RHdec()
				return
			default:
				panic(fmt.Sprintf("No command for '%s' under rh.\n", args[1]))
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
			panic("No command specified.")
		}

		wg.Wait()
		pretty.Pretty(data)
	}
}

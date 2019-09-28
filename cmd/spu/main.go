package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"../../src/fin"
	"../../src/pretty"
	"../../src/vol"
	// "../../src/tr"
)

func main() {
	args := os.Args[1:]

	var wg sync.WaitGroup
	var syms, data, trend []string

	switch v := len(args); {
	case v == 1:
		fmt.Println("YO DOGS")
		os.Exit(1)

	case v > 1:
		for _, v := range args {
			if len(v) > 5 {
				fmt.Printf("%s: invalid ticker format.\n", v)
				continue
			}

			if v == "-t" {
				trend = append(trend, "yo")
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
					ch_ := make(chan []string)

					go vol.Vol(v, &ch)
					go fin.Finviz(v, &ch_)

					_vol := <-ch
					_finviz := <-ch_

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

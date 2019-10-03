package pretty

import (
	"fmt"
	"sort"
)

func run() {

}

func Pretty(dump []string) {
	sort.Strings(dump)
	fmt.Println("\n\tChange %\tEarnings\tShort Interest\t7 Day Implied\tWeekly Implied\tWeekly Expiration")

	for _, v := range dump {
		fmt.Println(v)
	}

	fmt.Println("")
}

// generates random permutations
// using multithreading (coroutine)
package main

import (
	"fmt"
	"github.com/pkg/profile"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
)

const MAXSIZE = 1024

// Returns an array containing 1..n shuffled in a random manner
func shuffleFIX(n int) (perm [MAXSIZE]int) {
	for i := 0; i < n; i++ {
		perm[i] = i
	}

	for i := n - 1; i > 1; i-- {
		p := rand.Intn(i)
		perm[i], perm[p] = perm[p], perm[i]
	}
	return perm
}

func fill_seq(x []int) {
	for i := 0; i < len(x); i++ {
		x[i] = i
	}
}

// Randomly shuffle an array in place
func shuffle(perm []int) {
	for i := len(perm) - 1; i > 1; i-- {
		p := rand.Intn(i)
		perm[i], perm[p] = perm[p], perm[i]
	}
}

// cocoutine creating a radnom permutation
var c = make(chan []int)

func shuff(n int) {
	var perm = make([]int, n)
	fill_seq(perm)
	shuffle(perm)
	c <- perm
}

func print_perm(perm []int) {
	for i := 0; i < len(perm); i++ {
		fmt.Print(string(65 + perm[i]))
	}
	fmt.Println()
}

func usage() {
	fmt.Println("Usage: ", os.Args[0], "permsize npermutations")
}

func main() {

	defer profile.Start(profile.ProfilePath(".")).Stop()
	var permsize int
	var npermutations int
	var ncpus int
	var err error

	if len(os.Args) < 3 {
		usage()
		os.Exit(1)
	}

	permsize, err = strconv.Atoi(os.Args[1])
	if err != nil {
		usage()
		os.Exit(1)
	}

	npermutations, err = strconv.Atoi(os.Args[2])
	if err != nil {
		usage()
		os.Exit(1)
	}

	ncpus = runtime.NumCPU()
	fmt.Println("ncpus=", ncpus)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < npermutations; i++ {
		go shuff(permsize)
	}

	for p := 0; p < npermutations; p++ {
		perm := <-c
		if perm[0] == perm[1] {
			fmt.Println("bug")
		}
	}

	perm := make([]int, permsize)
	fill_seq(perm)
	for p := 0; p < npermutations; p++ {
		shuffle(perm)
		if perm[0] == perm[1] {
			fmt.Println("bug")
		}

		//for  i := 0; i < permsize; i++ {
		// 	fmt.Println(perm[i])
		// }
	}

	for p := 0; p < npermutations; p++ {
		perm := shuffleFIX(permsize)
		if perm[0] == perm[1] {
			fmt.Println("bug")
		}

		//for  i := 0; i < permsize; i++ {
		// 	fmt.Println(perm[i])
		// }
	}
}

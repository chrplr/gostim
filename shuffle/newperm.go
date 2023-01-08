// generating random numbers uniformatly on [0, 1[
package main

import (
	"fmt"
	"flag"
	"pgregory.net/rand"
)


func RandomU(n int, seed uint64) []float32 {
	x := make([]float32, n)
	r := rand.New(seed)	
	for i:=0; i<n; i++ {
		x[i] = r.Float32()
	}
	return x
}


func RandomIntn(n int, max int, seed uint64) []int {
	x := make([]int, n)
	r := rand.New(seed)	
	for i:=0; i<n; i++ {
		x[i] = r.Intn(max)
	}
	return x
}



func main() {
	var count int
	var max int
	var seed uint64
	
	flag.IntVar(&count, "n", 1, "sample size")
	flag.Uint64Var(&seed, "s", 0, "seed")
	flag.IntVar(&max, "m", 10, "max")

	flag.Parse()

	vectorFloats := RandomU(count, seed)
	fmt.Println(vectorFloats)

	vectorIntn := RandomIntn(count, max, seed)
	fmt.Println(vectorIntn)
	
}

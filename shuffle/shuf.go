// https://stackoverflow.com/questions/30226438/generate-all-permutations-in-go
// by Paul Hankin
// This code iterates over all permutations without generating them all first.
// The slice p keeps the intermediate state as offsets in a Fisher-Yates shuffle algorithm.
// This has the property that the zero value for p describes the identity permutation.

package main

import "fmt"

func nextPerm(p []int) {
    for i := len(p) - 1; i >= 0; i-- {
        if i == 0 || p[i] < len(p)-i-1 {
            p[i]++
            return
        }
        p[i] = 0
    }
}

func getPerm(orig, p []int) []int {
    result := append([]int{}, orig...)
    for i, v := range p {
        result[i], result[i+v] = result[i+v], result[i]
    }
    return result
}

func main() {
    orig := []int{1, 2, 3, 4, 5}
    for p := make([]int, len(orig)); p[0] < len(p); nextPerm(p) {
        fmt.Println(getPerm(orig, p))
    }
}

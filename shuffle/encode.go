package main

import "fmt"

type Code struct {
	s2i map[string]int
	i2s map[int]string
	n int
}

func (a *Code) init() {
	a.s2i = make(map[string]int)
	a.i2s = make(map[int]string)
	a.n = 0
}

func (a *Code) add(token string) int {
	if p, ok := a.s2i[token]; ok {
		return p
	} else {
		a.n++
		a.s2i[token] = a.n
		a.i2s[a.n] = token
		return a.n
	}
}

func (a *Code) get(code int) string {
	return a.i2s[code]
}


func (a *Code) EncodeText(text []string) []int {
	var encoded_text []int
	encoded_text = make([]int, len(text))
	for i, v := range text {
		encoded_text[i] = a.add(v)
	}
	return encoded_text
}


func main() {
	var a Code
	a.init()

	fmt.Println(a.EncodeText([]string{"bonjour", "hello", "bonjour"}))
	fmt.Println(a)
	for _, v := range [...]int{1,1,2,2,1} {
		fmt.Print(a.get(v), " ")
	}
 
}

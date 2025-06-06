package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"strconv"
)

// https://dev.to/jpoly1219/stacks-in-go-54k
type Stack struct {
	items []string
}

func main() {
	port := "4000"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/evaluate", evaluate)

	log.Println("Starting server on : ", port)
	err := http.ListenAndServe(":" + port, mux)
	log.Fatal(err)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Ello! From home!"))
}

// p = +
// s = -
// m = *
// d = /
func evaluate(w http.ResponseWriter, r *http.Request) {
	operators := Stack{}
	operands := Stack{}
	expr := r.URL.Query().Get("expr")
	for i := 0; i < len(expr); i++ {
		token := string(expr[i])
		if strings.Compare("p", token) == 0 {
			operators.Push("+")
		} else if strings.Compare("s", token) == 0 {
			operators.Push("-")
		} else if strings.Compare("m", token) == 0 {
			operators.Push("*");
		} else if strings.Compare(")", token) == 0 {
			op := operators.Pop()
			value, _ := strconv.Atoi(operands.Pop())
			temp, _ := strconv.Atoi(operands.Pop()) 
			if strings.Compare("+", op) == 0 {
				value = temp + value
			} else if strings.Compare("-", op) == 0 {
				value = temp - value
			} else if strings.Compare("*", op) == 0 {
				value = temp * value
			}
			operands.Push(strconv.Itoa(value))
		} else if strings.Compare("(", token) != 0 {
			operands.Push(token)
		}
	}
	// https://golang.cafe/blog/golang-int-to-string-conversion-example.html
	w.Write([]byte("Evaluation = " + operands.Pop()))
}

func (s *Stack) Push(data string) {
	s.items = append(s.items, data)
}

func (s *Stack) Pop() string {
	if s.IsEmpty() {
		return ""
	}
	top := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return top
}

func (s *Stack) IsEmpty() bool {
	if len(s.items) == 0 {
		return true
	}
	return false
}

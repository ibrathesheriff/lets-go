# NOTES

## Comparing strings
The `strings.Compare()` function is a built-in function in Golang that compares two strings lexicographically. It returns an integer value that indicates the relationship between the two strings.

Where `str1` and `str2` are the two strings to be compared:
+ If the two strings are equal, the function returns 0.
+ If the first string is lexicographically less than the second string, the function returns a negative value.
+ If the first string is lexicographically greater than the second string, the function returns a positive value.

```golang
package main

import (
  "fmt"
  "strings"
)

func main() {
  greater := strings.Compare("Z", "A")
  equal := strings.Compare("W", "W")
  lesser := strings.Compare("A", "Z")

  fmt.Println(greater)
  fmt.Println(equal)
  fmt.Println(lesser)
}
```

Output:

```shell
1
0
-1
```

## How to Extract URL Query Parameters in Go
Based on: https://freshman.tech/snippets/go/extract-url-query-params/

```golang
func evaluate(w http.ResponseWriter, r *http.Request) {
    // e.g. http://localhost:4000/evaluate?expr=(2p2)
    // expr = (2p2)
	expr := r.URL.Query().Get("expr")
	w.Write([]byte("Evaluate! " + expr))
}
```

To loop through all query parameters in a request:
```golang
values := r.URL.Query()
for k, v := range values {
	fmt.Println(k, " => ", v)
}
```
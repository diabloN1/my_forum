package main

import "fmt"

func main() {
    a := []int{1, 2, 3, 4, 5}
    b := make([]int, len(a)) // Destination slice with length 3

    n := copy(b, a) // Copies first 3 elements from src to dst

    fmt.Println("Copied elements:", n)
    fmt.Println("Destination slice:", b)
}

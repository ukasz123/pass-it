package main

import (
	"fmt"
	"pass-it/cache"
)

func main() {
	var cache = cache.NewCache[any]()
	cache.Set("MyKey", "MyValue")
	cache.Set("MyNumber", "MyNumber")

	if val := cache.Get("MyKey"); val != nil {
		fmt.Printf("Got it: %s, %s\n", val, val.Value())
	} else {
		fmt.Println("Not found")
	}

	if val := cache.Get("MyKey"); val != nil {
		fmt.Printf("Got it: %s, %s\n", val, val.Value())
	} else {
		fmt.Println("Not found")
	}

}

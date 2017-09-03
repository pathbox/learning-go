package main

import (
	"context"
	"fmt"
)

func main() {
	type favContextKey string

	f := func(ctx context.Context, k favContextKey) {
		if v := ctx.Value(k); v != nil {
			fmt.Println("found value: ", v)
			return
		}
		fmt.Println("key not found: ", k)
	}

	k := favContextKey("language")
	fmt.Println(k)
	ctx := context.WithValue(context.Background(), k, "Nice") // {"language": "Nice"} 存储到了ctx中。 而并没有key 是color的key。context此时存的是key=>value的值
	f(ctx, k)
	f(ctx, favContextKey("color"))
}

// https://deepzz.com/post/golang-context-package-notes.html

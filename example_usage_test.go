package pathmatch_test

import (
	"fmt"

	"github.com/abichinger/pathmatch"
)

func Example_usage() {
	p, err := pathmatch.Compile("/api/:version/*")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(p.Match("/foo/bar"))
	fmt.Println(p.Match("/api/v1/"))
	fmt.Println(p.FindSubmatch("/api/v1/foo"))
	fmt.Println(p.FindSubmatch("/api/v2/foo/bar"))

	// Output: false
	// true
	// map[$0:foo version:v1]
	// map[$0:foo/bar version:v2]
}

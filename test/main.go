package main

import (
	// "encoding/json"
	"fmt"
	"strings"

	"github.com/boon4681/smanchai/vm"
)

func main() {
	lexer := smanchai.NewLexer(strings.NewReader("@user.role.name + \"I\" == \"HII\""))
	parser := smanchai.NewParser(lexer)
	ast := parser.Parse()
	vm := smanchai.Compile(ast)
	vm.AddStatic("user", func() *smanchai.Data {
		data, err := smanchai.Reflect(
			struct {
				role struct {
					name string
				}
			}{
				role: struct{ name string }{
					name: "HI",
				},
			},
		)
		if err != nil {
			panic(err)
		}
		return data
	})
	result, err := vm.Run()
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("%s\n", result)
	}
}




// func main() {
// 	// lexer := smanchai.NewLexer(strings.NewReader("@user.role.name == \"HII\""))
// 	// lexer := smanchai.NewLexer(strings.NewReader("@user.role.name + \"I\" == \"HII\" and 1.5- 1 == 2"))
// 	lexer := smanchai.NewLexer(strings.NewReader("2.5 - 10 ** 2"))
// 	// lexer.Dg = true
// 	parser := smanchai.NewParser(lexer)
// 	ast := parser.Parse()
// 	// {
// 	// 	b, err := json.Marshal(ast)
// 	// 	if err != nil {
// 	// 		fmt.Println(err)
// 	// 		return
// 	// 	}
// 	// 	fmt.Println(string(b))
// 	// }
// 	vm := smanchai.Compile(ast)
// 	vm.AddStatic("user", func() *smanchai.Data {
// 		data, err := smanchai.Reflect(
// 			struct {
// 				role struct {
// 					name string
// 				}
// 			}{
// 				role: struct{ name string }{
// 					name: "HI",
// 				},
// 			},
// 		)
// 		if err != nil {
// 			panic(err)
// 		}
// 		// {
// 		// 	b, err := json.Marshal(data)
// 		// 	if err != nil {
// 		// 		fmt.Println(err)
// 		// 	}
// 		// 	fmt.Println(string(b))
// 		// }
// 		return data
// 	})
// 	// {
// 	// 	b, err := json.Marshal(vm)
// 	// 	if err != nil {
// 	// 		fmt.Println(err)
// 	// 		return
// 	// 	}
// 	// 	fmt.Println(string(b))
// 	// }
// 	result, err := vm.Run()
// 	if err != nil {
// 		panic(err)
// 	} else {
// 		fmt.Printf("%s\n", result)
// 	}
// 	// smanchai.TestVM().Run()
// }

package java_test

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/rwxrob/java"
)

//go:embed testdata/hello.java
var helloJava string

// fooJava from go:embed testdata/hello.java

func ExampleRunJava() {

	err := java.RunJava("testdata/hello.java")
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// Hello, World!
}

func ExampleRunString() {

	raw := `
class HelloWorld {
    public static void main(String[] args) {
        System.out.println("Hello, World!");
    }
}
`

	err := java.RunString(raw)
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// Hello, World!
}

func ExampleRunClass() {

	os.Setenv("CLASSPATH", "testdata")
	err := java.RunClass("HelloWorld")
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// Hello, World!
}

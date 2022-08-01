package java_test

import (
	"embed"
	_ "embed"
	"fmt"
	"os"

	"github.com/rwxrob/fs/file"
	"github.com/rwxrob/java"
)

//go:embed testdata/javafiles/hello.java
var helloJava string

//go:embed testdata/javafiles
var javafiles embed.FS

func ExampleExtract() {

	java.CacheDir = "testdata/tmpcache"
	defer os.RemoveAll("testdata/tmpcache")

	if err := java.Extract(javafiles, "testdata/javafiles"); err != nil {
		fmt.Println(err)
	}

	fmt.Println(java.Cached("hello.java"))
	fmt.Println(java.Cached("HelloWorld.class"))
	fmt.Println(file.Exists("testdata/tmpcache/hello.java"))
	fmt.Println(file.Exists("testdata/tmpcache/HelloWorld.class"))

	// Output:
	// testdata/tmpcache/hello.java
	// testdata/tmpcache/HelloWorld.class
	// true
	// true

}

func ExampleRunJava() {

	err := java.RunJava("testdata/javafiles/hello.java")
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

func ExampleRunClass_nocache() {

	defer os.Setenv("CLASSPATH", os.Getenv("CLASSPATH"))
	os.Setenv("CLASSPATH", "testdata/javafiles")

	err := java.RunClass("HelloWorld")
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// Hello, World!
}

func ExampleRunClass_cached() {

	java.CacheDir = "testdata/tmpcache"
	defer os.RemoveAll("testdata/tmpcache")
	if err := java.Extract(javafiles, "testdata/javafiles"); err != nil {
		fmt.Println(err)
	}

	err := java.RunClass("HelloWorld")
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// Hello, World!
}

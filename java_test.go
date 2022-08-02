package java_test

import (
	"embed"
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/rwxrob/fs/file"
	"github.com/rwxrob/java"
)

//go:embed testdata/javafiles/hello.java
var helloJava string

//go:embed testdata/javafiles
var javafiles embed.FS

func ExampleClass2Path() {

	fmt.Println(java.Class2Path("foo.bar.Some"))
	fmt.Println(java.Class2Path("foo.bar.Some.class"))

	// Output:
	// foo/bar/Some.class
	// foo/bar/Some.class
}

func ExampleParseCmd() {

	c := `-Dfoo=bar HelloClass some args here`
	parsed := java.ParseCmd(strings.Fields(c)...)

	fmt.Println(parsed.Name)
	fmt.Println(parsed.Options)
	fmt.Println(parsed.Args)

	// Output:
	// HelloClass
	// [-Dfoo=bar]
	// [some args here]
}

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

func ExampleExec_java() {

	err := java.Exec("testdata/javafiles/hello.java")
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// Hello, World!
}

func ExampleExec_jar() {

	err := java.Exec("-jar", "testdata/files.jar")
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// Hello, World!
}

func ExampleExec_class_nocache() {

	defer os.Setenv("CLASSPATH", os.Getenv("CLASSPATH"))
	os.Setenv("CLASSPATH", "testdata/javafiles")

	err := java.Exec("HelloWorld")
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// Hello, World!
}

func ExampleExec_class_Cached() {

	java.CacheDir = "testdata/tmpcache"
	defer os.RemoveAll("testdata/tmpcache")
	if err := java.Extract(javafiles, "testdata/javafiles"); err != nil {
		fmt.Println(err)
	}

	err := java.Exec("HelloWorld")
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// Hello, World!
}

func ExampleOut_java_with_Args() {

	out := java.Out("-Dfoo=bar", "testdata/javafiles/fooprop.java")
	fmt.Println(out)

	// Output:
	// bar
}

func ExampleOut_class() {

	defer os.Setenv("CLASSPATH", os.Getenv("CLASSPATH"))
	os.Setenv("CLASSPATH", "testdata/javafiles")
	fmt.Println(java.Out("HelloWorld"))

	// Output:
	// Hello, World!
}

func ExampleOut_jar() {

	fmt.Println(java.Out("-jar", "testdata/files.jar"))

	// Output:
	// Hello, World!
}

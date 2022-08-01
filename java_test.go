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

func ExampleExecJava() {

	err := java.ExecJava("testdata/javafiles/hello.java")
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// Hello, World!
}

func ExampleExecString() {

	raw := `
class HelloWorld {
    public static void main(String[] args) {
        System.out.println("Hello, World!");
    }
}
`

	err := java.ExecString(raw)
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// Hello, World!
}

func ExampleExecJar() {

	err := java.ExecJar("testdata/files.jar")
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// Hello, World!
}

func ExampleExecClass_nocache() {

	defer os.Setenv("CLASSPATH", os.Getenv("CLASSPATH"))
	os.Setenv("CLASSPATH", "testdata/javafiles")

	err := java.ExecClass("HelloWorld")
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// Hello, World!
}

func ExampleExecClass_cached() {

	java.CacheDir = "testdata/tmpcache"
	defer os.RemoveAll("testdata/tmpcache")
	if err := java.Extract(javafiles, "testdata/javafiles"); err != nil {
		fmt.Println(err)
	}

	err := java.ExecClass("HelloWorld")
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// Hello, World!
}

func ExampleOutJava_with_Args() {

	out := java.OutJava("testdata/javafiles/fooprop.java", "-Dfoo=bar")
	fmt.Println(out)

	// Output:
	// bar
}

func ExampleOutString_with_Args() {

	raw := `
import java.util.*;
class HelloWorld {
    public static void main(String[] args) {
        System.out.println("Hello, World!");
				Properties p = System.getProperties();
				String value = (String)p.get("foo");
				System.out.println(value);
    }
}
`

	out := java.OutString(raw, "-Dfoo=bar")
	fmt.Println(out)

	// Output:
	// Hello, World!
	// bar
}

func ExampleOutClass() {

	defer os.Setenv("CLASSPATH", os.Getenv("CLASSPATH"))
	os.Setenv("CLASSPATH", "testdata/javafiles")
	fmt.Println(java.OutClass("HelloWorld"))

	// Output:
	// Hello, World!
}

func ExampleOutJar() {

	fmt.Println(java.OutJar("testdata/files.jar"))

	// Output:
	// Hello, World!
}

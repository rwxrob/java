/*

Package java uses Go embed.FS so that packages can be created to
encapsulate Java JAR, class, and raw source files that have been
embedded into the package with the default java executable on the host
system. This package includes a caching mechanism implemented in the
Extract and Cached functions so that the files need not be extracted
with every execution.  The java command invocation depends entirely on
the version of java installed on the host system and observes CLASSPATH
and other java-specific environment variables.

Options beginning with dash passed as arguments before the main
class/jar/java file are preserved. Options must use the equals or colon
format to avoid confusion with the main identifier. Arguments following
the class/jar/java argument are passed as expected.  No shell expansion
is performed.

The Exec function maps the output of the java command to the system
stdin/out/err (which can be redirected to a file by assigning to
os.Stdin, etc.) while the Out function returns a string with stdout and
logs stderr (see internal/exec.go).

*/
package java

import (
	"embed"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwxrob/fs"
	"github.com/rwxrob/java/internal"
)

// Cmd is a java command line with options preceding the named
// class/jar/java file. Args come after.
type Cmd struct {
	Name    string
	Options []string
	Args    []string
}

// CacheDir is set to os.UserCacheDir() plus "gojavacache" by default at
// init time.
var CacheDir string

func init() {
	dir, err := os.UserCacheDir()
	if err == nil {
		CacheDir = filepath.Join(dir, "gojavacache")
	}
}

// careful not to call more than once since will duplicate
func updateCP() {
	if os.Getenv("CLASSPATH") == "" {
		os.Setenv("CLASSPATH", CacheDir)
		return
	}
	os.Setenv("CLASSPATH",
		CacheDir+string(os.PathListSeparator)+os.Getenv("CLASSPATH"))
}

// Extract explicitly extracts all of an embedded file system into the
// CacheDir starting from the root path passed. Files in the CacheDir
// always have priority over anything else on the system since CacheDir
// is added to the beginning of the CLASSPATH.
func Extract(fsys embed.FS, root string) error {
	os.MkdirAll(CacheDir, fs.ExtractDirPerms)
	if err := fs.ExtractEmbed(fsys, root, CacheDir); err != nil {
		return err
	}
	updateCP()
	return nil
}

// Cached returns the full path the extracted cache location of the file
// indicated by it.
func Cached(file string) string {
	path := filepath.Join(CacheDir, file)
	if fs.Exists(path) {
		return path
	}
	return ""
}

// ParseCmd parses a typical java command line with options beginning
// with dash (and containing no spaces). The first non-dashed argument
// is considered the Name, or main class/java/jar file (see Cmd). The
// remaining arguments are stored as arguments to the class/java/jar
// itself.
func ParseCmd(cmd ...string) *Cmd {
	c := new(Cmd)

	for _, it := range cmd {
		if !strings.HasPrefix(it, "-") {
			if c.Name == "" {
				c.Name = it
				continue
			}
		}
		if c.Name == "" {
			c.Options = append(c.Options, it)
		} else {
			c.Args = append(c.Args, it)
		}
	}

	return c
}

// Class2Path translates a simple string into a class name adding the
// ".class" suffix if needed and replacing the dots (.) with the
// os.PathSeparator.
func Class2Path(cl string) string {
	cl = strings.Replace(cl, ".", string(os.PathSeparator), -1)
	if strings.HasSuffix(cl, "/class") {
		return cl[:len(cl)-6] + ".class"
	}
	return cl + ".class"
}

// Exec takes the command line arguments to be passed to the first
// "java" command executable found on the local system path. It's
// usefulness is that it will automatically check for any extracted
// cache in addition to host file system. Any ".class", ".jar", or
// ".java" file is allowed and the same syntax rules from Java are implied.
//
// Since Java class names are indistinguishable from option values, and
// since options can usually be any number of things including
// those that have values separated by space, and since Java
// implementation may have different options completely, this function
// requires that all options begin with dash (-) and use one of the
// no-space forms for making the value assignment (-Dfoo=bar, -foo:bar).
//
// This first argument to not begin with a dash is used as the class
// name, jar, or java file.
//
// All arguments after the main class/jar/java argument are passed as
// arguments to the main argument itself.
func Exec(cmd ...string) error {
	c := ParseCmd(cmd...)
	main := c.Name

	if strings.HasSuffix(c.Name, ".java") || strings.HasSuffix(c.Name, ".jar") {
		if c := Cached(c.Name); c != "" {
			main = c
		}
	}

	args := []string{"java"}
	args = append(args, c.Options...)
	args = append(args, main)
	args = append(args, c.Args...)

	return internal.Exec(args...)
}

// Out is the same as Exec but returns the standard output as a string
// and logs any errors.
func Out(cmd ...string) string {
	c := ParseCmd(cmd...)
	main := c.Name

	if strings.HasSuffix(c.Name, ".java") || strings.HasSuffix(c.Name, ".jar") {
		if c := Cached(c.Name); c != "" {
			main = c
		}
	}

	args := []string{"java"}
	args = append(args, c.Options...)
	args = append(args, main)
	args = append(args, c.Args...)

	return internal.Out(args...)
}

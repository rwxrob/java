/*
Package java uses Go embed.FS so that packages can be created to encapsulate
Java JAR, class, and raw source files that have been embedded into the
package with the default java executable on the host system. This
package includes a caching mechanism implemented in the Extract and
Cached functions so that the files need not be extracted with every run.
The java command invocation depends entirely on the version of java
installed on the host system and depends on properties and CLASSPATH
being maintained outside of this package itself.
*/
package java

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwxrob/fs"
	"github.com/rwxrob/java/internal"
)

// CacheDir is set to os.UserCacheDir() plus "gojavacache" by default at
// init time.
var CacheDir string

func init() {
	dir, err := os.UserCacheDir()
	if err == nil {
		CacheDir = filepath.Join(dir, "gojavacache")
	}
}

// Extract explicitly extracts all of an embedded file system into the
// CacheDir starting from the root path passed.
func Extract(fsys embed.FS, root string) error {
	os.MkdirAll(CacheDir, fs.ExtractDirPerms)
	return fs.ExtractEmbed(fsys, root, CacheDir)
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

// Run is a convenience function that takes a file path ending with
// ".class", ".jar", ".java" and calls RunClass, RunJar, or RunJava. If
// none of these suffixes match, and the string has a length greater
// than 10, assumes RunString.
func Run(it string) error {
	switch {
	case strings.HasSuffix(it, ".class"):
		return RunClass(it)
	case strings.HasSuffix(it, ".jar"):
		return RunJar(it)
	case strings.HasSuffix(it, ".java"):
		return RunJava(it)
	case len(it) > 10:
		return RunString(it)
	default:
		return fmt.Errorf("Unable to run anything but .java|.class|.jar or Java source string")
	}
	return nil
}

// RunString writes the Java source passed to a temporary file and then
// calls RunJava on it.
func RunString(src string) error {
	tmpfile := filepath.Join(os.TempDir(), internal.Isonan()+`.java`)
	if err := os.WriteFile(tmpfile, []byte(src), 0600); err != nil {
		return err
	}
	defer os.Remove(tmpfile)
	return RunJava(tmpfile)
}

// RunJava dynamically compiles and runs the java file passed but Java
// 11 or higher to be installed in host system. The file must end with
// ".java". If the java.FS is defined it will first be checked before
// checking the local file system.
func RunJava(file string) error {
	cached := Cached(file)
	if cached != "" {
		file = cached
	}
	return internal.Exec("java", file)
}

// RunJar calls "java -jar <file>". The file must end with ".jar". If
// the file is found to be Cached, the cached copy will be used.
// Otherwise, the current file system is assumed. The Java CLASSPATH and
// other properties must be set using other methods when wanted.
func RunJar(file string) error {
	cached := Cached(file)
	if cached != "" {
		file = cached
	}
	return internal.Exec("java", "-jar", file)
}

// RunClass calls "java <name>". The name must be a file ending with
// ".class" in the CLASSPATH. If a <name>.class file is found in the
// cache (see Cached) the directory of the path returned from Cached
// will be added to the front of the current CLASSPATH. Otherwise, the
// regular Java rules for finding class files apply.
func RunClass(name string) error {
	cached := Cached(name + ".class")
	if cached != "" {
		cp := os.Getenv("CLASSPATH")
		dir := filepath.Dir(cached)
		if len(cp) > 0 {
			os.Setenv("CLASSPATH", dir+string(os.PathListSeparator)+cp)
		} else {
			os.Setenv("CLASSPATH", dir)
		}
	}
	return internal.Exec("java", name)
}

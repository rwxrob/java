/*
Package java uses Go embed.FS so that packages can be created to encapsulate
Java JAR, class, and raw source files that have been embedded into the
package with the default java executable on the host system. This
package includes a caching mechanism implemented in the Extract and
Cached functions so that the files need not be extracted with every
execution.  The java command invocation depends entirely on the version
of java installed on the host system and observes CLASSPATH and other java-specific environment variables. Additional arguments may be passed to any of the Exec* or Out* functions and will be directly passed to the java command command line (without any shell expansion, of course). The Exec* functions map stdin/out/err to that of the OS while the Out* functions return a string with stdout and log any stderr (see internal/exec.go).

*/
package java

import (
	"embed"
	"fmt"
	"log"
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

// Exec is a convenience function that takes a file path ending with
// ".class", ".jar", ".java" and calls ExecClass, ExecJar, or ExecJava. If
// none of these suffixes match, and the string has a length greater
// than 10, assumes ExecString.
func Exec(it string, args ...string) error {
	switch {
	case strings.HasSuffix(it, ".class"):
		return ExecClass(it, args...)
	case strings.HasSuffix(it, ".jar"):
		return ExecJar(it, args...)
	case strings.HasSuffix(it, ".java"):
		return ExecJava(it, args...)
	case len(it) > 10:
		return ExecString(it, args...)
	default:
		return fmt.Errorf("Unable to run anything but .java|.class|.jar or Java source string")
	}
	return nil
}

// ExecString writes the Java source passed to a temporary file and then
// calls ExecJava on it.
func ExecString(src string, args ...string) error {
	tmpfile := filepath.Join(os.TempDir(), internal.Isonan()+`.java`)
	if err := os.WriteFile(tmpfile, []byte(src), 0600); err != nil {
		return err
	}
	defer os.Remove(tmpfile)
	return ExecJava(tmpfile, args...)
}

// ExecJava dynamically compiles and runs the java file passed but Java
// 11 or higher to be installed in host system. The file must end with
// ".java". If the java.FS is defined it will first be checked before
// checking the local file system.
func ExecJava(file string, args ...string) error {
	cached := Cached(file)
	if cached != "" {
		file = cached
	}

	cmd := []string{"java"}
	cmd = append(cmd, args...)
	cmd = append(cmd, file)

	return internal.Exec(cmd...)
}

// ExecJar calls "java -jar <file>". The file must end with ".jar" and
// must contain a manifest identifying which main class to run. If the
// file is found to be Cached, the cached copy will be used.  Otherwise,
// the current file system is assumed. The Java CLASSPATH and other
// properties must be set using other methods when wanted.
func ExecJar(file string, args ...string) error {
	cached := Cached(file)
	if cached != "" {
		file = cached
	}

	cmd := []string{"java"}
	cmd = append(cmd, args...)
	cmd = append(cmd, "-jar", file)

	return internal.Exec(cmd...)
}

// ExecClass calls "java <name>". The name must be a file ending with
// ".class" in the CLASSPATH. If a <name>.class file is found in the
// cache (see Cached) the directory of the path returned from Cached
// will be added to the front of the current CLASSPATH. Otherwise, the
// regular Java rules for finding class files apply.
func ExecClass(name string, args ...string) error {
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

	cmd := []string{"java"}
	cmd = append(cmd, args...)
	cmd = append(cmd, name)

	return internal.Exec(cmd...)
}

// Out is a convenience function that takes a file path ending with
// ".class", ".jar", ".java" and calls OutClass, OutJar, or OutJava. If
// none of these suffixes match, and the string has a length greater
// than 10, assumes OutString.
func Out(it string, args ...string) string {
	switch {
	case strings.HasSuffix(it, ".class"):
		return OutClass(it, args...)
	case strings.HasSuffix(it, ".jar"):
		return OutJar(it, args...)
	case strings.HasSuffix(it, ".java"):
		return OutJava(it, args...)
	case len(it) > 10:
		return OutString(it, args...)
	default:
		log.Println("Unable to run anything but .java|.class|.jar or Java source string")
	}
	return ""
}

// OutString writes the Java source passed to a temporary file and then
// calls OutJava on it.
func OutString(src string, args ...string) string {
	tmpfile := filepath.Join(os.TempDir(), internal.Isonan()+`.java`)
	if err := os.WriteFile(tmpfile, []byte(src), 0600); err != nil {
		log.Println(err)
	}
	defer os.Remove(tmpfile)
	return OutJava(tmpfile, args...)
}

// OutJava dynamically compiles and runs the java file passed but Java
// 11 or higher to be installed in host system. The file must end with
// ".java". If the java.FS is defined it will first be checked before
// checking the local file system.
func OutJava(file string, args ...string) string {
	cached := Cached(file)
	if cached != "" {
		file = cached
	}

	cmd := []string{"java"}
	cmd = append(cmd, args...)
	cmd = append(cmd, file)

	return internal.Out(cmd...)
}

// OutJar calls "java -jar <file>". The file must end with ".jar" and
// must include a manifest so that the correct class containing "main"
// will be used. If the file is found to be Cached, the cached copy will
// be used.  Otherwise, the current file system is assumed. The Java
// CLASSPATH and other properties must be set using other methods when
// wanted.
func OutJar(file string, args ...string) string {
	cached := Cached(file)
	if cached != "" {
		file = cached
	}

	cmd := []string{"java"}
	cmd = append(cmd, args...)
	cmd = append(cmd, "-jar", file)

	return internal.Out(cmd...)
}

// OutClass calls "java <name>". The name must be a file ending with
// ".class" in the CLASSPATH. If a <name>.class file is found in the
// cache (see Cached) the directory of the path returned from Cached
// will be added to the front of the current CLASSPATH. Otherwise, the
// regular Java rules for finding class files apply.
func OutClass(name string, args ...string) string {
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

	cmd := []string{"java"}
	cmd = append(cmd, args...)
	cmd = append(cmd, name)

	return internal.Out(cmd...)
}

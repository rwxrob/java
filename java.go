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

// FS, when assigned, will first be used to locate all calls before
// looking at the host file system.
var FS embed.FS

var CacheDir = filepath.Join(os.UserCacheDir, "gojavacache")

var extracted bool

func cacheFS() {
	if extracted {
		return
	}
	var zerofs embed.FS
	if FS == zerofs {
		// TODO needs to be written
		// fs.ExtractEmbed(FS, CacheDir)
		extracted = true
	}
}

// Cached returns the full path the extracted cache location of the file
// indicated by it.
func Cached(it string) string {
	cacheFS()
	path := filepath.Join(CacheDir, it)
	if extracted == true && fs.Exists(path) {
		return path
	}
	return ""
}

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

func RunString(it string) error {
	tmpfile := filepath.Join(os.TempDir(), internal.Isonan()+`.java`)
	if err := os.WriteFile(tmpfile, []byte(it), 0600); err != nil {
		return err
	}
	defer os.Remove(tmpfile)
	return RunJava(tmpfile)
}

func RunJava(it string) error {
	cached := Cached(it)
	if cached != "" {
		it = cached
	}
	return internal.Exec("java", it)
}

func RunJar(it string) error {
	cached := Cached(it)
	if cached != "" {
		it = cached
	}
	return internal.Exec("java", "-jar", it)
}

func RunClass(it string) error {
	cached := Cached(it)
	if cached != "" {
		it = cached
	}
	return internal.Exec("java", it)
}

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/clementauger/monparcours/cmd/goreadme/globci"
)

func main() {

	var verbose bool
	var dryrun bool
	var recursive bool

	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.BoolVar(&dryrun, "d", false, "dry run")
	flag.BoolVar(&recursive, "r", false, "recursive lookup")

	flag.Parse()

	stmt := []byte("package main")
	notStmt := []byte("package ")
	autogenStmt := []byte("// goreadme autogen ")
	skips := map[string]string{}
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if !recursive && info.Name() != "." {
				return filepath.SkipDir
			}
			return nil
		}
		dir := filepath.Dir(path)
		if _, ok := skips[dir]; ok {
			return nil
		}
		if strings.HasSuffix(path, ".go") {

			var isMain bool
			if i, err := fileContainLine(path, stmt, notStmt); err != nil {
				return err
			} else if i == 0 {
				isMain = true
			}

			// dont, it might skip interesting sub directory
			// if !isMain {
			// 	return filepath.SkipDir
			// }
			if !isMain {
				skips[dir] = path
				if verbose {
					log.Printf("directory %q is skipped\n", dir)
				}
				return nil
			}

			if verbose {
				log.Printf("main file %q\n", path)
			}

			isNew := false
			isAutoGen := false
			docFile := filepath.Join(dir, "doc.go")

			{
				g, err := globci.Glob(docFile)
				if err != nil {
					return err
				}
				if len(g) > 0 {
					docFile = g[0]
				}
			}

			var rmdFile string
			{
				g, err := globci.Glob(filepath.Join(dir, "README.md"))
				if err != nil {
					return err
				}
				if len(g) > 0 {
					rmdFile = g[0]
				}
			}

			if verbose {
				log.Printf("md file %q\n", rmdFile)
			}

			if _, err := os.Stat(rmdFile); os.IsNotExist(err) {
				return nil
			}

			if _, err := os.Stat(docFile); !os.IsNotExist(err) {
				i, err := fileContainLine(docFile, autogenStmt)
				if err != nil {
					return err
				}
				isAutoGen = i == 0
			} else {
				isNew = true
			}

			if verbose {
				log.Printf("isAutoGen %v isNew %v\n", isAutoGen, isNew)
			}

			if !isAutoGen && !isNew {
				return nil
			}
			f, err := os.Open(rmdFile)
			if err != nil {
				return err
			}
			defer f.Close()
			sc := bufio.NewScanner(f)
			p := ""
			out := new(bytes.Buffer)
			for sc.Scan() {
				l := sc.Text()
				isHead := strings.HasPrefix(l, "#") && !strings.HasPrefix(l, "##")

				if isHead && p != "" {
					fmt.Fprintf(out, "//md\n/*\n%v\n*/\n\n", p)
					p = ""
				}
				p += l + "\n"
			}
			if p != "" {
				fmt.Fprintf(out, "//md\n/*\n%v\n*/\n\n", p)
			}
			out.WriteString("package main\n\n")
			out.Write(autogenStmt)
			out.WriteString("\n")
			if dryrun {
				log.Printf("doc file to write %v\n", docFile)
				if verbose {
					log.Println()
					log.Printf("%v\n", out.String())
				}
			} else {
				ioutil.WriteFile(docFile, out.Bytes(), os.ModePerm)
			}
			skips[dir] = path
		}
		return nil
	})
}

func fileContainLine(path string, prefix ...[]byte) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return -1, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		l := sc.Bytes()
		for i, p := range prefix {
			if bytes.HasPrefix(l, p) {
				return i, nil
			}
		}
	}
	return -1, nil
}

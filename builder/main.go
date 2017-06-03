//lightriders builder
package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"time"

	"github.com/jessevdk/go-flags"
)

var (
	opts Options
)

const (
	debugFile  = "debug.go"
	botPackage = "lightRiders-starterBot-go"
)

type Options struct {
	Output  string `short:"o" long:"output-file" default:"bot.zip" description:"the location of the output file"`
	Input   string `short:"i" long:"input-dir" description:"the location of the source, set only if not found automatically"`
	KeepTmp bool   `short:"k" long:"keep-temp" description:"don't delete temp folder at the end, but writes it's location to the console"`
}

func main() {
	var err error
	var location, tmpDir, fileName string
	var zipFile, fp *os.File
	var zipPart *zip.Writer
	sayKeepTmp := true //if it needs to be told that tmp dir stays

	var fileContent []byte //parsing input arguments
	if _, err = flags.ParseArgs(&opts, os.Args); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	defer func() { // clean up
		if zipPart != nil {
			zipPart.Close()
		}
		if zipFile != nil {
			zipFile.Close()
		}
		if len(tmpDir) > 0 {
			if opts.KeepTmp {
				if sayKeepTmp {
					fmt.Println("Temp directory is available at", tmpDir)
				}
			} else {
				os.RemoveAll(tmpDir)
			}
		}
	}()

	var ok bool //getting source location
	if _, location, _, ok = runtime.Caller(0); !ok {
		if len(opts.Input) > 0 {
			location = opts.Input
			location = path.Dir(location)
		} else {
			panic(errors.New("Program code location not found!"))
		}
	}
	sourceDir := path.Join(path.Dir(location), "..", botPackage)

	//getting function list from debug.go, for removing calls to it
	if fileContent, err = ioutil.ReadFile(path.Join(sourceDir, debugFile)); err != nil {
		panic(err)
	}
	re := regexp.MustCompile(`func (\w+)\(`)
	var b bytes.Buffer
	b.WriteString(`(utils\.|.*?((`)
	for _, fileContent = range re.FindAll(fileContent, -1) {
		b.WriteString(string(fileContent[5 : len(fileContent)-1]))
		b.WriteRune('|')
	}
	b.Truncate(b.Len() - 1)
	b.WriteString(`)\(|"github\.com/vendelin8/lightriders-starterbot-golang/utils").*)`)
	re = regexp.MustCompile(b.String())

	if tmpDir, err = ioutil.TempDir("", botPackage); err != nil {
		panic(err)
	}

	var fileInfos []os.FileInfo //copying modified main files to tmp folder
	var info os.FileInfo
	if fileInfos, err = ioutil.ReadDir(sourceDir); err != nil {
		panic(err)
	}
	for _, info = range fileInfos {
		fileName = info.Name()
		if filepath.Ext(fileName) != ".go" || fileName == debugFile {
			continue
		}
		if fileContent, err = ioutil.ReadFile(path.Join(sourceDir, fileName)); err != nil {
			panic(err)
		}
		fileContent = re.ReplaceAll(fileContent, []byte{})
		if err = ioutil.WriteFile(path.Join(tmpDir, fileName), fileContent, 0644); err != nil {
			panic(err)
		}
	}

	var bufWriter *bufio.Writer //copying util files
	oldPkgStr := "package utils"
	newPkgStr := "package main"
	sourceDir = path.Join(path.Dir(location), "..", "utils")
	if fileInfos, err = ioutil.ReadDir(sourceDir); err != nil {
		panic(err)
	}
	for _, info = range fileInfos {
		fileName = info.Name()
		if filepath.Ext(fileName) != ".go" || fileName == "const.go" {
			continue
		}
		if fileContent, err = ioutil.ReadFile(path.Join(sourceDir, fileName)); err != nil {
			panic(err)
		}
		if fp, err = os.Create(path.Join(tmpDir, fileName)); err != nil {
			panic(err)
		}
		defer fp.Close()
		bufWriter = bufio.NewWriter(fp)
		if _, err = bufWriter.WriteString(newPkgStr); err != nil {
			panic(err)
		}
		if _, err = bufWriter.Write(fileContent[len(oldPkgStr):]); err != nil {
			panic(err)
		}
		bufWriter.Flush()
	}

	cmd := exec.Command("gofmt", "-w", tmpDir) //reformatting
	cmd.Run()
	cmd = exec.Command("go", "build", "-o", "bot") //building for test if it works
	cmd.Dir = tmpDir
	if err = cmd.Run(); err != nil {
		fmt.Println("Something is wrong with your build. The temp folder will NOT be deleted.",
			"Try manually:")
		fmt.Println("cd", tmpDir)
		fmt.Println("go build -o bot")
		opts.KeepTmp = true
		sayKeepTmp = false
		panic(err)
	}
	os.Remove(path.Join(tmpDir, "bot"))

	if zipFile, err = os.Create(opts.Output); err != nil {
		panic(err)
	}
	zipPart = zip.NewWriter(zipFile)
	if fileInfos, err = ioutil.ReadDir(tmpDir); err != nil {
		panic(err)
	}
	var zipWriter io.Writer
	t := time.Now()
	for _, info = range fileInfos {
		fileName = info.Name()
		if fp, err = os.Open(path.Join(tmpDir, fileName)); err != nil {
			panic(err)
		}
		defer fp.Close()
		header := &zip.FileHeader{
			Name:   fileName,
			Method: zip.Deflate,
		}
		header.SetModTime(t)
		if zipWriter, err = zipPart.CreateHeader(header); err != nil {
			panic(err)
		}
		if _, err = io.Copy(zipWriter, fp); err != nil {
			panic(err)
		}
	}
}

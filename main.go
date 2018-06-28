package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"

	docopt "github.com/docopt/docopt-go"
)

const (
	version = "0.1.0-dev"
)

func main() {
	usage := `orga is a cli to work with golang webassembly compilation target.

Usage:
  orga generate (<directory> [--force|--file=<name>])
  orga serve (<directory> [--addr=<host:port>])
  orga -h | --help
  orga -v | --version

Options:
  -h --help            Show this screen.
  --version            Show version.
  --force              Overwrite existing file.
  --file=<name>        Name of the wasm file [default: main.wasm].
  --addr=<host:port>   Name of the wasm file [default: :3000].
	`
	arguments, err := docopt.ParseDoc(usage)
	if err != nil {
		docopt.PrintHelpAndExit(err, usage)
	}

	fmt.Printf("DEBUG: arguments %#v\n", arguments)

	if arguments["-v"] == true || arguments["--version"] == true {
		fmt.Println(version)
	}

	if arguments["generate"] == true {
		directory := arguments["<directory>"].(string)
		force := arguments["--force"].(bool)
		filename := arguments["--file"].(string)
		generateJsHtmlFiles(directory, filename, force)
	}

	if arguments["serve"] == true {
		addr := arguments["--addr"].(string)
		directory := arguments["<directory>"].(string)
		serve(directory, addr)
	}
}

func generateJsHtmlFiles(filepath, name string, force bool) {
	fmt.Printf("Generating files: generateJsHtmlFiles(%s, %s, %v)\n", filepath, name, force)
	data := make(map[string]string)
	for k, v := range defaultData {
		if k == "wasmFilename" {
			data["wasmFilename"] = name
		} else {
			data[k] = v
		}
	}

	templatedFiles := make(map[string]templatedFile)
	templatedFiles["go_js_wasm_exec"] = templatedFile{
		content: go_js_wasm_exec_tmpl,
		data:    data,
	}
	templatedFiles["index.html"] = templatedFile{
		content: wasm_exec_html_tmpl,
		data:    data,
	}
	templatedFiles["wasm_exec.js"] = templatedFile{
		content: wasm_exec_js_tmpl,
		data:    data,
	}

	filepath = path.Clean(filepath)

	var err error
	var tpl *template.Template
	var content strings.Builder
	for filename, tf := range templatedFiles {
		filename = path.Join(filepath, filename)
		tpl, err = template.New(filename).Parse(tf.content)
		if err != nil {
			log.Fatal("could not parse: ", filename, err)
		}
		fmt.Printf("file: %s ; data %#v", filename, tf.data)
		err = tpl.Execute(&content, tf.data)
		if err != nil {
			log.Fatal("Could not Execute the template", err)
		}

		err = writeFileFromTemplate(filename, content.String(), force)
		if err != nil {
			log.Fatal("Could not write file: ", err)
		}
		content.Reset()
	}

}

func writeFileFromTemplate(filename, content string, force bool) error {
	fmt.Printf("Generating file: %s, %v\n", filename, force)

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	stats, err := f.Stat()
	if err != nil {
		return err
	}
	if stats.Size() != 0 && !force {
		fmt.Printf("noop on file: %s", filename)
		return nil
	}

	_, err = f.WriteString(content)
	return err
}

func serve(directory, addr string) {
	fmt.Printf("Starting  an http server: %s", addr)

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(directory)))
	log.Fatal(http.ListenAndServe(addr, mux))

}

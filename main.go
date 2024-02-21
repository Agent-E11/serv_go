package main

import (
	"flag"
	"fmt"
	"net/http"
	"text/template"
)

func main() {
    port := flag.Int("p", 8000, "port to listen on")
    // TODO: When true, it should dynamically search for any url given to it
    permissive := flag.Bool("permissive", false, "whether it should allow routes to any file")
    // TODO: Might be unnecessary, because it seems to always listen externally
    external := flag.Bool("ext", false, "whether it should listen on 0.0.0.0")
    flag.Parse()

    fmt.Println("Non-flag arguments:", flag.Args())

    
    // Debug stuff
    if *permissive {
        fmt.Println("Permissive")
    } else {
        fmt.Println("Not permissive")
    }

    // Set the file to the first parameter (if specified)
    var file string
    if flag.Arg(0) == "" {
        file = "index.html"
    } else {
        file = flag.Arg(0)
    }
    fmt.Printf("The root file is `%s`\n", file)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("Serving /")
        tmpl, err := template.ParseFiles(file)
        if err != nil {
            fmt.Printf("The file %s does not exist\n", file)
            return
        }

        err = tmpl.Execute(w, nil)
        if err != nil {
            fmt.Printf("Error executing template: %v", err)
        }
    })

    // TODO: Search for files using filepath.Glob,
    // And then make handlers for every file found
    for _, f := range flag.Args() {
        http.HandleFunc(
            fmt.Sprintf("/%s", f),
            generateHandler(f),
        )
    }

    var portString string
    if *external {
        portString = fmt.Sprintf("0.0.0.0:%d", *port)
    } else {
        portString = fmt.Sprintf(":%d", *port)
    }

    fmt.Println("Listening on port", portString)
    http.ListenAndServe(portString, nil)
}

// generateHandler generates a basic http.HandleFunc given a file name (e.g. path/to/index.html)
// Will use the file's name to look for a template, and will execute that template
// without passing it any data
func generateHandler(file string) http.HandlerFunc {
    fmt.Printf("Generating handler for %s\n", file)
    return func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("Serving /%s\n", file)
        tmpl, err := template.ParseFiles(file)
        if err != nil {
            fmt.Printf("The file %s does not exist\n", file)
            return
        }

        err = tmpl.Execute(w, nil)
        if err != nil {
            fmt.Printf("Error executing template: %v", err)
        }
    }
}
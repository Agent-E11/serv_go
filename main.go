package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "text/template"
)

func main() {
    port := flag.Int("p", 8000, "port to listen on")
    // TODO: When true, it should dynamically search for any url given to it
    permissive := flag.Bool("permissive", false, "whether it should allow routes to any file")
    // TODO: Might be unnecessary, because it seems to always listen externally
    external := flag.Bool("ext", false, "whether it should listen on 0.0.0.0")
    flag.Parse()

    log.Println("Non-flag arguments:", flag.Args())

    
    // Debug stuff
    if *permissive {
        log.Println("Permissive")
    } else {
        log.Println("Not permissive")
    }

    // Set the file to the first parameter
    // if specified, and if it's not a glob pattern. Else,
    // set to index.html
    var file string
    if flag.Arg(0) == "" {
        file = "index.html"
    } else if strings.Contains(flag.Arg(0), "*") ||
        strings.Contains(flag.Arg(0), "?") ||
        strings.Contains(flag.Arg(0), "[") {
        file = "index.html"
    } else {
        file = flag.Arg(0)
    }
    log.Printf("The root file is `%s`\n", file)

    // Create the default route
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.Println("Serving /")
        log.Printf("Value of r.URL.Path: %v\n", r.URL.Path)

        // If the path is not the root, return a 404
        if r.URL.Path != "/" {
            log.Println("Before not found error")
            http.NotFound(w, r)
            log.Println("After not found error")
            return
        }

        tmpl, err := template.ParseFiles(file)
        if err != nil {
            log.Printf("The file %s does not exist\n", file)
            return
        }

        err = tmpl.Execute(w, nil)
        if err != nil {
            log.Printf("Error executing template: %v\n", err)
        }
    })

    // Search for files using a glob pattern
    // and then make handlers for every file found
    for _, descriptor := range flag.Args() {

        files, err := filepath.Glob(descriptor)
        if err != nil {
            log.Printf("Error reading glob: %v\n", err)
            continue
        }
        if files == nil {
            log.Println("No files described by", descriptor)
            continue
        }

        for _, f := range files {
            // TODO: Only make handler if file is not a directory

            fileInfo, err := os.Stat(f)
            if err != nil {
                log.Println("Error getting file info:", err)
                continue
            }

            if !fileInfo.IsDir() {
                http.HandleFunc(
                    fmt.Sprintf("/%s", f),
                    generateHandler(f),
                )
            } else {
                log.Println(f, "is a directory, not creating handler")
            }
        }
    }

    // Create portString
    var portString string
    if *external {
        portString = fmt.Sprintf("0.0.0.0:%d", *port)
    } else {
        portString = fmt.Sprintf(":%d", *port)
    }

    // Listen for http requests
    log.Println("Listening on port", portString)
    http.ListenAndServe(portString, nil)
}

// generateHandler generates a basic http.HandleFunc given a file name (e.g. path/to/index.html)
// Will use the file's name to look for a template, and will execute that template
// without passing it any data
func generateHandler(file string) http.HandlerFunc {
    log.Printf("Generating handler for %s\n", file)
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Serving /%s\n", file)
        tmpl, err := template.ParseFiles(file)
        if err != nil {
            log.Printf("The file %s does not exist\n", file)
            return
        }

        err = tmpl.Execute(w, nil)
        if err != nil {
            log.Printf("Error executing template: %v\n", err)
        }
    }
}

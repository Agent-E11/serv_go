package main

import (
    "flag"
    "fmt" // TODO: Change all debug `fmt`s prints to `log`
    "net/http"
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

    fmt.Println("Non-flag arguments:", flag.Args())

    
    // Debug stuff
    if *permissive {
        fmt.Println("Permissive")
    } else {
        fmt.Println("Not permissive")
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
    fmt.Printf("The root file is `%s`\n", file)

    // Create the default route
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("Serving /")
        fmt.Printf("Value of r.URL.Path: %v\n", r.URL.Path)

        // If the path is not the root, return a 404
        if r.URL.Path != "/" {
            fmt.Println("Before not found error")
            http.NotFound(w, r)
            fmt.Println("After not found error")
            return
        }

        tmpl, err := template.ParseFiles(file)
        if err != nil {
            fmt.Printf("The file %s does not exist\n", file)
            return
        }

        err = tmpl.Execute(w, nil)
        if err != nil {
            fmt.Printf("Error executing template: %v\n", err)
        }
    })

    // Search for files using a glob pattern
    // and then make handlers for every file found
    for _, descriptor := range flag.Args() {

        files, err := filepath.Glob(descriptor)
        if err != nil {
            fmt.Printf("Error reading glob: %v\n", err)
            continue
        }
        if files == nil {
            fmt.Println("No files described by", descriptor)
            continue
        }

        for _, f := range files {
            // TODO: Only make handler if file is not a directory
            http.HandleFunc(
                fmt.Sprintf("/%s", f),
                generateHandler(f),
            )
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
            fmt.Printf("Error executing template: %v\n", err)
        }
    }
}

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
    var port int
    flag.IntVar(&port, "p", 8000, "port to listen on")
    var permissive bool
    flag.BoolVar(&permissive, "permissive", false, "whether it should allow routes to any file")
    // TODO: Might be unnecessary, because it seems to always listen externally
    var external bool
    flag.BoolVar(&external, "ext", false, "whether it should listen on 0.0.0.0")
    flag.Parse()

    log.Println("Non-flag arguments:", flag.Args())

    // Debug stuff
    if permissive {
        log.Println("Permissive")
    } else {
        log.Println("Not permissive")
    }

    // Set the file to the first parameter if specified, and if it's not a
    // glob pattern. Else, set to index.html
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

        path := strings.TrimLeft(r.URL.Path, "/.")

        if path == "" {
            tmpl, err := template.ParseFiles(file)
            if err != nil {
                log.Printf("The file %s does not exist\n", file)
                return
            }

            tmpl.Execute(w, nil)
            return
        } else if !permissive {
            log.Println("Not permissive, not attempting to find file:", path)
            http.NotFound(w, r)
            return
        }

        tmpl, err := template.ParseFiles(path)
        if err != nil {
            log.Println("Couldn't find file:", path)
            log.Println("Error:", err)
            http.NotFound(w, r)
            return
        }
        log.Println("Serving dynamic file:", path)
        tmpl.Execute(w, nil)
    })

    // Search for files using a glob pattern
    // and then make handlers for every file found
    var files []string
    for _, descriptor := range deDuplicate(flag.Args()) {
        fs, err := filepath.Glob(descriptor)
        if err != nil {
            log.Printf("Error reading glob: %v\n", err)
            continue
        }
        files = append(files, fs...)
    }
    log.Printf("Files: %v\n", files)
    log.Printf("Files: %v\n", deDuplicate(files))

    for _, f := range deDuplicate(files) {
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

    // Create portString
    var portString string
    if external {
        portString = fmt.Sprintf("0.0.0.0:%d", port)
    } else {
        portString = fmt.Sprintf(":%d", port)
    }

    // Listen for http requests
    log.Println("Listening on port", portString)
    log.Fatal(http.ListenAndServe(portString, nil))
}

// generateHandler generates a basic http.HandleFunc given a file name (e.g. path/to/index.html).
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

// deDuplicate removes duplicate elements from a slice
func deDuplicate[T comparable](sliceList []T) []T {
    keys := make(map[T]bool)
    list := []T{}
    for _, item := range sliceList {
        if _, value := keys[item]; !value {
            keys[item] = true
            list = append(list, item)
        }
    }
    return list
}

package main

import (
	"flag"
	"fmt"
	"net/http"
	"text/template"
)

func main() {
    port := flag.Int("p", 8000, "port to listen on")
    permissive := flag.Bool("permissive", false, "whether it should allow all routes to all files")
    external := flag.Bool("ext", false, "whether it should listen on 0.0.0.0")
    flag.Parse()

    fmt.Println("Non-flag arguments:", flag.Args())

    if *permissive {
        fmt.Println("Permissive")
    } else {
        fmt.Println("Not permissive")
    }

    var file string
    if flag.Arg(0) == "" {
        file = "index.html"
    } else {
        file = flag.Arg(0)
    }
    fmt.Printf("The file is `%s`", file)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("Serving /")
        tmpl, err := template.ParseFiles(file)
        if err != nil {
            fmt.Printf("The file %s does not exist\n", file)
            return
        }

        tmpl.Execute(w, nil)
    })

    var portString string
    if *external {
        portString = fmt.Sprintf("0.0.0.0:%d", *port)
    } else {
        portString = fmt.Sprintf(":%d", *port)
    }

    fmt.Println("Listening on port", portString)
    http.ListenAndServe(portString, nil)
}

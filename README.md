## Codis

Codis is a search engine specifically build for source code. It's built as a terminal user interface so it can run in the command line directly. It doesn't require any runtime dependencies and it's written fully in Go (hence it's quite performant). This version is meant as a proof-of-concept where I validate my initial idea and collect feedback. Hence, feel free to download this program and share feedback. I'll add some videos here to explain how the tool works.

Abstract: [VIDEO](https://www.loom.com/share/bed8033b20bd4692b0866f58d84285ec?sid=d9863c27-f677-4556-808c-a8470379b308)  
How to install/use: clone the git repository and cd into its root directory. Next, install the dependencies (only bubbletea framework) with `go get ,`. Thereafter, you can build/run the executable using `go build main.go` or run it through the interpreter using `go run main.go`. Instructions/keyboard shortcuts are available through pressing `:` (which brings you in command mode) and typing `help`.

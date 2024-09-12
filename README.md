# Serve NextJS app from Go

This little library makes it easy to serve NextJS apps that have been built using the `output: "export"` option from your Go app.

The main purpose is to resolve dynamic route segments and match the request path to the correct HTML file.
For this it does some basic pattern matching which probably doesn't cover all edge cases.

I mainly built this so I can bundle a NextJS app together with the Go backend in a single binary.

```go
import gonextstatic "github.com/merlinfuchs/go-next-static"

func main() {
    // Use go:embed instead to embed the files in the binary
	nextFiles := os.DirFS("out")

	handler, err := gonextstatic.NewHandler(nextFiles.(fs.StatFS))
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(":8080", handler)
}
```

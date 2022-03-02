# hasher

It's a small tool which makes http requests and prints the address of the request along with the MD5 hash of the response.
# Usage

```
Usage: go run main.go [OPTIONS] URLS 

Options:

  -parallel
    	specifies max number of goroutines (10 by default)

Examples:
    go run main.go adjust.com google.com vk.com
    go run main.go -parallel 3 adjust.com google.com vk.com yahoo.com twitter.com
    go run main.go -parallel=5 adjust.com google.com vk.com
```
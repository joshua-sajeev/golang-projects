## Hello World Web Server: Create a basic web server that responds with "Hello, World!".


```go
import (
	"fmt"
	"log"
	"net/http"
)
```
These are the necessary imports 

Package `fmt` implements formatted I/O with functions analogous to C's printf and scanf.  The format 'verbs' are derived from C's but are simpler.

Package `log` implements a simple logging package.

Package `http` provides HTTP client and server implementations.


```go
name := r.URL.Path[1:]
```
`r.URL` is a struct of type `*url.URL` that represents the parsed URL of the request.
It's a string that represents the portion of the URL after the domain and the protocol (e.g., http://example.com), excluding the query parameters and fragment.

For example

```
GET /greet/john HTTP/1.1
Host: example.com
```
The r.URL.Path will be /greet/john.

So for the above code
GET / HTTP/1.1
will be 
Hello, World!

GET /alice HTTP/1.1
will be
Hello, alice!


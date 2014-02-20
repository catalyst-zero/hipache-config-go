hipache-config-go
=================

[![GoDoc](https://godoc.org/github.com/catalyst-zero/hipache-config-go?status.png)](https://godoc.org/github.com/catalyst-zero/hipache-config-go)

Go library to configure hipache backends via redis.

## Usage

```
import hipache "github.com/catalyst-zero/hipache-config-go"

client, err := hipache.DialHipacheConfig("127.0.0.1:6379")

client.BindingCreate("www.example.com")
client.BindingAddHost("www.example.com", "10.0.0.2:80")
client.BindingAddHost("www.example.com", "10.0.0.3:80")
```

For the complete API, visit [godoc.org](http://godoc.org/github.com/catalyst-zero/hipache-config-go).

There is also an in-memory implementation available via `hipache.NewLocalHipacheConfig()` for testing.

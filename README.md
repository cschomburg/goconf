goconf
======

Package conf provides a simple way to read and write config files.
It supports JSON by default, but can work with every encoding that has
Marshal / Unmarshal functions, e.g. most of the standard library.

Documentation available under: http://godoc.org/github.com/xconstruct/goconf

### Install ###

	go get "github.com/xconstruct/goconf"

### Example ###

```go
// Your struct containing the config fields
type TestConfig struct {
	String string
	Number int
	Sub struct {
		Field string
	}
}


// Creates JSON config under ~/.config/goconftest/config.json
cfg := conf.Build().App("goconftest").JSON.Create()

// Write config data
myConfig := TestConfig{"Just testing", 123, struct{Field string}{"test"}}
err := cfg.Write(myConfig)

// Read config data
var readConfig TestConfig
err := conf.Read(&readConfig)
```

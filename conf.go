// Package conf provides a simple way to read and write config files.
// It supports JSON by default, but can work with every encoding that has
// Marshal / Unmarshal functions, e.g. most of the standard library.
//
// Normally, you want to use the builder to specify your config format.
//
//    // Sets up a JSON file "config.json" in your current directory.
//    cfg := conf.Build().JSON().Create()
//
//
//    // Sets up a JSON file "~/.config/yourapp/config.json".
//    cfg := conf.Build().App("yourapp").JSON().Create()
//
// You can then read and write this file from/into your own data types.
//    myValue := struct{One, Two string}{"Hello", "World"}
//
//	  err := cfg.Write(myValue)
//	  err := cfg.Read(&myValue)
package conf

import (
	"io/ioutil"
	"os"
	"encoding/json"
	"errors"
)

// Marshal returns a encoding of v.
type MarshalFunc func(v interface{}) ([]byte, error)
// Unmarshal parses the data and stores the result in the value pointed to by v.
type UnmarshalFunc func(data []byte, v interface{}) error

// Context holds all information to access a specific config file
type Context struct {
	Directory string
	File string
	Marshal MarshalFunc
	Unmarshal UnmarshalFunc
}

var (
	ErrNoMarshal = errors.New("Context has no marshal func")
	ErrNoUnmarshal = errors.New("Context has no Unmarshal func")
)

// Read reads the config file into the value pointed to by conf.
func (c *Context) Read(conf interface{}) error {
	f, err := os.Open(c.Directory + "/" + c.File)
	if err != nil {
		path := err.(*os.PathError)
		if path != nil && path.Err == os.ErrNotExist {
			return nil
		}
		return err
	}
	defer f.Close()
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	if c.Unmarshal == nil {
		return ErrNoUnmarshal
	}
	return c.Unmarshal(bytes, conf)
}

// Write writes conf into the config file of the context.
func (c *Context) Write(conf interface{}) error {
	if err := os.MkdirAll(c.Directory, 0777); err != nil {
		return err
	}
	f, err := os.OpenFile(c.Directory + "/" + c.File, os.O_WRONLY | os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	if c.Marshal == nil {
		return ErrNoMarshal
	}
	bytes, err := c.Marshal(conf)
	if err != nil {
		return err
	}
	_, err = f.Write(bytes)
	return err
}

// Builder helps create contexts.
type Builder struct {
	ctx Context
}

// Directory sets the directory of the config file.
func (b *Builder) Directory(dir string) *Builder {
	b.ctx.Directory = dir
	return b
}

// File sets the name of the config file.
func (b *Builder) File(file string) *Builder {
	b.ctx.File = file
	return b
}

// App sets the directory of a config file to the appName in the
// user config directory, e.g. ~/.config/appName
func (b *Builder) App(appName string) *Builder {
	b.ctx.Directory = os.Getenv("XDG_CONFIG_HOME") + "/" + appName
	return b
}

// Marshaller sets the functions to use for encoding/decoding.
func (b *Builder) Marshaller(m MarshalFunc, u UnmarshalFunc) *Builder {
	b.ctx.Marshal = m
	b.ctx.Unmarshal = u
	return b
}

func jsonMarshalIndent(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "    ")
}

// JSON sets the encoding to the JSON format.
func (b *Builder) JSON() *Builder {
	if b.ctx.File == "" {
		b.ctx.File = "config.json"
	}
	return b.Marshaller(jsonMarshalIndent, json.Unmarshal)
}

// Create creates the context.
func (b *Builder) Create() *Context {
	if b.ctx.Directory == "" {
		b.ctx.Directory = "."
	}
	if b.ctx.File == "" {
		b.ctx.File = "config"
	}
	return &b.ctx
}

// Returns a new builder to build a config context.
func Build() *Builder {
	return &Builder{}
}

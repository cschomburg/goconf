package conf

import (
	"testing"
)

type TestConfig struct {
	String string
	Number int
	Sub struct {
		Field string
	}
}

func TestReadWrite(t *testing.T) {
	// Initialize config context
	conf := Build().JSON().Create()

	// Example config
	cfg := TestConfig{"Just testing", 123, struct{Field string}{"test"}}

	// Write config
	err := conf.Write(cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Read config
	var cfgRead TestConfig
	err = conf.Read(&cfgRead)
	if err != nil {
		t.Fatal(err)
	}

	if cfg != cfgRead {
		t.Errorf("Configs differ: %v, %v", cfg, cfgRead)
	}
}

func TestGlobal(t *testing.T) {
	// Initialize config context
	conf := Build().App("goconftest").JSON().Create()

	// Example config
	cfg := TestConfig{"Just testing", 123, struct{Field string}{"test"}}

	// Write config
	err := conf.Write(cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Read config
	var cfgRead TestConfig
	err = conf.Read(&cfgRead)
	if err != nil {
		t.Fatal(err)
	}

	if cfg != cfgRead {
		t.Errorf("Configs differ: %v, %v", cfg, cfgRead)
	}
}

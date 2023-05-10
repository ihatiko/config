package config

import (
	"fmt"
	"github.com/ihatiko/config/example"
	"os"
	"testing"
)

func TestGetConfig(t *testing.T) {
	cases := []struct {
		in       func()
		env      string
		excepted func(config *example.Config)
	}{
		{
			env: "CASE1_NAME",
			excepted: func(config *example.Config) {
				if config.Case1.Name != 777 {
					panic("os.Setenv(\"CASE1_NAME\", \"321\")")
				}
			},
			in: func() {
				os.Setenv("CASE1_NAME", "777")
			},
		},
		{
			env: "TEST1_0",
			excepted: func(config *example.Config) {
				if config.TestArray[0] != "321" && config.TestArray[1] != "321" && config.TestArray[2] != "321" {
					panic(os.Setenv("TESTARRAY_0", "321"))
				}
			},
			in: func() {
				os.Setenv("TESTARRAY_0", "321")
				os.Setenv("TESTARRAY_1", "321")
				os.Setenv("TESTARRAY_2", "321")
			},
		},
		{
			env: "TESTARRAY",
			excepted: func(config *example.Config) {
				if config.TestArray[0] != "1" && config.TestArray[1] != "2" && config.TestArray[2] != "3" {
					panic("os.Setenv(\"TESTARRAY\", \"1,2,3\")")
				}
			},
			in: func() {
				os.Setenv("TESTARRAY", "1,2,3")
			},
		},
		{
			env: "TESTARRAY",
			excepted: func(config *example.Config) {
				if config.TestArray[0] != "1" && config.TestArray[1] != "2" && config.TestArray[2] != "3" {
					panic("os.Setenv(\"TESTARRAY\", \"1,2,3\")")
				}
			},
			in: func() {
				os.Setenv("TESTARRAY", "1,2,3,4,5,6")
			},
		},
		{
			env: "os.Setenv(\"CASE3_0_NAME\", \"321\")",
			excepted: func(config *example.Config) {
				if config.Case3[0].Name != 321 {
					panic("os.Setenv(\"TESTARRAY\", \"1,2,3\")")
				}
			},
			in: func() {
				os.Setenv("CASE3_0_NAME", "321")
			},
		},
	}

	for _, caseTest := range cases {
		caseTest.in()
		cfg, err := GetConfig[example.Config]("./example/config")
		if err != nil {
			panic(err)
		}
		caseTest.excepted(cfg)
		os.Clearenv()
		fmt.Println(cfg)
	}

}

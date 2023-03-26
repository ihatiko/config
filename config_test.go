package config

import (
	"fmt"
	"github.com/ihatiko/config/example"
	"os"
	"testing"
)

func TestGetConfig(t *testing.T) {
	//os.Setenv("FIELD1_TEST", "HEELO")
	//os.Setenv("CASE2_ARRAYSTRUCT_0_NAME", "000")
	//os.Setenv("CASE2_ARRAYSTRUCT_1_NAME", "000")
	//os.Setenv("CASE2_ARRAYSTRUCT_1_VALUE2", "000")
	os.Setenv("TEST1_0", "TEST12313123")
	os.Setenv("CASE1_ARRAY_0", "HLLOWORLDLASDASD")
	os.Setenv("CASE1_ARRAY_1", "HLLOWORLDLASDASD")
	os.Setenv("CASE1_NAME", "321")
	cfg, err := GetConfig[example.Config]("./example/config")
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)
}

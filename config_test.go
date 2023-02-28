package config

import (
	"fmt"
	"github.com/ihatiko/config/example"
	"os"
	"testing"
)

func TestGetConfig(t *testing.T) {
	os.Setenv("CASE2_ARRAYSTRUCT_1_NAME", "HLLOWORLDLASDASD")
	os.Setenv("CASE1_ARRAY_1", "HLLOWORLDLASDASD")
	cfg, err := GetConfig[example.Config]("./example/config")
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)
}

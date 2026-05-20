package main

import (
	"github.com/adelmoradian/terraformer/pkg/blocks"
	"github.com/adelmoradian/terraformer/pkg/module"
)

func main() {
	mainBlock := blocks.NewTerraformBlock(">= 1.9.0, < 2.0.0").
		AddProvider(blocks.Provider{
			Name:    "gitlab",
			Source:  "gitlabhq/gitlab",
			Version: "~> 18.0"}).
		AddProvider(blocks.Provider{
			Name:    "null",
			Source:  "hashicorp/null",
			Version: "~> 3.0"}).SetRemoteBackend("http", map[string]string{"foo": "bar"})

	varBlock := blocks.NewVarBlock("test", "this is a test").SetType(`list(object({
    internal = number
    external = number
    protocol = string
  }))`).SetValidation(`length(var.image_id) > 4 && substr(var.image_id, 0, 4) == "ami-"`, "some msg")

	someBlock := blocks.NewLocals().
		Add("k1", "value").
		Add("k2", 123).
		Add("k3", false).
		Add("k4", []any{1, 2.1}).
		Add("k5", map[string]any{"a": "bar", "b": "foo"})

	someRes := blocks.NewResource("label", "tag").
		Add("k1", "value").
		Add("k2", 123).
		Add("k3", false).
		Add("k4", []any{1, 2.1}).
		Add("k5", map[string]any{"a": "bar", "b": "foo"})

	someBlock.AddStatement("k6", someBlock.KeyRef("k4")+"[1]")
	module.New("./modules").AddFile("main.tf", mainBlock, varBlock, someBlock, someRes).Create()
}

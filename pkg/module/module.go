package module

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adelmoradian/terraformer/pkg/blocks"
)

type tfModule struct {
	moduleDir string
	tfFiles   []tfFile
}

type tfFile struct {
	fileName string
	blocks   []blocks.TFBlock
}

func New(moduleDir string) *tfModule {
	return &tfModule{moduleDir: moduleDir}
}

func (x *tfModule) AddFile(fileName string, blocks ...blocks.TFBlock) *tfModule {
	x.tfFiles = append(x.tfFiles, tfFile{fileName: fileName, blocks: blocks})
	return x
}

func (x *tfModule) Create() {
	if err := os.MkdirAll(x.moduleDir, os.ModePerm); err != nil {
		fmt.Println("could not create module dir", err.Error())
		os.Exit(1)
	}

	for _, f := range x.tfFiles {
		if len(f.blocks) == 0 {
			return
		}

		path := filepath.Join(x.moduleDir, f.fileName)
		content := f.blocks[0].String()
		for i := 1; i < len(f.blocks); i++ {
			content = fmt.Sprintf("%s\n%s", content, f.blocks[i].String())
		}

		if err := os.WriteFile(path, []byte(content), os.ModePerm); err != nil {
			fmt.Println("could not write", err.Error())
			os.Exit(1)
		}
	}
}

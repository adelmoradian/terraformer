package blocks

import (
	"fmt"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclwrite"
)

type dataBlock struct {
	*hclwrite.File
	keys          []string
	resourceType  string
	resourceLabel string
}

func NewData(resourceType, resourceLabel string) *dataBlock {
	file := hclwrite.NewEmptyFile()
	body := file.Body()
	body.AppendNewBlock("data", []string{resourceType, resourceLabel})
	return &dataBlock{file, nil, resourceType, resourceLabel}
}

func (x *dataBlock) Add(key string, value any) *dataBlock {
	x.keys = append(x.keys, key)
	body := x.Body()
	lb := body.Blocks()[0]
	lbBody := lb.Body()
	lbBody.SetAttributeValue(key, toCtyVal(value))
	return x
}

func (x *dataBlock) AddStatement(key string, value string) *dataBlock {
	body := x.Body()
	lb := body.Blocks()[0]
	lbBody := lb.Body()
	lbBody.SetAttributeTraversal(key, hcl.Traversal{hcl.TraverseRoot{Name: value}})
	return x
}

func (x *dataBlock) Self() string {
	return fmt.Sprintf("data.%s.%s", x.resourceType, x.resourceLabel)
}

func (x *dataBlock) Type() string {
	return "data"
}

func (x *dataBlock) String() string {
	return string(x.Bytes())
}

package blocks

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclwrite"
)

type resourceBlock struct {
	*hclwrite.File
	keys          []string
	resourceType  string
	resourceLabel string
}

func NewResource(resourceType, resourceLabel string) *resourceBlock {
	return &resourceBlock{hclwrite.NewEmptyFile(), nil, resourceType, resourceLabel}
}

func (x *resourceBlock) Add(key string, value any) *resourceBlock {
	x.keys = append(x.keys, key)
	body := x.Body()
	lb := &hclwrite.Block{}
	if len(body.Blocks()) == 0 {
		lb = body.AppendNewBlock("resource", []string{x.resourceType, x.resourceLabel})
	} else {
		lb = body.Blocks()[0]
	}
	lbBody := lb.Body()
	lbBody.SetAttributeValue(key, toCtyVal(value))
	return x
}

func (x *resourceBlock) AddStatement(key string, value string) *resourceBlock {
	body := x.Body()
	lb := &hclwrite.Block{}
	if len(body.Blocks()) == 0 {
		lb = body.AppendNewBlock("resource", []string{x.resourceType, x.resourceLabel})
	} else {
		lb = body.Blocks()[0]
	}
	lbBody := lb.Body()
	lbBody.SetAttributeTraversal(key, hcl.Traversal{hcl.TraverseRoot{Name: value}})
	return x
}

func (x *resourceBlock) KeyRef(k string) string {
	for _, key := range x.keys {
		if key == k {
			return "local." + key
		}
	}
	fmt.Println("could not find key", k)
	os.Exit(1)
	return ""
}

func (x *resourceBlock) Type() string {
	return "resource"
}

func (x *resourceBlock) String() string {
	return string(x.Bytes())
}

package blocks

import (
	"fmt"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type resourceBlock struct {
	*hclwrite.File
	keys          []string
	resourceType  string
	resourceLabel string
}

func NewResource(resourceType, resourceLabel string) *resourceBlock {
	file := hclwrite.NewEmptyFile()
	body := file.Body()
	body.AppendNewBlock("resource", []string{resourceType, resourceLabel})
	return &resourceBlock{file, nil, resourceType, resourceLabel}
}

func (x *resourceBlock) WithImport(id string) *resourceBlock {
	importBlockBody := x.Body().AppendNewBlock("import", nil).Body()
	importBlockBody.SetAttributeValue("id", cty.StringVal(id))
	importBlockBody.SetAttributeTraversal("to", hcl.Traversal{hcl.TraverseRoot{Name: x.resourceType}, hcl.TraverseAttr{Name: x.resourceLabel}})
	return x
}

func (x *resourceBlock) Add(key string, value any) *resourceBlock {
	x.keys = append(x.keys, key)
	body := x.Body()
	lb := body.Blocks()[0]
	lbBody := lb.Body()
	lbBody.SetAttributeValue(key, toCtyVal(value))
	return x
}

func (x *resourceBlock) AddStatement(key string, value string) *resourceBlock {
	body := x.Body()
	lb := body.Blocks()[0]
	lbBody := lb.Body()
	lbBody.SetAttributeTraversal(key, hcl.Traversal{hcl.TraverseRoot{Name: value}})
	return x
}

func (x *resourceBlock) ForEach(m string) *resourceBlock {
	x.AddStatement("for_each", m)
	return x
}

func (x *resourceBlock) Self() string {
	return fmt.Sprintf("%s.%s", x.resourceType, x.resourceLabel)
}

func (x *resourceBlock) Type() string {
	return "resource"
}

func (x *resourceBlock) String() string {
	return string(x.Bytes())
}

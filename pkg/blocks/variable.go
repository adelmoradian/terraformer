package blocks

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type variableBlock struct {
	name            string
	description     string
	varType         string
	validationBlock validation
	isSensitive     bool
	isNullable      bool
	isEphemeral     bool
	isConstant      bool
}

type validation struct {
	condition, errMessage string
}

func (x *variableBlock) Type() string {
	return "variable"
}

func NewVarBlock(name, description string) *variableBlock {
	return &variableBlock{name: name, description: description}
}

func (x *variableBlock) SetType(t string) *variableBlock {
	x.varType = t
	return x
}

func (x *variableBlock) SetValidation(condition, errMessage string) *variableBlock {
	x.validationBlock.condition = condition
	x.validationBlock.errMessage = errMessage
	return x
}

func (x *variableBlock) IsSensitive() *variableBlock {
	x.isSensitive = true
	return x
}

func (x *variableBlock) IsNullable() *variableBlock {
	x.isNullable = true
	return x
}

func (x *variableBlock) IsEphemeral() *variableBlock {
	x.isEphemeral = true
	return x
}

func (x *variableBlock) IsConstant() *variableBlock {
	x.isConstant = true
	return x
}

func (x *variableBlock) String() string {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	varBlock := rootBody.AppendNewBlock(x.Type(), []string{x.name})
	varBlockBody := varBlock.Body()
	varBlockBody.SetAttributeTraversal("type", hcl.Traversal{hcl.TraverseRoot{Name: x.varType}})
	varBlockBody.AppendNewline()
	varBlockBody.SetAttributeValue("description", cty.StringVal(x.description))
	varBlockBody.SetAttributeValue("sensitive", cty.BoolVal(x.isSensitive))
	varBlockBody.SetAttributeValue("nullable", cty.BoolVal(x.isNullable))
	varBlockBody.SetAttributeValue("ephemeral", cty.BoolVal(x.isEphemeral))
	varBlockBody.SetAttributeValue("const", cty.BoolVal(x.isConstant))
	if x.validationBlock.condition != "" {
		varBlockBody.AppendNewline()
		validationBlock := varBlockBody.AppendNewBlock("validation", nil)
		validationBlockBody := validationBlock.Body()
		validationBlockBody.SetAttributeTraversal("condition", hcl.Traversal{hcl.TraverseRoot{Name: x.validationBlock.condition}})
		validationBlockBody.SetAttributeValue("error_message", cty.StringVal(x.validationBlock.errMessage))
	}
	return string(hclwrite.Format(f.Bytes()))
}

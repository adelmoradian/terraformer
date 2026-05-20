package blocks

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type localsBlock struct {
	*hclwrite.File
	keys []string
}

func NewLocals() *localsBlock {
	return &localsBlock{hclwrite.NewEmptyFile(), nil}
}

func toCtyVal(x any) cty.Value {
	switch v := x.(type) {
	case string:
		return cty.StringVal(v)
	case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64:
		i, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
		if err != nil {
			fmt.Println("could not convert number", err.Error())
			os.Exit(1)
		}
		return cty.NumberFloatVal(i)
	case bool:
		return cty.BoolVal(v)
	case []any:
		vList := make([]cty.Value, 0, len(v))
		for _, vAny := range v {
			vList = append(vList, toCtyVal(vAny))
		}
		return cty.ListVal(vList)
	case map[string]any:
		vMap := make(map[string]cty.Value)
		for k, vAny := range v {
			vMap[k] = toCtyVal(vAny)
		}
		return cty.MapVal(vMap)
	default:
		fmt.Println("could not parse value", v)
		os.Exit(1)
	}
	return cty.NullVal(cty.NilType)
}

func (x *localsBlock) Add(key string, value any) *localsBlock {
	x.keys = append(x.keys, key)
	body := x.Body()
	lb := &hclwrite.Block{}
	if len(body.Blocks()) == 0 {
		lb = body.AppendNewBlock("locals", nil)
	} else {
		lb = body.Blocks()[0]
	}
	lbBody := lb.Body()
	lbBody.SetAttributeValue(key, toCtyVal(value))
	return x
}

func (x *localsBlock) AddStatement(key string, value string) *localsBlock {
	body := x.Body()
	lb := &hclwrite.Block{}
	if len(body.Blocks()) == 0 {
		lb = body.AppendNewBlock("locals", nil)
	} else {
		lb = body.Blocks()[0]
	}
	lbBody := lb.Body()
	lbBody.SetAttributeTraversal(key, hcl.Traversal{hcl.TraverseRoot{Name: value}})
	return x
}

func (x *localsBlock) KeyRef(k string) string {
	for _, key := range x.keys {
		if key == k {
			return "local." + key
		}
	}
	fmt.Println("could not find key", k)
	os.Exit(1)
	return ""
}

func (x *localsBlock) Type() string {
	return "locals"
}

func (x *localsBlock) String() string {
	return string(x.Bytes())
}

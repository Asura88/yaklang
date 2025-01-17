package statements

import "github.com/yaklang/yaklang/common/javaclassparser/decompiler/core/class_context"

type CustomStatement struct {
	Name       string
	Info       any
	StringFunc func(funcCtx *class_context.ClassContext) string
}

func (v *CustomStatement) String(funcCtx *class_context.ClassContext) string {
	return v.StringFunc(funcCtx)
}
func NewCustomStatement(stringFun func(funcCtx *class_context.ClassContext) string) *CustomStatement {
	return &CustomStatement{
		StringFunc: stringFun,
	}
}

package ssa

import (
	"github.com/yaklang/yaklang/common/log"
	"github.com/yaklang/yaklang/common/utils"
)

func (pkg *Program) GetBluePrint(name string) *Blueprint {
	if pkg == nil {
		return nil
	}
	return pkg.GetClassBlueprintEx(name, "")
}

func (b *FunctionBuilder) GetBluePrint(name string) *Blueprint {
	p := b.prog
	return p.GetBluePrint(name)
}

func (b *FunctionBuilder) SetClassBluePrint(name string, class *Blueprint) {
	p := b.prog
	if _, ok := p.ClassBluePrint[name]; ok {
		log.Errorf("SetClassBluePrint: this class redeclare")
	}
	p.ClassBluePrint[name] = class
}

// CreateClassBluePrint will create object template (maybe class)
// in dynamic and classless language, we can create object without class
// but because of the 'this/super', we will still keep the concept 'Class'
// for ref the method/function, the blueprint is a container too,
// saving the static variables and util methods.

func (b *FunctionBuilder) CreateBluePrintWithPkgName(name string, tokenizer ...CanStartStopToken) *Blueprint {
	prog := b.prog
	blueprint := NewClassBluePrint(name)
	if prog.ClassBluePrint == nil {
		prog.ClassBluePrint = make(map[string]*Blueprint)
	}
	blueprint.GeneralUndefined = func(s string) *Undefined {
		return b.EmitUndefined(s)
	}
	prog.ClassBluePrint[name] = blueprint
	klassvar := b.CreateVariable(name, tokenizer...)
	klassContainer := b.EmitEmptyContainer()
	b.AssignVariable(klassvar, klassContainer)
	if err := blueprint.InitializeWithContainer(klassContainer); err != nil {
		log.Errorf("CreateClassBluePrint.InitializeWithContainer error: %s", err)
	}
	return blueprint
}
func (b *FunctionBuilder) CreateBluePrint(name string, tokenizer ...CanStartStopToken) *Blueprint {
	return b.CreateBluePrintWithPkgName(name, tokenizer...)
}

// ReadSelfMember  用于读取当前类成员，包括静态成员和普通成员和方法。
// 其中使用MarkedThisClassBlueprint标识当前在哪个类中。
func (b *FunctionBuilder) ReadSelfMember(name string) Value {
	var value Value
	defer func() {
		if !utils.IsNil(value) {
			b.AssignVariable(b.CreateVariable(name), value)
		}
	}()
	if class := b.MarkedThisClassBlueprint; class != nil {
		variable := b.GetStaticMember(class, name)
		if _value := b.PeekValueByVariable(variable); _value != nil {
			return _value
		}
		if val := class.GetStaticMember(name); !utils.IsNil(val) {
			return val
		}
		if normalMember := class.GetNormalMember(name); !utils.IsNil(normalMember) {
			return normalMember
		}
		if method_ := class.GetNormalMethod(name); !utils.IsNil(method_) {
			return method_
		}
	}
	return nil
}
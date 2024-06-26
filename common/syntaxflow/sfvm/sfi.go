package sfvm

import (
	"fmt"
	"strconv"
)

type SFVMOpCode int

const (
	OpPass SFVMOpCode = iota

	// enter/exit statement
	OpEnterStatement
	OpExitStatement

	// push input
	OpPushInput
	// duplicate the top of stack
	OpDuplicate

	// OpPushSearchExact can push data from origin
	OpPushSearchExact
	OpPushSearchGlob
	OpPushSearchRegexp

	// handle function call
	opGetCall
	OpGetCallArgs
	OpGetAllCallArgs

	// use def chain
	OpGetUsers
	OpGetBottomUsers
	OpGetDefs
	OpGetTopDefs

	// ListOperation
	OpListIndex

	// => variable
	OpNewRef
	OpUpdateRef

	// OpPushNumber and OpPushString and OpPushBool can push literal into stack
	OpPushNumber
	OpPushString
	OpPushBool
	OpPop

	// Condition
	// use the []bool  && []Value of stack top, push result into stack
	OpCondition
	OpCompareOpcode
	OpCompareString

	/*
		Binary Operator
		Fetch TWO in STACK, calc result, push result into stack
	*/
	OpEq
	OpNotEq
	OpGt
	OpGtEq
	OpLt
	OpLtEq
	OpLogicAnd
	OpLogicOr
	OpLogicBang

	/*
		Unary Operator: Fetch ONE in STACK, calc result, push result into stack
	*/
	OpReMatch
	OpGlobMatch
	OpNot

	// OpCheckParams check the params in vm context
	// if not match, record error
	// matched, use 'then expr' (if exists)
	OpCheckParams

	// OpAddDescription add description to current context
	OpAddDescription
)

type SFI struct {
	OpCode           SFVMOpCode
	UnaryInt         int
	UnaryStr         string
	Desc             string
	Values           []string
	SyntaxFlowConfig []*RecursiveConfigItem
}

func (s *SFI) ValueByIndex(i int) string {
	if i < 0 || i >= len(s.Values) {
		return ""
	}
	return s.Values[i]
}

const verboseLen = "%-12s"

func (s *SFI) String() string {
	switch s.OpCode {
	case OpEnterStatement:
		return "- enter -"
	case OpExitStatement:
		return "- exit -"
	case OpPass:
		return "- pass -"
	case OpPushBool:
		v := "false"
		if s.UnaryInt > 0 {
			v = "true"
		}
		return fmt.Sprintf(verboseLen+" %v", "push", v)
	case OpPushString:
		return fmt.Sprintf(verboseLen+" (len:%v) %v", "push", len(s.UnaryStr), strconv.Quote(s.UnaryStr))
	case OpPushNumber:
		return fmt.Sprintf(verboseLen+" %v", "push", s.UnaryInt)
	case OpDuplicate:
		return fmt.Sprintf(verboseLen+" %v", "duplicate", s.UnaryStr)
	case OpPushInput:
		return fmt.Sprintf(verboseLen+" %v", "push$input", s.UnaryStr)

	case OpPushSearchGlob:
		return fmt.Sprintf(verboseLen+" %v isMember[%v]", "push$glob", s.UnaryStr, s.UnaryInt)
	case OpPushSearchExact:
		return fmt.Sprintf(verboseLen+" %v isMember[%v]", "push$exact", s.UnaryStr, s.UnaryInt)
	case OpPushSearchRegexp:
		return fmt.Sprintf(verboseLen+" %v isMember[%v]", "push$regexp", s.UnaryStr, s.UnaryInt)
	case opGetCall:
		return fmt.Sprintf(verboseLen+" %v", "getCall", s.UnaryStr)
	case OpGetAllCallArgs:
		return fmt.Sprintf(verboseLen+" %v", "getAllCallArgs", s.UnaryStr)
	case OpGetCallArgs:
		return fmt.Sprintf(verboseLen+" %v", "getCallArgs", s.UnaryInt)
	case OpGetUsers:
		return fmt.Sprintf(verboseLen+" %v", "users", s.UnaryStr)
	case OpGetDefs:
		return fmt.Sprintf(verboseLen+" %v", "defs", s.UnaryStr)
	case OpGetTopDefs:
		return fmt.Sprintf(verboseLen+" %v", "topDefs", s.UnaryStr)
	case OpGetBottomUsers:
		return fmt.Sprintf(verboseLen+" %v", "bottomUse", s.UnaryStr)
	case OpListIndex:
		return fmt.Sprintf(verboseLen+" %v", "listIndex", s.UnaryStr)
	case OpNewRef:
		return fmt.Sprintf(verboseLen+" %v", "new$ref", s.UnaryStr)
	case OpUpdateRef:
		return fmt.Sprintf(verboseLen+" %v", "update$ref", s.UnaryStr)
	case OpCompareOpcode:
		return fmt.Sprintf(verboseLen+" %v", "compare opcode", s.Values)
	case OpCompareString:
		return fmt.Sprintf(verboseLen+" %v", "compare opcode", s.Values)
	case OpCondition:
		return fmt.Sprintf(verboseLen+" %v", "condition", s.UnaryStr)
	case OpEq:
		return fmt.Sprintf(verboseLen+" %v", "(operator) ==", s.UnaryStr)
	case OpNotEq:
		return fmt.Sprintf(verboseLen+" %v", "(operator) !=", s.UnaryStr)
	case OpGt:
		return fmt.Sprintf(verboseLen+" %v", "(operator) >", s.UnaryStr)
	case OpGtEq:
		return fmt.Sprintf(verboseLen+" %v", "(operator) >=", s.UnaryStr)
	case OpLt:
		return fmt.Sprintf(verboseLen+" %v", "(operator) <", s.UnaryStr)
	case OpLtEq:
		return fmt.Sprintf(verboseLen+" %v", "(operator) <=", s.UnaryStr)
	case OpReMatch:
		return fmt.Sprintf(verboseLen+" %v", "(operator) ~=", s.UnaryStr)
	case OpGlobMatch:
		return fmt.Sprintf(verboseLen+" %v", "(operator) *~", s.UnaryStr)
	case OpNot:
		return fmt.Sprintf(verboseLen+" %v", "(operator) !", s.UnaryStr)
	case OpLogicAnd:
		return fmt.Sprintf(verboseLen+" %v", "(operator) &&", s.UnaryStr)
	case OpLogicOr:
		return fmt.Sprintf(verboseLen+" %v", "(operator) ||", s.UnaryStr)
	case OpLogicBang:
		return fmt.Sprintf(verboseLen+" %v", "(operator) !", s.UnaryStr)
	case OpPop:
		return fmt.Sprintf(verboseLen+" %v", "pop", s.UnaryStr)
	case OpCheckParams:
		var suffix string
		if ret := s.ValueByIndex(0); ret != "" {
			suffix += " then: " + ret
		}
		if ret := s.ValueByIndex(1); ret != "" {
			suffix += ", else: " + ret
		}
		return fmt.Sprintf(verboseLen+" $%v"+suffix, "check", s.UnaryStr)
	case OpAddDescription:
		var suffix string
		if ret := s.ValueByIndex(0); ret != "" {
			suffix += " value: " + ret
		}
		return fmt.Sprintf(verboseLen+" %v"+suffix, "desc", s.UnaryStr)
	default:
		panic("unhandled default case")
	}
	return ""
}

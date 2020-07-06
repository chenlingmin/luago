package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	. "luago/api"
	. "luago/binchunk"
	. "luago/compiler/lexer"
	"luago/compiler/parser"
	. "luago/vm"
	"os"
)

func main() {

	if len(os.Args) > 1 {
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		testLexer(string(data), os.Args[1])
		testParser(string(data), os.Args[1])
	}

	//
	//if len(os.Args) > 1 {
	//	data, err := ioutil.ReadFile(os.Args[1])
	//	if err != nil {
	//		panic(err)
	//	}
	//	ls := state.New()
	//	ls.Register("print", print)
	//	ls.Register("fail", fail)
	//	ls.Register("getmetatable", getMetatable)
	//	ls.Register("setmetatable", setMetatable)
	//	ls.Register("next", next)
	//	ls.Register("pairs", pairs)
	//	ls.Register("ipairs", iPairs)
	//	ls.Register("error", error)
	//	ls.Register("pcall", pCall)
	//	ls.Load(data, os.Args[1], "b")
	//	ls.Call(0, 0)
	//}

	//if len(os.Args) > 1 {
	//	data, err := ioutil.ReadFile(os.Args[1])
	//	if err != nil {
	//		panic(err)
	//	}
	//	proto := Undump(data)
	//	luaMain(proto)
	//}
	//}
	//ls := state.New()
	//ls.PushInteger(1)
	//ls.PushString("2.0")
	//ls.PushString("3.0")
	//ls.PushString("4.0")
	//printStack(ls)
	//
	//ls.Arith(LUA_OPADD)
	//printStack(ls)
	//ls.Arith(LUA_OPBNOT)
	//printStack(ls)
	//ls.Len(2)
	//printStack(ls)
	//ls.Concat(3)
	//printStack(ls)
	//ls.PushBoolean(ls.Compare(1, 2, LUA_OPEQ))
	//printStack(ls)

	//
	//ls := state.New()
	//
	//ls.PushBoolean(true)
	//printStack(ls)
	//ls.PushInteger(10)
	//printStack(ls)
	//ls.PushNil()
	//printStack(ls)
	//ls.PushString("hello")
	//printStack(ls)
	//ls.PushValue(-4)
	//printStack(ls)
	//ls.Replace(3)
	//printStack(ls)
	//ls.SetTop(6)
	//printStack(ls)
	//ls.Remove(-3)
	//printStack(ls)
	//ls.SetTop(-5)
	//printStack(ls)
}

func testParser(chunk, chunkName string) {
	ast := parser.Parse(chunk, chunkName)
	b, err := json.Marshal(ast)
	if err != nil {
		panic(err)
	}
	println(string(b))
}

func testLexer(chunk, chunkName string) {
	lexer := NewLexer(chunk, chunkName)
	for {
		line, kind, token := lexer.NextToken()
		fmt.Printf("[%2d] [%-10s] %s\n", line, kindToCategory(kind), token)
		if kind == TOKEN_EOF {
			break
		}
	}
}

func kindToCategory(kind int) string {
	switch {
	case kind < TOKEN_SEP_SEMI:
		return "other"
	case kind <= TOKEN_SEP_RCURLY:
		return "separator"
	case kind <= TOKEN_OP_NOT:
		return "operator"
	case kind <= TOKEN_KW_WHILE:
		return "keyword"
	case kind == TOKEN_IDENTIFIER:
		return "identifier"
	case kind == TOKEN_NUMBER:
		return "number"
	case kind == TOKEN_STRING:
		return "string"
	default:
		return "other"
	}
}

//
//func luaMain(proto *Prototype) {
//	nRegs := int(proto.MaxStackSize)
//	ls := state.New(nRegs+8, proto)
//	ls.SetTop(nRegs)
//	for {
//		pc := ls.PC()
//		inst := Instruction(ls.Fetch())
//		if inst.Opcode() != OP_RETURN {
//			inst.Execute(ls)
//			fmt.Printf("[%20d] %s ", pc+1, inst.OpName())
//			printStack(ls)
//		} else {
//			break
//		}
//	}
//}

func error(ls LuaState) int {
	return ls.Error()
}

func pCall(ls LuaState) int {
	nArgs := ls.GetTop() - 1
	status := ls.PCall(nArgs, -1, 0)
	ls.PushBoolean(status == LUA_OK)
	ls.Insert(1)
	return ls.GetTop()
}

func iPairs(ls LuaState) int {
	ls.PushGoFunction(_isPairsAux)
	ls.PushValue(1)
	ls.PushInteger(0)
	return 3
}
func _isPairsAux(ls LuaState) int {
	i := ls.ToInteger(2) + 1
	ls.PushInteger(i)
	if ls.GetI(1, i) == LUA_TNIL {
		return 1
	} else {
		return 2
	}
}

func pairs(ls LuaState) int {
	ls.PushGoFunction(next)
	ls.PushValue(1)
	ls.PushNil()
	return 3
}

func next(ls LuaState) int {
	ls.SetTop(2) //
	if ls.Next(1) {
		return 2
	} else {
		ls.PushNil()
		return 1
	}
}

func getMetatable(ls LuaState) int {
	if !ls.GetMetatable(1) {
		ls.PushNil()
	}
	return 1
}

func setMetatable(ls LuaState) int {
	ls.SetMetatable(1)
	return 1
}

func fail(ls LuaState) int {
	return 0
}

func print(ls LuaState) int {
	nArgs := ls.GetTop()
	for i := 1; i <= nArgs; i++ {
		if ls.IsBoolean(i) {
			fmt.Printf("%t", ls.ToBoolean(i))
		} else if ls.IsString(i) {
			fmt.Print(ls.ToString(i))
		} else {
			fmt.Print(ls.TypeName(ls.Type(i)))
		}
		if i < nArgs {
			fmt.Print("\t")
		}
	}
	fmt.Println()
	return 0
}

func printStack(ls LuaState) {
	top := ls.GetTop()
	for i := 1; i <= top; i++ {
		t := ls.Type(i)
		switch t {
		case LUA_TBOOLEAN:
			fmt.Printf("[%t]", ls.ToBoolean(i))
		case LUA_TNUMBER:
			fmt.Printf("[%g]", ls.ToNumber(i))
		case LUA_TSTRING:
			fmt.Printf("[%q]", ls.ToString(i))
		default: // other values
			fmt.Printf("[%s]", ls.TypeName(t))
		}
	}
	fmt.Println()
}

func list(f *Prototype) {
	printHeader(f)
	printCode(f)
	printDetail(f)
	for _, p := range f.Protos {
		list(p)
	}
}

func printHeader(f *Prototype) {
	funcType := "main"
	if f.LineDefined > 0 {
		funcType = "function"
	}

	varargFlag := ""
	if f.IsVararg > 0 {
		varargFlag = "+"
	}

	fmt.Printf("\n%s <%s:%d,%d> (%d instructions)\n",
		funcType, f.Source, f.LineDefined, f.LastLineDefined, len(f.Code))

	fmt.Printf("%d%s params, %d slots, %d upvalues, ",
		f.NumParams, varargFlag, f.MaxStackSize, len(f.Upvalues))

	fmt.Printf("%d locals, %d constants, %d functions\n",
		len(f.LocVars), len(f.Constants), len(f.Protos))
}

func printCode(f *Prototype) {
	for pc, c := range f.Code {
		line := "-"
		if len(f.LineInfo) > 0 {
			line = fmt.Sprintf("%d", f.LineInfo[pc])
		}

		i := Instruction(c)

		fmt.Printf("\t%d\t[%s]\t%s \t", pc+1, line, i.OpName())
		printOperands(i)
		fmt.Printf("\n")
	}
}

func printOperands(i Instruction) {
	switch i.OpMode() {
	case IABC:
		a, b, c := i.ABC()
		fmt.Printf("%d", a)
		if i.BMode() != OpArgN {
			if b > 0xFF {
				fmt.Printf(" %d", -1-b&0xFF)
			} else {
				fmt.Printf(" %d", b)
			}
		}
		if i.CMode() != OpArgN {
			if c > 0xFF {
				fmt.Printf(" %d", -1-c&0xFF)
			} else {
				fmt.Printf(" %d", c)
			}
		}
	case IABx:
		a, bx := i.ABx()

		fmt.Printf("%d", a)
		if i.BMode() == OpArgK {
			fmt.Printf(" %d", -1-bx)
		} else if i.BMode() == OpArgU {
			fmt.Printf(" %d", bx)
		}
	case IAsBx:
		a, sbx := i.AsBx()
		fmt.Printf("%d %d", a, sbx)
	case IAx:
		ax := i.Ax()
		fmt.Printf("%d", -1-ax)
	}

}

func printDetail(f *Prototype) {
	fmt.Printf("constants (%d):\n", len(f.Constants))
	for i, k := range f.Constants {
		fmt.Printf("\t%d\t%s\n", i+1, constantToString(k))
	}

	fmt.Printf("locals (%d):\n", len(f.LocVars))
	for i, locVar := range f.LocVars {
		fmt.Printf("\t%d\t%s\t%d\t%d\n",
			i, locVar.VarName, locVar.StartPC+1, locVar.EndPC+1)
	}

	fmt.Printf("upvalues (%d):\n", len(f.Upvalues))
	for i, upval := range f.Upvalues {
		fmt.Printf("\t%d\t%s\t%d\t%d\n",
			i, upvalName(f, i), upval.Instack, upval.Idx)
	}
}

func constantToString(k interface{}) string {
	switch k.(type) {
	case nil:
		return "nil"
	case bool:
		return fmt.Sprintf("%t", k)
	case float64:
		return fmt.Sprintf("%g", k)
	case int64:
		return fmt.Sprintf("%d", k)
	case string:
		return fmt.Sprintf("%q", k)
	default:
		return "?"
	}
}

func upvalName(f *Prototype, idx int) string {
	if len(f.UpvalueNames) > 0 {
		return f.UpvalueNames[idx]
	}
	return "-"
}

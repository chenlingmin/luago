package state

import (
	"fmt"
	"luago/binchunk"
	"luago/vm"
)

func (self *luaState) Load(chunk []byte, chunkName, mode string) int {
	proto := binchunk.Undump(chunk)
	c := newLuaClosure(proto)
	self.stack.push(c)
	return 0
}

func (self *luaState) Call(nArgs, nResults int) {
	val := self.stack.get(-(nArgs + 1))
	if c, ok := val.(*closure); ok {
		if c.proto != nil {
			fmt.Printf("call %s<%d, %d>\n", c.proto.Source, c.proto.LineDefined, c.proto.LastLineDefined)
			self.callLuaClosure(nArgs, nResults, c)
		} else { // 调用 go 函数
			self.callGoClosure(nArgs, nResults, c)
		}
	} else {
		panic("not function! ")
	}
}

func (self *luaState) callGoClosure(nArgs, nResults int, c *closure) {
	newStack := newLuaStack(nArgs+20, self)
	newStack.closure = c

	args := self.stack.popN(nArgs)
	newStack.pushN(args, nArgs)
	self.stack.pop()

	self.pushLuaStack(newStack)
	r := c.goFunc(self)
	self.popLuaStack()

	if nResults != 0 {
		results := newStack.popN(r)
		self.stack.check(len(results))
		self.stack.pushN(results, nResults)
	}

}

func (self *luaState) callLuaClosure(nArgs, nResults int, c *closure) {
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1

	// create new lua stack
	newStack := newLuaStack(nRegs+20, self)
	newStack.closure = c

	// pass args, pop func
	funcAndArgs := self.stack.popN(nArgs + 1)
	newStack.pushN(funcAndArgs[1:], nParams)
	newStack.top = nRegs
	if nArgs > nParams && isVararg {
		newStack.varargs = funcAndArgs[nParams+1:]
	}

	// run closure
	self.pushLuaStack(newStack)
	self.runLuaClosure()
	self.popLuaStack()

	// return results
	if nResults != 0 {
		results := newStack.popN(newStack.top - nRegs)
		self.stack.check(len(results))
		self.stack.pushN(results, nResults)
	}
}

func (self *luaState) runLuaClosure() {
	for {
		inst := vm.Instruction(self.Fetch())
		inst.Execute(self)
		if inst.Opcode() == vm.OP_RETURN {
			break
		}
	}
}

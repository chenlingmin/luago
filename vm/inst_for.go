package vm

import "luago/api"

// R(A)-=R(A+2); pc+=sBx
func forPrep(i Instruction, vm api.LuaVM) {
	a, sBx := i.AsBx()
	a += 1
	if vm.Type(a) == api.LUA_TSTRING {
		vm.PushNumber(vm.ToNumber(a))
		vm.Replace(a)
	}
	if vm.Type(a+1) == api.LUA_TSTRING {
		vm.PushNumber(vm.ToNumber(a + 1))
		vm.Replace(a + 1)
	}
	if vm.Type(a+2) == api.LUA_TSTRING {
		vm.PushNumber(vm.ToNumber(a + 2))
		vm.Replace(a + 2)
	}
	vm.PushValue(a)
	vm.PushValue(a + 2)
	vm.Arith(api.LUA_OPSUB)
	vm.Replace(a)
	vm.AddPC(sBx)
}

// R(A)+=R(A+2);
// if R(A) <?= R(A+1) then {
//   pc+=sBx; R(A+3)=R(A)
// }
func forLoop(i Instruction, vm api.LuaVM) {
	a, sBx := i.AsBx()
	a += 1

	// R(A)+=R(A+2);
	vm.PushValue(a + 2)
	vm.PushValue(a)
	vm.Arith(api.LUA_OPADD)
	vm.Replace(a)

	isPositiveStep := vm.ToNumber(a+2) >= 0
	if isPositiveStep && vm.Compare(a, a+1, api.LUA_OPLE) ||
		!isPositiveStep && vm.Compare(a+1, a, api.LUA_OPLE) {

		// pc+=sBx; R(A+3)=R(A)
		vm.AddPC(sBx)
		vm.Copy(a, a+3)
	}
}

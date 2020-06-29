package state

import (
	"fmt"
	"luago/api"
)

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_typename
func (self *luaState) TypeName(tp api.LuaType) string {
	switch tp {
	case api.LUA_TNONE:
		return "no value"
	case api.LUA_TNIL:
		return "nil"
	case api.LUA_TBOOLEAN:
		return "boolean"
	case api.LUA_TNUMBER:
		return "number"
	case api.LUA_TSTRING:
		return "string"
	case api.LUA_TTABLE:
		return "table"
	case api.LUA_TFUNCTION:
		return "function"
	case api.LUA_TTHREAD:
		return "thread"
	default:
		return "userdata"
	}
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_type
func (self *luaState) Type(idx int) api.LuaType {
	if self.stack.isValid(idx) {
		val := self.stack.get(idx)
		return typeOf(val)
	}
	return api.LUA_TNONE
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isnone
func (self *luaState) IsNone(idx int) bool {
	return self.Type(idx) == api.LUA_TNONE
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isnil
func (self *luaState) IsNil(idx int) bool {
	return self.Type(idx) == api.LUA_TNIL
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isnoneornil
func (self *luaState) IsNoneOrNil(idx int) bool {
	return self.Type(idx) <= api.LUA_TNIL
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isboolean
func (self *luaState) IsBoolean(idx int) bool {
	return self.Type(idx) == api.LUA_TBOOLEAN
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_istable
func (self *luaState) IsTable(idx int) bool {
	return self.Type(idx) == api.LUA_TTABLE
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isfunction
func (self *luaState) IsFunction(idx int) bool {
	return self.Type(idx) == api.LUA_TFUNCTION
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isthread
func (self *luaState) IsThread(idx int) bool {
	return self.Type(idx) == api.LUA_TTHREAD
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isstring
func (self *luaState) IsString(idx int) bool {
	t := self.Type(idx)
	return t == api.LUA_TSTRING || t == api.LUA_TNUMBER
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isnumber
func (self *luaState) IsNumber(idx int) bool {
	_, ok := self.ToNumberX(idx)
	return ok
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isinteger
func (self *luaState) IsInteger(idx int) bool {
	val := self.stack.get(idx)
	_, ok := val.(int64)
	return ok
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_iscfunction
func (self *luaState) IsGoFunction(idx int) bool {
	val := self.stack.get(idx)
	if c, ok := val.(*closure); ok {
		return c.goFunc != nil
	}
	return false
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_toboolean
func (self *luaState) ToBoolean(idx int) bool {
	val := self.stack.get(idx)
	return convertToBoolean(val)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tointeger
func (self *luaState) ToInteger(idx int) int64 {
	i, _ := self.ToIntegerX(idx)
	return i
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tointegerx
func (self *luaState) ToIntegerX(idx int) (int64, bool) {
	val := self.stack.get(idx)
	return convertToInteger(val)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tonumber
func (self *luaState) ToNumber(idx int) float64 {
	n, _ := self.ToNumberX(idx)
	return n
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tonumberx
func (self *luaState) ToNumberX(idx int) (float64, bool) {
	val := self.stack.get(idx)
	return convertToFloat(val)
}

// [-0, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_tostring
func (self *luaState) ToString(idx int) string {
	s, _ := self.ToStringX(idx)
	return s
}

func (self *luaState) ToStringX(idx int) (string, bool) {
	val := self.stack.get(idx)

	switch x := val.(type) {
	case string:
		return x, true
	case int64, float64:
		s := fmt.Sprintf("%v", x) // todo
		self.stack.set(idx, s)
		return s, true
	default:
		return "", false
	}
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tocfunction
func (self *luaState) ToGoFunction(idx int) api.GoFunction {
	val := self.stack.get(idx)
	if c, ok := val.(*closure); ok {
		return c.goFunc
	}
	return nil
}

package binchunk

const (
	LUA_SIGNATURE    = "\x1bLua"
	LUAC_VERSION     = 0x53
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	CINT_SIZE        = 4
	CSIZET_SIZE      = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

const (
	TAG_NIL       = 0x00
	TAG_BOOLEAN   = 0x01
	TAG_NUMBER    = 0x03
	TAG_INTEGER   = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

type binaryChunk struct {
	header                  //头部
	sizeUpvalues byte       // 主函数 upvalue 数量
	mainFunc     *Prototype // 主函数原型
}

type header struct {
	signature       [4]byte // 0x1B4C7561
	version         byte    // 0x53
	format          byte
	luacData        [6]byte
	cintSize        byte
	sizetSize       byte
	instructionSize byte
	luaIntegerSize  byte
	luaNumberSize   byte
	luacInt         int64
	luacNum         float64
}

type Prototype struct {
	Source          string        // 源文件名
	LineDefined     uint32        // 起止行号
	LastLineDefined uint32        // 起止行号
	NumParams       byte          // 固定参数个数
	IsVararg        byte          // 是否是 Vararg 函数
	MaxStackSize    byte          // 寄存器数量, 这个字段也被叫作MaxStackSize，为什么这样叫呢？这是因为Lua虚拟机在执行函数时，真正使用的其实是一种栈结构，这种栈结构除了可以进行常规地推入和弹出操作以外，还可以按索引访问，所以可以用来模拟寄存器
	Code            []uint32      // 指令表
	Constants       []interface{} // 常量表
	Upvalues        []Upvalue
	Protos          []*Prototype // 子函数原型
	LineInfo        []uint32     // 行号表
	LocVars         []LocVar     // 局部变量表
	UpvalueNames    []string     // Upvalue名列表
}

type Upvalue struct {
	Instack byte
	Idx     byte
}

type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}

func Undump(data []byte) *Prototype {
	reader := &reader{data}
	reader.checkHeader()        // 检验头部
	reader.readByte()           // 跳过 Upvalue 数量
	return reader.readProto("") // 读取函数原型
}

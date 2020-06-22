package vm

/*
** 这种编码方式和浮点数的编码方式类似，只是仅用一个字节。
** 具体来说，如果把某个字节用二进制写成eeeeexxx，那么当eeeee == 0时该字节表示的整数就是xxx，
** 否则该字节表示的整数是(1xxx) ＊ 2^(eeeee -1)。
** Floating Point Byte 编码
** converts an integer to a "floating point byte", represented as
** (eeeeexxx), where the real value is (1xxx) * 2^(eeeee - 1) if
** eeeee != 0 and (xxx) otherwise.
 */
func Int2fb(x int) int {
	e := 0 /* exponent */
	if x < 8 {
		return x
	}
	for x >= (8 << 4) {
		x = (x + 0xf) >> 4
		e += 4
	}
	for x >= (8 << 1) {
		x = (x + 1) >> 1
		e++
	}
	return ((e + 1) << 3) | (x - 8)
}

/* converts back*/
func Fb2int(x int) int {
	if x < 8 {
		return x
	} else {
		return ((x & 7) + 8) << uint((x>>3)-1)
	}
}

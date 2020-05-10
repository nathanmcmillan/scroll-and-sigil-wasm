package fast

import "unsafe"

// ByteWriter struct
type ByteWriter struct {
	buffer []byte
	dex    int
}

// ByteWriterInit func
func ByteWriterInit(estimate int) *ByteWriter {
	writer := &ByteWriter{}
	writer.buffer = make([]byte, estimate)
	return writer
}

func (me *ByteWriter) verify(x int) {
	size := len(me.buffer)
	if me.dex+x >= size {
		temp := make([]byte, size+x+64)
		copy(temp, me.buffer)
		me.buffer = temp
	}
}

// PutBytes func
func (me *ByteWriter) PutBytes(x []byte) {
	size := len(x)
	me.verify(size)
	pos := me.dex
	for i := 0; i < size; i++ {
		me.buffer[pos] = x[i]
		pos++
	}
	me.dex = pos
}

// PutByte func
func (me *ByteWriter) PutByte(x byte) {
	me.verify(1)
	me.buffer[me.dex] = x
	me.dex++
}

// PutInt8 func
func (me *ByteWriter) PutInt8(x int8) {
	me.verify(1)
	me.buffer[me.dex] = byte(x)
	me.dex++
}

// PutUint8 func
func (me *ByteWriter) PutUint8(x uint8) {
	me.verify(1)
	me.buffer[me.dex] = byte(x)
	me.dex++
}

// PutUint16 func
func (me *ByteWriter) PutUint16(x uint16) {
	me.verify(2)
	me.buffer[me.dex] = byte(x)
	me.buffer[me.dex+1] = byte(x >> 8)
	me.dex += 2
}

// PutInt32 func
func (me *ByteWriter) PutInt32(x int32) {
	me.verify(4)
	me.buffer[me.dex] = byte(x)
	me.buffer[me.dex+1] = byte(x >> 8)
	me.buffer[me.dex+2] = byte(x >> 16)
	me.buffer[me.dex+3] = byte(x >> 24)
	me.dex += 4
}

// PutUint32 func
func (me *ByteWriter) PutUint32(x uint32) {
	me.verify(4)
	me.buffer[me.dex] = byte(x)
	me.buffer[me.dex+1] = byte(x >> 8)
	me.buffer[me.dex+2] = byte(x >> 16)
	me.buffer[me.dex+3] = byte(x >> 24)
	me.dex += 4
}

// PutFloat32 func
func (me *ByteWriter) PutFloat32(f float32) {
	me.verify(4)
	x := *(*int32)(unsafe.Pointer(&f))
	me.buffer[me.dex] = byte(x)
	me.buffer[me.dex+1] = byte(x >> 8)
	me.buffer[me.dex+2] = byte(x >> 16)
	me.buffer[me.dex+3] = byte(x >> 24)
	me.dex += 4
}

// SetUint8 func
func (me *ByteWriter) SetUint8(pos int, x uint8) {
	me.buffer[pos] = byte(x)
}

// Bytes func
func (me *ByteWriter) Bytes() []byte {
	return me.buffer[:me.dex]
}

// Position func
func (me *ByteWriter) Position(x int) {
	me.dex = x
}

// Reset func
func (me *ByteWriter) Reset() {
	me.dex = 0
}

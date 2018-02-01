package utils

import (
	"encoding/binary"
	"math"
	"strconv"
)

type Float16 uint16

// This code is adapted from Paul "Phernost" Tessier's answer on Stack Overflow:
// http://stackoverflow.com/questions/1659440/32-bit-to-16-bit-floating-point-conversion
//
// Ported from C to Go and generalized to 64-bit by William MacKay

const mantissaLength16 = 10
const mantissaLength32 = 23
const mantissaLength64 = 52

const shift32 = mantissaLength32 - mantissaLength16
const shift64 = mantissaLength64 - mantissaLength16

const shiftSign32 = 32 - 16
const shiftSign64 = 64 - 16

const infN32 = 0x7F800000         // flt32 infinity
const infN64 = 0x7FF0000000000000 // flt64 infinity
const maxN32 = 0x477FE000         // max flt16 normal as a flt32
const maxN64 = 0x40EFFC0000000000 // max flt16 normal as a flt64
const minN32 = 0x38800000         // min flt16 normal as a flt32
const minN64 = 0x3F10000000000000 // min flt16 normal as a flt64
const signN32 = 1 << (32 - 1)     // flt32 sign bit
const signN64 = 1 << (64 - 1)     // flt64 sign bit

const infC32 = infN32 >> shift32
const infC64 = infN64 >> shift64
const nanN32 = (infC32 + 1) << shift32 // minimum flt16 nan as a flt32
const nanN64 = (infC64 + 1) << shift64 // minimum flt16 nan as a flt64
const maxC32 = maxN32 >> shift32
const maxC64 = maxN64 >> shift64
const minC32 = minN32 >> shift32
const minC64 = minN64 >> shift64
const signC = 1 << (16 - 1) // flt16 sign bit

const mulN32 = 0x52000000         // (1 << mantissaLength32) / minN32
const mulN64 = 0x4410000000000000 // (1 << mantissaLength64) / minN64
const mulC32 = 0x33800000         // minN32 / (1 << mantissaLength16)
const mulC64 = 0x3E70000000000000 // minN64 / (1 << mantissaLength16)

const subC = ^(-1 << mantissaLength16) // max flt16 subnormal down shifted
const norC = 1 << mantissaLength16     // min flt16 normal down shifted

const maxD32 = infC32 - maxC32 - 1
const maxD64 = infC64 - maxC64 - 1
const minD32 = minC32 - subC - 1
const minD64 = minC64 - subC - 1

// converts a float32 to half-precision float
func From32(x float32) Float16 {
	v := math.Float32bits(x)
	sign := v & signN32
	v ^= sign
	sign >>= shiftSign32 // logical shift

	if minN32 > v {
		f := math.Float32frombits(mulN32)
		f *= math.Float32frombits(v) // correct subnormals
		v = uint32(f)
	}
	if infN32 > v && v > maxN32 {
		v = infN32
	}
	if nanN32 > v && v > infN32 {
		v = nanN32
	}
	v >>= shift32 // logical shift
	if v > maxC32 {
		v -= maxD32
	}
	if v > subC {
		v -= minD32
	}

	v |= sign
	return Float16(v)
}

// converts a float64 to half-precision float
func From64(x float64) Float16 {
	v := math.Float64bits(x)
	sign := v & signN64
	v ^= sign
	sign >>= shiftSign64 // logical shift

	if minN64 > v {
		f := math.Float64frombits(mulN64)
		f *= math.Float64frombits(v) // correct subnormals
		v = uint64(f)
	}
	if infN64 > v && v > maxN64 {
		v = infN64
	}
	if nanN64 > v && v > infN64 {
		v = nanN64
	}
	v >>= shift64 // logical shift
	if v > maxC64 {
		v -= maxD64
	}
	if v > subC {
		v -= minD64
	}

	v |= sign
	return Float16(v)
}

// converts a half-precision float to float32
func (x Float16) To32() float32 {
	v := uint32(x)
	sign := v & signC
	v ^= sign
	sign <<= shiftSign32
	if v > subC {
		v += minD32
	}
	if v > maxC32 {
		v += maxD32
	}

	v <<= shift32
	if norC > v {
		f := math.Float32frombits(mulC32)
		f *= float32(v)
		v = math.Float32bits(f)
	}

	v |= sign
	return math.Float32frombits(v)
}

// converts a half-precision float to float64
func (x Float16) To64() float64 {
	v := uint64(x)
	sign := v & signC
	v ^= sign
	sign <<= shiftSign64
	if v > subC {
		v += minD64
	}
	if v > maxC64 {
		v += maxD64
	}

	v <<= shift64
	if norC > v {
		f := math.Float64frombits(mulC64)
		f *= float64(v)
		v = math.Float64bits(f)
	}

	v |= sign
	return math.Float64frombits(v)
}

func FromBigEndian(b []byte) Float16 {
	x := binary.BigEndian.Uint16(b)
	return Float16(x)
}

func FromLittleEndian(b []byte) Float16 {
	x := binary.LittleEndian.Uint16(b)
	return Float16(x)
}

func (x Float16) PutBigEndian(b []byte) {
	binary.BigEndian.PutUint16(b, uint16(x))
}

func (x Float16) PutLittleEndian(b []byte) {
	binary.LittleEndian.PutUint16(b, uint16(x))
}

func (x Float16) BigEndian() []byte {
	b := make([]byte, 2)
	x.PutBigEndian(b)
	return b
}

func (x Float16) LittleEndian() []byte {
	b := make([]byte, 2)
	x.PutLittleEndian(b)
	return b
}

// stringer
func (x Float16) String() string {
	return strconv.FormatFloat(x.To64(), 'f', 7, 64)
}

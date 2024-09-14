package strct

import (
	"errors"
	"strconv"
	"strings"
)

const (
	BIG_NUM_LENTH = 2048
	BLOCK_NUM     = BIG_NUM_LENTH / 32
	INT_SIZE      = 4294967296
	SIZE          = 32
)

type BigInt struct {
	NumSlice [BLOCK_NUM]uint64 // 64 * 64 / 2 = 2048
}

func NewBigNumFromStringHex(hexNum string) (*BigInt, error) {
	bigNum := BigInt{}
	stringLen := len(hexNum)
	if stringLen > 512 {
		return nil, errors.New("number out of range")
	} else {
		k := 0
		j := 0
		for i := stringLen; i >= 8; i -= 8 {
			end := i - 8
			if end < 0 {
				end = stringLen
			}
			num, err := strconv.ParseUint(hexNum[end:i], 16, 64)
			if err != nil {
				return nil, err
			}
			bigNum.NumSlice[63-k] = num
			k++
			j = i
		}
		if j-8 > 0 {
			num, err := strconv.ParseUint(hexNum[0:j-8], 16, 64)
			if err != nil {
				return nil, err
			}
			bigNum.NumSlice[63-k] = num
		}
		return &bigNum, nil
	}
}

func formatUintWithLeadingZeros(num uint64, width int) string {
	hexStr := strconv.FormatUint(num, 16)

	if len(hexStr) < width {
		hexStr = strings.Repeat("0", width-len(hexStr)) + hexStr
	}

	return hexStr
}

func (bi *BigInt) ToHex() string {
	hexString := ""
	isCarryBit := true
	for i := 0; i < 64; i++ {
		num := bi.NumSlice[i]
		if isCarryBit == true {
			if num == 0 {
				continue
			} else {
				isCarryBit = false
			}
		}
		hexString += formatUintWithLeadingZeros(num, 8)

	}
	if hexString == "" {
		return "0"
	}
	return hexString
}

func LongCmp(num1, num2 *BigInt) int64 {
	borrow, res := LongSub(num1, num2)

	if borrow == 0 {
		if res.ToHex() == "0" {
			return 0
		} else {
			return 1
		}
	} else {
		return -1
	}
}

func LongMulOneDigit(bigNum *BigInt, num uint64) (*BigInt, uint64) {
	res := BigInt{}
	var carry uint64 = 0
	for i := 63; i >= 0; i-- {
		temp := bigNum.NumSlice[i]*num + carry
		res.NumSlice[i] = temp & (INT_SIZE - 1)
		carry = temp >> SIZE
	}
	return &res, carry
}

func LongAdd(bigInt1, bigInt2 *BigInt) (uint64, BigInt) {
	resBigNum := BigInt{}
	var carry uint64 = 0
	for i := 63; i >= 0; i-- {
		temp := bigInt1.NumSlice[i] + bigInt2.NumSlice[i] + carry
		resBigNum.NumSlice[i] = temp&INT_SIZE - 1
		carry = temp >> 32
		//fmt.Println(carry)
	}

	return carry, resBigNum
}

//func LongMul(num1, num2 *BigInt) *BigInt {
//	resBigNum := BigInt{}
//	for i := 63; i >= 0; i-- {
//		temp, tempCarry := LongMulOneDigit(num1, num2.NumSlice[i])
//
//	}
//}

func LongSub(bigInt1, bigInt2 *BigInt) (int64, *BigInt) {
	resBigNum := BigInt{}
	var borrow int64 = 0
	for i := 63; i >= 0; i-- {
		temp := int64(bigInt1.NumSlice[i]-bigInt2.NumSlice[i]) - borrow
		if temp >= 0 {
			resBigNum.NumSlice[i] = uint64(temp)
			borrow = 0
		} else {
			resBigNum.NumSlice[i] = uint64(INT_SIZE + temp)
			borrow = 1
		}
		//fmt.Println(bigInt1.NumSlice[i], "-", bigInt2.NumSlice[i], temp)
		//fmt.Println(resBigNum.NumSlice[i])
		//fmt.Println(borrow)
	}
	return borrow, &resBigNum
}

// "0xffffffffffffffff" => 18 446 744 073 709 551 615
//784abde982aa027d6fb10d409256f4ef79389842c6e41388735d12326d1829d9bfa21c15170e07eead079eb90eccc2602d17c924df43829d047624747a1cd312c32e66b1b35fc20254fb75c51234b32f7f11aa513fe4bcdfcc2549cfbde1dc625bc92ea559141dfab69b0fe973536935dde762c1f7036afe607876fafc7ffa5a9093322e59cdb58ca09858e87fc75e6936b6c23f
//784abde982aa027d6fb10d409256f4ef79389842c6e41388735d12326d1829d9bfa21c15170e07eead079eb9eccc2602d17c924df43829d47624747a1cd312c32e66b1b35fc20254fb75c51234b32f7f11aa513fe4bcdfcc2549cfbde1dc625bc92ea559141dfab69b0fe973536935dde762c1f7036afe607876fafc7ffa5a9093322e59cdb58ca09858e87fc75e6936b6c23f
//0xffff - 16bit
//0xffffffff - 32bit
//0xffffffffffffffff - 64bit
// 1 { 0xffffffffffffffff 0xffffffffffffffff 0xffffffffffffffff 0xffffffffffffffff ... 0xffffffffffffffff } 32 - 2048bit
// 1945670

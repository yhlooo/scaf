package units

import (
	"fmt"
	"math"
	"strconv"
)

// Format 单位格式
type Format string

const (
	// FormatSI International System of Units
	// 十进制
	FormatSI Format = "SI"
	// FormatIEC International Electrotechnical Commission
	// 二进制
	FormatIEC Format = "IEC"
)

// 常用单位
const (
	K = int64(1000)
	M = 1000 * K
	G = 1000 * M
	T = 1000 * G
	P = 1000 * T
	E = 1000 * P

	Ki = int64(1 << 10)
	Mi = int64(1 << 20)
	Gi = int64(1 << 30)
	Ti = int64(1 << 40)
	Pi = int64(1 << 50)
	Ei = int64(1 << 60)
)

var (
	siUnits  = []string{"K", "M", "G", "T", "P", "E"}
	iecUnits = []string{"Ki", "Mi", "Gi", "Ti", "Pi", "Ei"}
)

// NewIECValue 创建使用 IEC 格式单位的值
func NewIECValue(v int64) Value {
	return Value{
		Value:  v,
		Format: FormatIEC,
	}
}

// NewSIValue 创建使用 SI 格式单位的值
func NewSIValue(v int64) Value {
	return Value{
		Value:  v,
		Format: FormatSI,
	}
}

// Value 值
type Value struct {
	Value  int64
	Format Format
}

// String 返回字符串形式表示
func (v Value) String() string {
	if v.Value == 0 {
		return "0"
	}

	var unitNames []string
	var step int64
	switch v.Format {
	case FormatIEC:
		unitNames = iecUnits
		step = Ki
	default:
		unitNames = siUnits
		step = K
	}

	i := 0
	value := v.Value
	for ; i < len(unitNames); i++ {
		if value%step != 0 {
			break
		}
		value /= step
	}
	u := ""
	if i > 0 {
		u = unitNames[i-1]
	}
	return strconv.FormatInt(value, 10) + u
}

// RoundString 返回包含 ndigits 位小数的字符串形式
// 单位取尽可能大，但保证整数部分绝对值大于 0
func (v Value) RoundString(ndigits int) string {
	if ndigits < 0 {
		ndigits = 0
	}

	var unitNames []string
	var step float64
	switch v.Format {
	case FormatIEC:
		unitNames = iecUnits
		step = float64(Ki)
	default:
		unitNames = siUnits
		step = float64(K)
	}

	i := 0
	value := float64(v.Value)
	for ; i < len(unitNames); i++ {
		if math.Abs(value) < step {
			break
		}
		value /= step
	}
	u := ""
	if i > 0 {
		u = unitNames[i-1]
	}

	return fmt.Sprintf(fmt.Sprintf("%%.%df", ndigits), value) + u
}

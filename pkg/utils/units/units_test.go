package units

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValue_String 测试 Value.String 方法
func TestValue_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("0", Value{Value: 0, Format: FormatSI}.String())
	a.Equal("1", Value{Value: 1, Format: FormatSI}.String())
	a.Equal("-233", Value{Value: -233, Format: FormatSI}.String())
	a.Equal("1K", Value{Value: K, Format: FormatSI}.String())
	a.Equal("20M", Value{Value: 20 * M, Format: FormatSI}.String())
	a.Equal("999G", Value{Value: 999 * G, Format: FormatSI}.String())
	a.Equal("-32T", Value{Value: -32 * T, Format: FormatSI}.String())
	a.Equal("-72P", Value{Value: -72 * P, Format: FormatSI}.String())
	a.Equal("-5E", Value{Value: -5 * E, Format: FormatSI}.String())
	a.Equal("3001", Value{Value: 3*K + 1, Format: FormatSI}.String())

	a.Equal("0", Value{Value: 0, Format: FormatIEC}.String())
	a.Equal("1", Value{Value: 1, Format: FormatIEC}.String())
	a.Equal("-233", Value{Value: -233, Format: FormatIEC}.String())
	a.Equal("1Ki", Value{Value: Ki, Format: FormatIEC}.String())
	a.Equal("20Mi", Value{Value: 20 * Mi, Format: FormatIEC}.String())
	a.Equal("999Gi", Value{Value: 999 * Gi, Format: FormatIEC}.String())
	a.Equal("-32Ti", Value{Value: -32 * Ti, Format: FormatIEC}.String())
	a.Equal("-72Pi", Value{Value: -72 * Pi, Format: FormatIEC}.String())
	a.Equal("-5Ei", Value{Value: -5 * Ei, Format: FormatIEC}.String())
	a.Equal("3073", Value{Value: 3*Ki + 1, Format: FormatIEC}.String())
}

// TestValue_RoundString 测试 Value.RoundString 方法
func TestValue_RoundString(t *testing.T) {
	a := assert.New(t)
	a.Equal("0", Value{Value: 0, Format: FormatSI}.RoundString(0))
	a.Equal("0.00", Value{Value: 0, Format: FormatSI}.RoundString(2))
	a.Equal("1", Value{Value: 1, Format: FormatSI}.RoundString(0))
	a.Equal("0.0", Value{Value: 0, Format: FormatSI}.RoundString(1))
	a.Equal("233", Value{Value: 233, Format: FormatSI}.RoundString(0))
	a.Equal("-233.00", Value{Value: -233, Format: FormatSI}.RoundString(2))
	a.Equal("-1.23K", Value{Value: -1230, Format: FormatSI}.RoundString(2))
	a.Equal("1.23M", Value{Value: 1230000, Format: FormatSI}.RoundString(2))
	a.Equal("2.34Mi", Value{Value: 2452619, Format: FormatIEC}.RoundString(2))
	a.Equal("-1.00Gi", Value{Value: -1073741824, Format: FormatIEC}.RoundString(2))
}

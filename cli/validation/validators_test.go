package validation

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomError(t *testing.T) {
	assert.NotEmpty(t, Var("42", Max(0).Error("")))
	assert.Error(t, errors.New("test"), Var("42", Max(0).Error("test")))
}

func TestCustomErrorf(t *testing.T) {
	assert.NotEmpty(t, Var("42", Max(0).Errorf("")))
	assert.Error(t, errors.New("test"), Var("42", Max(0).Errorf("%s", "test")))
}

func TestWhen(t *testing.T) {
	assert.NoError(t, Var("42", Max(0).When(false == true)))
	assert.Error(t, Var("42", Max(0).When(false == false)))
}

func TestTags(t *testing.T) {
	assert.NoError(t, Var("42", Tags("")))
	assert.NoError(t, Var(42, Tags("max=42")))
	assert.NoError(t, Var(nil, Tags("omitempty,min=42")))
}

func TestOmitEmpty(t *testing.T) {
	assert.Error(t, Var(1, OmitEmpty(), Min(42)))
	assert.NoError(t, Var(nil, OmitEmpty(), Min(42)))
	assert.NoError(t, Var(0, OmitEmpty(), Min(42)))
	assert.NoError(t, Var([]string{}, OmitEmpty(), Min(42)))
	assert.NoError(t, Var(&[]string{}, OmitEmpty(), Min(42)))
}

func TestSkip(t *testing.T) {
	assert.NoError(t, Var(nil, Skip(), Required()))
	assert.Error(t, Var(nil, Required(), Skip()))
}

func TestRequired(t *testing.T) {
	assert.Error(t, Var(nil, Required()))
	assert.Error(t, Var("", Required()))
	assert.NoError(t, Var("test", Required()))
}

func TestMin(t *testing.T) {
	assert.NoError(t, Var(42, Min(42)))
	assert.Error(t, Var(0, Min(42)))

	assert.Error(t, Var([]int{}, Min(1)))
	assert.NoError(t, Var([]int{42}, Min(1)))
}

func TestMax(t *testing.T) {
	assert.NoError(t, Var(42, Max(42)))
	assert.Error(t, Var(42, Max(0)))

	assert.NoError(t, Var([]int{}, Max(1)))
	assert.Error(t, Var([]int{42, 42}, Max(1)))
}

func TestLen(t *testing.T) {
	assert.NoError(t, Var([]string{"1"}, Len(1)))
	assert.Error(t, Var([]int{42, 42}, Len(1)))
}

func TestMinLen(t *testing.T) {
	assert.NoError(t, Var(42, MinLen(42)))
	assert.Error(t, Var(0, MinLen(42)))

	assert.Error(t, Var([]int{}, MinLen(1)))
	assert.NoError(t, Var([]int{42}, MinLen(1)))
}

func TestMaxLen(t *testing.T) {
	assert.NoError(t, Var(42, MaxLen(42)))
	assert.Error(t, Var(42, MaxLen(0)))

	assert.NoError(t, Var([]int{}, MaxLen(1)))
	assert.Error(t, Var([]int{42, 42}, MaxLen(1)))
}

func TestIP(t *testing.T) {
	assert.Error(t, Var("42", IP()))
	assert.NoError(t, Var("192.168.113.1", IP()))
	assert.NoError(t, Var("2001:db8:3333:4444:5555:6666:7777:8888", IP()))
}

func TestIPv4(t *testing.T) {
	assert.Error(t, Var("42", IPv4()))
	assert.NoError(t, Var("192.168.113.1", IPv4()))
	assert.Error(t, Var("2001:db8:3333:4444:5555:6666:7777:8888", IPv4()))
}

func TestIPv6(t *testing.T) {
	assert.Error(t, Var("42", IPv6()))
	assert.Error(t, Var("192.168.113.1", IPv6()))
	assert.NoError(t, Var("2001:db8:3333:4444:5555:6666:7777:8888", IPv6()))
}

func TestCIDR(t *testing.T) {
	assert.Error(t, Var("42", CIDR()))
	assert.NoError(t, Var("192.168.113.1/24", CIDR()))
	assert.NoError(t, Var("2001:db8:3333:4444:5555:6666:7777:8888/64", CIDR()))
}

func TestCIDRv4(t *testing.T) {
	assert.Error(t, Var("42", CIDRv4()))
	assert.NoError(t, Var("192.168.113.1/24", CIDRv4()))
	assert.Error(t, Var("2001:db8:3333:4444:5555:6666:7777:8888/64", CIDRv4()))
}

func TestCIDRv6(t *testing.T) {
	assert.Error(t, Var("42", CIDRv6()))
	assert.Error(t, Var("192.168.113.1/24", CIDRv6()))
	assert.NoError(t, Var("2001:db8:3333:4444:5555:6666:7777:8888/64", CIDRv6()))
}

func TestIPInRange(t *testing.T) {
	assert.Error(t, Var("192.168.113.1", IPInRange("")))
	assert.Error(t, Var("", IPInRange("192.168.113.1")))
	assert.Error(t, Var("192.168.112.0", IPInRange("192.168.113.1/24")))
	assert.NoError(t, Var("192.168.113.0", IPInRange("192.168.113.1/24")))
	assert.NoError(t, Var("192.168.113.113", IPInRange("192.168.113.1/24")))
}

func TestMAC(t *testing.T) {
	assert.Error(t, Var("42", MAC()))
	assert.NoError(t, Var("AA:BB:CC:DD:EE:FF", MAC()))
}

func TestOneOf(t *testing.T) {
	assert.NoError(t, Var("42", OneOf(1, 42, 24)))
	assert.NoError(t, Var(24, OneOf(1, 42, 24)))
	assert.Error(t, Var(7, OneOf(1, 42, 24)))
	assert.Error(t, Var(7, OneOf()))
}

func TestAlpha(t *testing.T) {
	assert.Error(t, Var("42", Alpha()))
	assert.Error(t, Var(42, Alpha()))
	assert.Error(t, Var("42.0", Alpha()))
	assert.Error(t, Var(42.0, Alpha()))
	assert.Error(t, Var("42aAbB", Alpha()))
	assert.NoError(t, Var("aAbB", Alpha()))
	assert.Error(t, Var("-", Alpha()))
	assert.Error(t, Var("_", Alpha()))
	assert.Error(t, Var(nil, Alpha()))
}

func TestNumeric(t *testing.T) {
	assert.NoError(t, Var("42", Numeric()))
	assert.NoError(t, Var(42, Numeric()))
	assert.NoError(t, Var("42.0", Numeric()))
	assert.NoError(t, Var(42.0, Numeric()))
	assert.Error(t, Var("42aAbB", Numeric()))
	assert.Error(t, Var("-", Numeric()))
	assert.Error(t, Var("_", Numeric()))
	assert.Error(t, Var(nil, Alpha()))
}

func TestAlphaNumeric(t *testing.T) {
	assert.NoError(t, Var("42", AlphaNumeric()))
	assert.Error(t, Var(42, AlphaNumeric()))
	assert.Error(t, Var("42.0", AlphaNumeric()))
	assert.Error(t, Var(42.0, AlphaNumeric()))
	assert.NoError(t, Var("42aAbB", AlphaNumeric()))
	assert.NoError(t, Var("aAbB", AlphaNumeric()))
	assert.Error(t, Var("-", AlphaNumeric()))
	assert.Error(t, Var("_", AlphaNumeric()))
	assert.Error(t, Var(nil, AlphaNumeric()))
}

func TestAlphaNumericHyp(t *testing.T) {
	assert.NoError(t, Var("42", AlphaNumericHyp()))
	assert.Error(t, Var(42, AlphaNumericHyp()))
	assert.Error(t, Var("42.0", AlphaNumericHyp()))
	assert.Error(t, Var(42.0, AlphaNumericHyp()))
	assert.NoError(t, Var("42aAbB", AlphaNumericHyp()))
	assert.NoError(t, Var("aAbB", AlphaNumericHyp()))
	assert.NoError(t, Var("-", AlphaNumericHyp()))
	assert.Error(t, Var("_", AlphaNumericHyp()))
	assert.Error(t, Var(nil, AlphaNumericHyp()))
}

func TestAlphaNumericHypUS(t *testing.T) {
	assert.NoError(t, Var("42", AlphaNumericHypUS()))
	assert.Error(t, Var(42, AlphaNumericHypUS()))
	assert.Error(t, Var("42.0", AlphaNumericHypUS()))
	assert.Error(t, Var(42.0, AlphaNumericHypUS()))
	assert.NoError(t, Var("42aAbB", AlphaNumericHypUS()))
	assert.NoError(t, Var("aAbB", AlphaNumericHypUS()))
	assert.NoError(t, Var("-", AlphaNumericHypUS()))
	assert.NoError(t, Var("_", AlphaNumericHypUS()))
	assert.Error(t, Var(nil, AlphaNumericHypUS()))
}

func TestLowercase(t *testing.T) {
	assert.NoError(t, Var("42", Lowercase()))
	assert.Error(t, Var("42aAbB", Lowercase()))
	assert.Error(t, Var("aAbB", Lowercase()))
	assert.Error(t, Var("AB", Lowercase()))
	assert.NoError(t, Var("ab", Lowercase()))
	assert.NoError(t, Var("-", Lowercase()))
	assert.NoError(t, Var("_", Lowercase()))
	assert.Error(t, Var(nil, Lowercase()))
}

func TestUppercase(t *testing.T) {
	assert.NoError(t, Var("42", Uppercase()))
	assert.Error(t, Var("42aAbB", Uppercase()))
	assert.Error(t, Var("aAbB", Uppercase()))
	assert.NoError(t, Var("AB", Uppercase()))
	assert.Error(t, Var("ab", Uppercase()))
	assert.NoError(t, Var("-", Uppercase()))
	assert.NoError(t, Var("_", Uppercase()))
	assert.Error(t, Var(nil, Uppercase()))
}

func TestFileExists(t *testing.T) {
	assert.NoError(t, Var("./validators_test.go", FileExists()))
	assert.Error(t, Var("./non-existing-file-test", FileExists()))
}

func TestURL(t *testing.T) {
	assert.NoError(t, Var("https://kubitect.io", URL()))
	assert.Error(t, Var("kubitect.io", URL()))
	assert.Error(t, Var(nil, URL()))
}

func TestSemVer(t *testing.T) {
	assert.Error(t, Var("v1.2.3", SemVer()))
	assert.NoError(t, Var("1.2.3", SemVer()))
	assert.Error(t, Var("1.2", SemVer()))
	assert.Error(t, Var("1", SemVer()))
	assert.Error(t, Var("", SemVer()))
	assert.Error(t, Var("1.2.*", SemVer()))
	assert.Error(t, Var("a.b.c", SemVer()))
	assert.Error(t, Var(nil, SemVer()))
}

func TestVSemVer(t *testing.T) {
	assert.NoError(t, Var("v1.2.3", VSemVer()))
	assert.Error(t, Var("1.2.3", VSemVer()))
	assert.Error(t, Var("1.2", VSemVer()))
	assert.Error(t, Var("1", VSemVer()))
	assert.Error(t, Var("", VSemVer()))
	assert.Error(t, Var("1.2.*", VSemVer()))
	assert.Error(t, Var("a.b.c", VSemVer()))
	assert.Error(t, Var(nil, VSemVer()))
}

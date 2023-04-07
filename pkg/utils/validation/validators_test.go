package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomError(t *testing.T) {
	assert.NotEmpty(t, Var("42", Max(0).Error("")))
	assert.EqualError(t, Var("42", Max(0).Error("test")), "test")
}

func TestCustomErrorf(t *testing.T) {
	assert.NotEmpty(t, Var("42", Max(0).Errorf("")))
	assert.EqualError(t, Var("42", Max(0).Errorf("%s", "test")), "test")
}

func TestWhen(t *testing.T) {
	assert.NoError(t, Var("42", Max(0).When(false == true)))
	assert.Error(t, Var("42", Max(0).When(false == false)))
}

func TestNoneValidator(t *testing.T) {
	assert.NoError(t, Var(nil, None))
	assert.NoError(t, Var(5, None))
}

func TestCustomValidator(t *testing.T) {
	key := "test"

	RegisterCustomValidator(key, Min(42))

	assert.Error(t, Var(0, Custom(key)))
	assert.NoError(t, Var(42, Custom(key)))

	ClearCustomValidators()

	assert.NoError(t, Var(0, Custom(key)))
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

func TestFail(t *testing.T) {
	assert.Error(t, Var(nil, Fail()))
	assert.Error(t, Var(42, Fail()))
}

func TestRequired(t *testing.T) {
	type Struct struct {
		value int
	}

	assert.Error(t, Var(nil, Required()))
	assert.Error(t, Var("", Required()))
	assert.NoError(t, Var(Struct{}, Required()))
	assert.NoError(t, Var("test", Required()))
}

func TestNotEmpty(t *testing.T) {
	type Struct struct {
		value int
	}

	s1 := Struct{
		value: 0,
	}

	s2 := Struct{
		value: 1,
	}

	assert.Error(t, Var(nil, NotEmpty()))
	assert.Error(t, Var("", NotEmpty()))
	assert.Error(t, Var(Struct{}, NotEmpty()))
	assert.Error(t, Var(s1, NotEmpty()))
	assert.NoError(t, Var(s2, NotEmpty()))
	assert.NoError(t, Var("test", NotEmpty()))
}

func TestUnique(t *testing.T) {
	assert.Error(t, Var(nil, Unique()))
	assert.Error(t, Var([]string{"a", "a"}, Unique()))
	assert.NoError(t, Var([]string{"a", "b"}, Unique()))
	assert.NoError(t, Var([]string{}, Unique()))
}

func TestUniqueField(t *testing.T) {
	type S struct {
		X string
		V int
	}

	type ps struct {
		v int
	}

	fn1 := func() {
		Var([]ps{{42}, {42}}, UniqueField("v"))
	}

	fn2 := func() {
		Var([]ps{{}}, UniqueField(""))
	}

	assert.PanicsWithError(t, ErrorExportInterface.Error(), fn1)
	assert.PanicsWithError(t, ErrorFieldNotFound.Error(), fn2)

	assert.Error(t, Var(nil, UniqueField("V")))
	assert.Error(t, Var([]S{{}, {}}, UniqueField("V")))
	assert.Error(t, Var([]S{{V: 42}, {V: 42}}, UniqueField("V")))
	assert.NoError(t, Var([]S{{V: 41}, {V: 42}}, UniqueField("V")))
	assert.NoError(t, Var([]S{}, UniqueField("")))
	assert.NoError(t, Var([]S{}, UniqueField("V")))
	assert.NoError(t, Var([]S{{}}, UniqueField("V")))
	assert.NoError(t, Var([]S{{V: 42}}, UniqueField("V")))
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

func TestFilePath(t *testing.T) {
	assert.NoError(t, Var("./validators_test.go", FilePath()))
	assert.NoError(t, Var("./non-existing-file-test", FilePath()))
	assert.Error(t, Var("", FilePath()))
	// Fails the test (nil pointer)
	// assert.NoError(t, Var("/etc", FilePath()))
}

func TestDirPath(t *testing.T) {
	assert.NoError(t, Var("/etc", DirPath()))
	assert.NoError(t, Var("/", DirPath()))
	assert.Error(t, Var("", DirPath()))
	// Fails the test
	// assert.NoError(t, Var("/dir", DirPath()))
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

func TestRegexAny(t *testing.T) {
	regex := []string{"^[0-9][0-9]?$", "abc"}
	assert.NoError(t, Var("abc", RegexAny(regex...)))
	assert.NoError(t, Var("1", RegexAny(regex...)))
	assert.NoError(t, Var("13", RegexAny(regex...)))
	assert.EqualError(t, Var("a", RegexAny(regex...)), "Field does not match any regex expression [^[0-9][0-9]?$ abc]. (actual: a)")
}

func TestRegexAll(t *testing.T) {
	regex := []string{"^*[0-9]$", "^abc"}
	assert.NoError(t, Var("abc4", RegexAll(regex...)))
	assert.EqualError(t, Var("abc", RegexAll(regex...)), "Field does not match all regex expressions [^*[0-9]$ ^abc]. (actual: abc)")
}

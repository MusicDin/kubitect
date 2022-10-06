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

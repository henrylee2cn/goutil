package tpack

import (
	"reflect"
	"testing"
	"time"
	"unsafe"
)

func TestRuntimeTypeID(t *testing.T) {
	type (
		GoTime = time.Time
		Time   time.Time
		I2     interface {
			String() string
		}
		I1 interface {
			UnixNano() int64
			I2
		}
	)
	t0 := new(time.Time)
	t1 := Unpack(t0).RuntimeTypeID()
	t2 := Unpack(new(GoTime)).RuntimeTypeID()
	t3 := Unpack(new(Time)).RuntimeTypeID()
	t.Log(t1, t2, t3)
	e0 := time.Time{}
	e1 := Unpack(e0).RuntimeTypeID()
	e2 := Unpack(GoTime{}).RuntimeTypeID()
	e3 := Unpack(Time{}).RuntimeTypeID()
	i := Unpack(I2(I1(&GoTime{}))).RuntimeTypeID()
	if t1 != t2 || t1 != e1 || t1 != e2 || t1 != i || t3 != e3 {
		t.FailNow()
	}
	t.Log(e1, e2, e3, i, RuntimeTypeID(reflect.TypeOf(t0)), Unpack(t0.String).RuntimeTypeID())
}

func TestKind(t *testing.T) {
	type X struct {
		A int16
		B string
	}
	var x X
	k := Unpack(&x).Kind()
	t.Log(k)
	if k != reflect.Ptr {
		t.FailNow()
	}
	k = Unpack(x).Kind()
	t.Log(k)
	if k != reflect.Struct {
		t.FailNow()
	}
	f := func() {}
	k = Unpack(f).Kind()
	t.Log(k)
	if k != reflect.Func {
		t.FailNow()
	}
	k = Unpack(t.Name).Kind()
	t.Log(k)
	if k != reflect.Func {
		t.FailNow()
	}
}

func TestPointer(t *testing.T) {
	type X struct {
		A int16
		B string
	}
	x := X{A: 12345, B: "test"}
	if Unpack(&x).Pointer() != reflect.ValueOf(&x).Pointer() {
		t.FailNow()
	}
	elemPtr := Unpack(x).Pointer()
	a := *(*int16)(unsafe.Pointer(elemPtr))
	if a != x.A {
		t.FailNow()
	}
	b := *(*string)(unsafe.Pointer(elemPtr + unsafe.Offsetof(x.B)))
	if b != x.B {
		t.FailNow()
	}
	f := func() {}
	if Unpack(f).Pointer() != reflect.ValueOf(f).Pointer() {
		t.FailNow()
	}
	if Unpack(t.Name).Pointer() != reflect.ValueOf(t.Name).Pointer() {
		t.FailNow()
	}
	s := []string{""}
	if Unpack(s).Pointer() != reflect.ValueOf(s).Pointer() {
		t.FailNow()
	}
}

func BenchmarkUnpack_tpack(b *testing.B) {
	b.StopTimer()
	type T struct {
		a int
	}
	var t = new(T)
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = Unpack(t).RuntimeTypeID()
	}
}

func BenchmarkTypeOf_go(b *testing.B) {
	b.StopTimer()
	type T struct {
		a int
	}
	var t = new(T)
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = reflect.TypeOf(t).Elem().String()
	}
}

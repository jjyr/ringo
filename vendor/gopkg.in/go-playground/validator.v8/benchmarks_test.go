package validator

import (
	sql "database/sql/driver"
	"testing"
	"time"
)

func BenchmarkFieldSuccess(b *testing.B) {
	for n := 0; n < b.N; n++ {
		validate.Field("1", "len=1")
	}
}

func BenchmarkFieldFailure(b *testing.B) {
	for n := 0; n < b.N; n++ {
		validate.Field("2", "len=1")
	}
}

func BenchmarkFieldDiveSuccess(b *testing.B) {
	for n := 0; n < b.N; n++ {
		validate.Field([]string{"val1", "val2", "val3"}, "required,dive,required")
	}
}

func BenchmarkFieldDiveFailure(b *testing.B) {
	for n := 0; n < b.N; n++ {
		validate.Field([]string{"val1", "", "val3"}, "required,dive,required")
	}
}

func BenchmarkFieldCustomTypeSuccess(b *testing.B) {

	validate.RegisterCustomTypeFunc(ValidateValuerType, (*sql.Valuer)(nil), valuer{})

	val := valuer{
		Name: "1",
	}

	for n := 0; n < b.N; n++ {
		validate.Field(val, "len=1")
	}
}

func BenchmarkFieldCustomTypeFailure(b *testing.B) {

	// validate := New(Config{TagName: "validate"})
	validate.RegisterCustomTypeFunc(ValidateValuerType, (*sql.Valuer)(nil), valuer{})

	val := valuer{}

	for n := 0; n < b.N; n++ {
		validate.Field(val, "len=1")
	}
}

func BenchmarkFieldOrTagSuccess(b *testing.B) {
	for n := 0; n < b.N; n++ {
		validate.Field("rgba(0,0,0,1)", "rgb|rgba")
	}
}

func BenchmarkFieldOrTagFailure(b *testing.B) {
	for n := 0; n < b.N; n++ {
		validate.Field("#000", "rgb|rgba")
	}
}

func BenchmarkStructSimpleCustomTypeSuccess(b *testing.B) {

	validate.RegisterCustomTypeFunc(ValidateValuerType, (*sql.Valuer)(nil), valuer{})

	val := valuer{
		Name: "1",
	}

	type Foo struct {
		Valuer   valuer `validate:"len=1"`
		IntValue int    `validate:"min=5,max=10"`
	}

	validFoo := &Foo{Valuer: val, IntValue: 7}

	for n := 0; n < b.N; n++ {
		validate.Struct(validFoo)
	}
}

func BenchmarkStructSimpleCustomTypeFailure(b *testing.B) {

	validate.RegisterCustomTypeFunc(ValidateValuerType, (*sql.Valuer)(nil), valuer{})

	val := valuer{}

	type Foo struct {
		Valuer   valuer `validate:"len=1"`
		IntValue int    `validate:"min=5,max=10"`
	}

	validFoo := &Foo{Valuer: val, IntValue: 3}

	for n := 0; n < b.N; n++ {
		validate.Struct(validFoo)
	}
}

func BenchmarkStructPartialSuccess(b *testing.B) {

	type Test struct {
		Name     string `validate:"required"`
		NickName string `validate:"required"`
	}

	test := &Test{
		Name: "Joey Bloggs",
	}

	for n := 0; n < b.N; n++ {
		validate.StructPartial(test, "Name")
	}
}

func BenchmarkStructPartialFailure(b *testing.B) {

	type Test struct {
		Name     string `validate:"required"`
		NickName string `validate:"required"`
	}

	test := &Test{
		Name: "Joey Bloggs",
	}

	for n := 0; n < b.N; n++ {
		validate.StructPartial(test, "NickName")
	}
}

func BenchmarkStructExceptSuccess(b *testing.B) {

	type Test struct {
		Name     string `validate:"required"`
		NickName string `validate:"required"`
	}

	test := &Test{
		Name: "Joey Bloggs",
	}

	for n := 0; n < b.N; n++ {
		validate.StructPartial(test, "Nickname")
	}
}

func BenchmarkStructExceptFailure(b *testing.B) {

	type Test struct {
		Name     string `validate:"required"`
		NickName string `validate:"required"`
	}

	test := &Test{
		Name: "Joey Bloggs",
	}

	for n := 0; n < b.N; n++ {
		validate.StructPartial(test, "Name")
	}
}

func BenchmarkStructSimpleCrossFieldSuccess(b *testing.B) {

	type Test struct {
		Start time.Time
		End   time.Time `validate:"gtfield=Start"`
	}

	now := time.Now().UTC()
	then := now.Add(time.Hour * 5)

	test := &Test{
		Start: now,
		End:   then,
	}

	for n := 0; n < b.N; n++ {
		validate.Struct(test)
	}
}

func BenchmarkStructSimpleCrossFieldFailure(b *testing.B) {

	type Test struct {
		Start time.Time
		End   time.Time `validate:"gtfield=Start"`
	}

	now := time.Now().UTC()
	then := now.Add(time.Hour * -5)

	test := &Test{
		Start: now,
		End:   then,
	}

	for n := 0; n < b.N; n++ {
		validate.Struct(test)
	}
}

func BenchmarkStructSimpleCrossStructCrossFieldSuccess(b *testing.B) {

	type Inner struct {
		Start time.Time
	}

	type Outer struct {
		Inner     *Inner
		CreatedAt time.Time `validate:"eqcsfield=Inner.Start"`
	}

	now := time.Now().UTC()

	inner := &Inner{
		Start: now,
	}

	outer := &Outer{
		Inner:     inner,
		CreatedAt: now,
	}

	for n := 0; n < b.N; n++ {
		validate.Struct(outer)
	}
}

func BenchmarkStructSimpleCrossStructCrossFieldFailure(b *testing.B) {

	type Inner struct {
		Start time.Time
	}

	type Outer struct {
		Inner     *Inner
		CreatedAt time.Time `validate:"eqcsfield=Inner.Start"`
	}

	now := time.Now().UTC()
	then := now.Add(time.Hour * 5)

	inner := &Inner{
		Start: then,
	}

	outer := &Outer{
		Inner:     inner,
		CreatedAt: now,
	}

	for n := 0; n < b.N; n++ {
		validate.Struct(outer)
	}
}

func BenchmarkStructSimpleSuccess(b *testing.B) {

	type Foo struct {
		StringValue string `validate:"min=5,max=10"`
		IntValue    int    `validate:"min=5,max=10"`
	}

	validFoo := &Foo{StringValue: "Foobar", IntValue: 7}

	for n := 0; n < b.N; n++ {
		validate.Struct(validFoo)
	}
}

func BenchmarkStructSimpleFailure(b *testing.B) {

	type Foo struct {
		StringValue string `validate:"min=5,max=10"`
		IntValue    int    `validate:"min=5,max=10"`
	}

	invalidFoo := &Foo{StringValue: "Fo", IntValue: 3}

	for n := 0; n < b.N; n++ {
		validate.Struct(invalidFoo)
	}
}

func BenchmarkStructSimpleSuccessParallel(b *testing.B) {

	type Foo struct {
		StringValue string `validate:"min=5,max=10"`
		IntValue    int    `validate:"min=5,max=10"`
	}

	validFoo := &Foo{StringValue: "Foobar", IntValue: 7}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			validate.Struct(validFoo)
		}
	})
}

func BenchmarkStructSimpleFailureParallel(b *testing.B) {

	type Foo struct {
		StringValue string `validate:"min=5,max=10"`
		IntValue    int    `validate:"min=5,max=10"`
	}

	invalidFoo := &Foo{StringValue: "Fo", IntValue: 3}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			validate.Struct(invalidFoo)
		}
	})
}

func BenchmarkStructComplexSuccess(b *testing.B) {

	tSuccess := &TestString{
		Required:  "Required",
		Len:       "length==10",
		Min:       "min=1",
		Max:       "1234567890",
		MinMax:    "12345",
		Lt:        "012345678",
		Lte:       "0123456789",
		Gt:        "01234567890",
		Gte:       "0123456789",
		OmitEmpty: "",
		Sub: &SubTest{
			Test: "1",
		},
		SubIgnore: &SubTest{
			Test: "",
		},
		Anonymous: struct {
			A string `validate:"required"`
		}{
			A: "1",
		},
		Iface: &Impl{
			F: "123",
		},
	}

	for n := 0; n < b.N; n++ {
		validate.Struct(tSuccess)
	}
}

func BenchmarkStructComplexFailure(b *testing.B) {

	tFail := &TestString{
		Required:  "",
		Len:       "",
		Min:       "",
		Max:       "12345678901",
		MinMax:    "",
		Lt:        "0123456789",
		Lte:       "01234567890",
		Gt:        "1",
		Gte:       "1",
		OmitEmpty: "12345678901",
		Sub: &SubTest{
			Test: "",
		},
		Anonymous: struct {
			A string `validate:"required"`
		}{
			A: "",
		},
		Iface: &Impl{
			F: "12",
		},
	}

	for n := 0; n < b.N; n++ {
		validate.Struct(tFail)
	}
}

func BenchmarkStructComplexSuccessParallel(b *testing.B) {

	tSuccess := &TestString{
		Required:  "Required",
		Len:       "length==10",
		Min:       "min=1",
		Max:       "1234567890",
		MinMax:    "12345",
		Lt:        "012345678",
		Lte:       "0123456789",
		Gt:        "01234567890",
		Gte:       "0123456789",
		OmitEmpty: "",
		Sub: &SubTest{
			Test: "1",
		},
		SubIgnore: &SubTest{
			Test: "",
		},
		Anonymous: struct {
			A string `validate:"required"`
		}{
			A: "1",
		},
		Iface: &Impl{
			F: "123",
		},
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			validate.Struct(tSuccess)
		}
	})
}

func BenchmarkStructComplexFailureParallel(b *testing.B) {

	tFail := &TestString{
		Required:  "",
		Len:       "",
		Min:       "",
		Max:       "12345678901",
		MinMax:    "",
		Lt:        "0123456789",
		Lte:       "01234567890",
		Gt:        "1",
		Gte:       "1",
		OmitEmpty: "12345678901",
		Sub: &SubTest{
			Test: "",
		},
		Anonymous: struct {
			A string `validate:"required"`
		}{
			A: "",
		},
		Iface: &Impl{
			F: "12",
		},
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			validate.Struct(tFail)
		}
	})
}

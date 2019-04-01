package main

import (
	"testing"
)

/*
 * This test is checking whatever it's worth to pass something via pointer to speedup
 */


type subdata struct {
	 number1 int64
	 fl1     float64
	 b1      bool
	 s1      string
 }

type bigStruct struct {
	number1 int64
	number2 int64
	number3 int64
	number4 int64
	number5 int64
	number6 int64
	number7 int64
	number8 int64
	number9 int64
	number10 int64
	fl1 float64
	fl2 float64
	fl3 float64
	fl4 float64
	fl5 float64
	fl6 float64
	fl7 float64
	fl8 float64
	fl9 float64
	fl10 float64
	b1 bool
	b2 bool
	b3 bool
	b4 bool
	b5 bool
	b6 bool
	b7 bool
	b8 bool
	b9 bool
	b10 bool
	s1 string
	s2 string
	s3 string
	s4 string
	s5 string
	s6 string
	s7 string
	s8 string
	s9 string
	s10 string
	sub1 subdata
	sub2 subdata
	sub3 subdata
	sub4 subdata
}


func BenchmarkStructCopying(t *testing.B) {
	calc := func(st bigStruct) int64 {
		return st.number1
	}
	calcPtr := func(st *bigStruct) int64 {
		return st.number1
	}

	var bs bigStruct

	t.Run("by ptr", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			calcPtr(&bs)
		}
	})

	t.Run("by value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			calc(bs)
		}
	})
}
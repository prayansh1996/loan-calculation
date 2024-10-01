package main

import (
	"math"
	"testing"
)

func TestCalculateMonthlyInterestRate(t *testing.T) {
	calculation := calculateMonthlyInterestRate(11.0)
	t.Logf("\nCalculation: %v", calculation)

	want := 1.00873459382
	if math.Abs(calculation-want) > 0.000001 {
		t.Fatalf("calculateMonthlyInterestRate calculation incorrect, want: %v, actual %v", want, calculation)
	}
}

func TestCalculateMonthlyEmi(t *testing.T) {
	calculation := calculateMonthlyEmi(20_00_000.0, 8.0, 7.0)
	t.Logf("\nCalculation: %v", calculation)

	want := 31_172.0
	if math.Abs(calculation-want) > 1 {
		t.Fatalf("calculateMonthlyEmi calculation incorrect, want: %v, actual %v", want, calculation)
	}
}

func TestCalculateSipReturns(t *testing.T) {
	calculation := calculateSipReturns(10000.0, 11.0, 7.0)
	t.Logf("\nCalculation: %v", calculation)

	want := 12_68_471.0
	if math.Abs(calculation-want) > 1 {
		t.Fatalf("calculateSipReturns calculation incorrect, want: %v, actual %v", want, calculation)
	}
}

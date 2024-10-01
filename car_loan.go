package main

import (
	"fmt"
	"math"

	"github.com/go-playground/locales/currency"
	"github.com/go-playground/locales/en_IN"
)

func padding(tabs int) string {
	s := ""
	for i := 0; i < tabs; i++ {
		s += "\t"
	}
	return s
}

func calculateMonthlyInterestRate(yearlyInterestPercent float64) float64 {
	return math.Pow(1.0+yearlyInterestPercent/100.0, 1.0/12.0)
}

func calculateMonthlyEmi(
	loanAmount float64,
	yearlyLoanInterestPercent float64,
	loanDurationYears float64,
) float64 {
	monthlyLoanInterestRate := yearlyLoanInterestPercent / (12.0 * 100.0)
	months := loanDurationYears * 12.0
	exponent := math.Pow(1.0+monthlyLoanInterestRate, months)
	return loanAmount * monthlyLoanInterestRate * exponent / (exponent - 1.0)
}

func calculateSipReturns(
	monthlySipAmount float64,
	yearlyInvestmentReturnPercent float64,
	investmentDurationYears float64,
) float64 {
	monthlyInterestRate := yearlyInvestmentReturnPercent / (12.0 * 100.0)
	numerator := monthlySipAmount * (math.Pow(1.0+monthlyInterestRate, 12.0*investmentDurationYears) - 1.0) * (1.0 + monthlyInterestRate)
	return numerator / monthlyInterestRate
}

func calculate_intialCapitalInvested_liquidatedMonthlyToPayEmi(
	yearlyInvestmentReturnPercent float64,
	yearlyLoanInterestPercent float64,
	loanAmount float64,
	downPayment float64,
	loanDurationYears float64,
) float64 {
	totalMoneyRemaining := loanAmount - downPayment
	monthlyInvestmentReturnRate := calculateMonthlyInterestRate(yearlyInvestmentReturnPercent)
	monthlyEmi := calculateMonthlyEmi(totalMoneyRemaining, yearlyLoanInterestPercent, loanDurationYears)

	for i := 0.0; i < 12.0*loanDurationYears; i++ {
		totalMoneyRemaining = monthlyInvestmentReturnRate*totalMoneyRemaining - monthlyEmi
	}
	return totalMoneyRemaining
}

func calculate_intialCapitalInvested_liquidatedYearlyToPayEmi(
	yearlyInvestmentReturnPercent float64,
	yearlyLoanInterestPercent float64,
	loanAmount float64,
	downPayment float64,
	loanDurationYears float64,
) float64 {
	totalMoneyRemaining := loanAmount - downPayment
	yearlyInvestmentReturnRate := 1 + yearlyInvestmentReturnPercent/100.0
	monthlyEmi := calculateMonthlyEmi(totalMoneyRemaining, yearlyLoanInterestPercent, loanDurationYears)

	for i := 0.0; i < loanDurationYears; i++ {
		totalMoneyRemaining = totalMoneyRemaining*yearlyInvestmentReturnRate - monthlyEmi*12.0
	}
	return totalMoneyRemaining
}

func calculate_intialCapitalInvested_emiPaidFromPocket(
	yearlyInvestmentReturnPercent float64,
	yearlyLoanInterestPercent float64,
	loanAmount float64,
	downPayment float64,
	loanDurationYears float64,
) float64 {
	totalMoneyRemaining := loanAmount - downPayment
	yearlyInvestmentReturnRate := 1 + yearlyInvestmentReturnPercent/100.0
	monthlyEmi := calculateMonthlyEmi(totalMoneyRemaining, yearlyLoanInterestPercent, loanDurationYears)

	for i := 0.0; i < loanDurationYears; i++ {
		totalMoneyRemaining = totalMoneyRemaining * yearlyInvestmentReturnRate
	}
	return totalMoneyRemaining - 12*loanDurationYears*monthlyEmi
}

func calculate_boughtFromInitalCapital_emiAmountInvestedInSip(
	yearlyInvestmentReturnPercent float64,
	yearlyLoanInterestPercent float64,
	loanAmount float64,
	loanDurationYears float64,
) float64 {
	totalMoneyRemaining := loanAmount
	monthlyEmi := calculateMonthlyEmi(totalMoneyRemaining, yearlyLoanInterestPercent, loanDurationYears)
	return calculateSipReturns(monthlyEmi, yearlyInvestmentReturnPercent, loanDurationYears) - totalMoneyRemaining
}

func printStrategyReturns(
	yearlyInvestmentReturnPercent float64,
	yearlyLoanInterestPercent float64,
	loanAmount float64,
	downPayment float64,
	loanDurationYears float64,
	tabs int,
) {
	l := en_IN.New()
	pad := padding(tabs)

	amountRemaining := calculate_intialCapitalInvested_liquidatedMonthlyToPayEmi(
		yearlyInvestmentReturnPercent,
		yearlyLoanInterestPercent,
		loanAmount,
		downPayment,
		loanDurationYears,
	)
	fmt.Printf("%sintialCapitalInvested_liquidatedMonthlyToPayEmi: %v\n", pad, l.FmtCurrency(amountRemaining, 0, currency.INR))

	amountRemaining = calculate_intialCapitalInvested_liquidatedYearlyToPayEmi(
		yearlyInvestmentReturnPercent,
		yearlyLoanInterestPercent,
		loanAmount,
		downPayment,
		loanDurationYears,
	)
	fmt.Printf("%sintialCapitalInvested_liquidatedYearlyToPayEmi: %v\n", pad, l.FmtCurrency(amountRemaining, 0, currency.INR))

	amountRemaining = calculate_intialCapitalInvested_emiPaidFromPocket(
		yearlyInvestmentReturnPercent,
		yearlyLoanInterestPercent,
		loanAmount,
		downPayment,
		loanDurationYears,
	)
	fmt.Printf("%sintialCapitalInvested_emiPaidFromPocket: %v\n", pad, l.FmtCurrency(amountRemaining, 0, currency.INR))

	amountRemaining = calculate_boughtFromInitalCapital_emiAmountInvestedInSip(
		yearlyInvestmentReturnPercent,
		yearlyLoanInterestPercent,
		loanAmount,
		loanDurationYears,
	)
	fmt.Printf("%sboughtFromInitialCapital_emiAmountInvestedInSip: %v\n", pad, l.FmtCurrency(amountRemaining, 0, currency.INR))
}

func printReturnsForValues() {
	var (
		yearlyInvestmentReturnPercent = 11.0
		yearlyLoanInterestPercent     = 8.0
		loanAmount                    = 20_00_000.0
		downPayment                   = 0.0
		loanDurationYears             = 7.0
	)
	printStrategyReturns(
		yearlyInvestmentReturnPercent,
		yearlyLoanInterestPercent,
		loanAmount,
		downPayment,
		loanDurationYears,
		0,
	)
}

func optimizeReturns() {
	var (
		yearlyInvestmentReturnPercent = 11.0
		yearlyLoanInterestPercent     = 8.0
		loanAmount                    = 20_00_000.0
		downPayment                   = 0.0
	)

	for loanDurationYears := 3.0; loanDurationYears <= 10.0; loanDurationYears++ {
		fmt.Printf("\nLoan Duration Years: %v\n", loanDurationYears)
		printStrategyReturns(yearlyInvestmentReturnPercent,
			yearlyLoanInterestPercent,
			loanAmount,
			downPayment,
			loanDurationYears,
			1,
		)
	}
}

func main() {
	printReturnsForValues()
}

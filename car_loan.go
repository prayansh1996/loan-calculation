package main

import (
	"fmt"
	"math"

	"github.com/go-playground/locales/currency"
	"github.com/go-playground/locales/en_IN"
)

const (
	yearlyInvestmentReturnPercent = 11.0
	yearlyLoanInterestPercent     = 8.0
	loanAmount                    = 20_00_000.0
	loanDurationYears             = 7.0
)

func calculateMonthlyInterestRate(yearlyInterestPercent float64) float64 {
	return math.Pow(1.0+yearlyInterestPercent/100.0, 1.0/12.0)
}

func calculateMonthlyEmi(loanAmountLocal float64, yearlyLoanInterestPercentLocal float64, loanDurationYearsLocal float64) float64 {
	monthlyLoanInterestRate := yearlyLoanInterestPercentLocal / (12.0 * 100.0)
	months := loanDurationYearsLocal * 12.0
	exponent := math.Pow(1.0+monthlyLoanInterestRate, months)
	return loanAmountLocal * monthlyLoanInterestRate * exponent / (exponent - 1.0)
}

func calculateSipReturns(monthlySipAmount float64, yearlyInvestmentReturnPercentLocal float64, investmentDurationYears float64) float64 {
	monthlyInterestRate := yearlyInvestmentReturnPercentLocal / (12.0 * 100.0)
	numerator := monthlySipAmount * (math.Pow(1.0+monthlyInterestRate, 12.0*investmentDurationYears) - 1.0) * (1.0 + monthlyInterestRate)
	return numerator / monthlyInterestRate
}

func calculate_intialCapitalInvested_liquidatedMonthlyToPayEmi() float64 {
	totalMoneyRemaining := loanAmount
	monthlyInvestmentReturnRate := calculateMonthlyInterestRate(yearlyInvestmentReturnPercent)
	monthlyEmi := calculateMonthlyEmi(loanAmount, yearlyLoanInterestPercent, loanDurationYears)

	for i := 0; i < 12.0*loanDurationYears; i++ {
		totalMoneyRemaining = monthlyInvestmentReturnRate*totalMoneyRemaining - monthlyEmi
	}
	return totalMoneyRemaining
}

func calculate_intialCapitalInvested_liquidatedYearlyToPayEmi() float64 {
	totalMoneyRemaining := loanAmount
	yearylyInvestmentReturnRate := 1 + yearlyInvestmentReturnPercent/100.0
	monthlyEmi := calculateMonthlyEmi(loanAmount, yearlyLoanInterestPercent, loanDurationYears)

	for i := 0; i < loanDurationYears; i++ {
		totalMoneyRemaining = totalMoneyRemaining*yearylyInvestmentReturnRate - monthlyEmi*12.0
	}
	return totalMoneyRemaining
}

func calculate_intialCapitalInvested_emiPaidFromPocket() float64 {
	totalMoneyRemaining := loanAmount
	yearylyInvestmentReturnRate := 1 + yearlyInvestmentReturnPercent/100.0
	monthlyEmi := calculateMonthlyEmi(loanAmount, yearlyLoanInterestPercent, loanDurationYears)

	for i := 0; i < loanDurationYears; i++ {
		totalMoneyRemaining = totalMoneyRemaining * yearylyInvestmentReturnRate
	}
	return totalMoneyRemaining - 12*loanDurationYears*monthlyEmi
}

func calculate_boughtFromInitalCapital_emiAmountInvestedInSip() float64 {
	monthlyEmi := calculateMonthlyEmi(loanAmount, yearlyLoanInterestPercent, loanDurationYears)
	return calculateSipReturns(monthlyEmi, yearlyInvestmentReturnPercent, loanDurationYears) - loanAmount
}

func main() {
	l := en_IN.New()

	amountRemaining := calculate_intialCapitalInvested_liquidatedMonthlyToPayEmi()
	fmt.Printf("intialCapitalInvested_liquidatedMonthlyToPayEmi: %v\n", l.FmtCurrency(amountRemaining, 0, currency.INR))

	amountRemaining = calculate_intialCapitalInvested_liquidatedYearlyToPayEmi()
	fmt.Printf("intialCapitalInvested_liquidatedYearlyToPayEmi: %v\n", l.FmtCurrency(amountRemaining, 0, currency.INR))

	amountRemaining = calculate_intialCapitalInvested_emiPaidFromPocket()
	fmt.Printf("intialCapitalInvested_rmiPaidFromPocket: %v\n", l.FmtCurrency(amountRemaining, 0, currency.INR))

	amountRemaining = calculate_boughtFromInitalCapital_emiAmountInvestedInSip()
	fmt.Printf("boughtFromInitialCapital_emiAmountInvestedInSip: %v\n", l.FmtCurrency(amountRemaining, 0, currency.INR))
}

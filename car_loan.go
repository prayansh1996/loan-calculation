package main

import (
	"fmt"
	"math"

	"github.com/go-playground/locales/currency"
	"github.com/go-playground/locales/en_IN"
)

type Calculation struct {
	strategy     string
	finalCapital float64
}

func (c Calculation) String() string {
	l := en_IN.New()
	return fmt.Sprintf("%s: %v\n", c.strategy, l.FmtCurrency(c.finalCapital, 0, currency.INR))
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
) Calculation {
	totalMoneyRemaining := loanAmount - downPayment
	monthlyInvestmentReturnRate := calculateMonthlyInterestRate(yearlyInvestmentReturnPercent)
	monthlyEmi := calculateMonthlyEmi(totalMoneyRemaining, yearlyLoanInterestPercent, loanDurationYears)

	for i := 0.0; i < 12.0*loanDurationYears; i++ {
		totalMoneyRemaining = monthlyInvestmentReturnRate*totalMoneyRemaining - monthlyEmi
	}
	return Calculation{"Initial capital invested in market, liquidated monthly to pay EMI", totalMoneyRemaining}
}

func calculate_intialCapitalInvested_liquidatedYearlyToPayEmi(
	yearlyInvestmentReturnPercent float64,
	yearlyLoanInterestPercent float64,
	loanAmount float64,
	downPayment float64,
	loanDurationYears float64,
) Calculation {
	totalMoneyRemaining := loanAmount - downPayment
	yearlyInvestmentReturnRate := 1 + yearlyInvestmentReturnPercent/100.0
	monthlyEmi := calculateMonthlyEmi(totalMoneyRemaining, yearlyLoanInterestPercent, loanDurationYears)

	for i := 0.0; i < loanDurationYears; i++ {
		totalMoneyRemaining = totalMoneyRemaining*yearlyInvestmentReturnRate - monthlyEmi*12.0
	}
	return Calculation{"Initial capital invested in market, liquidated yearly to pay EMI", totalMoneyRemaining}
}

func calculate_intialCapitalInvested_emiPaidFromPocket(
	yearlyInvestmentReturnPercent float64,
	yearlyLoanInterestPercent float64,
	loanAmount float64,
	downPayment float64,
	loanDurationYears float64,
) Calculation {
	totalMoneyRemaining := loanAmount - downPayment
	yearlyInvestmentReturnRate := 1 + yearlyInvestmentReturnPercent/100.0
	monthlyEmi := calculateMonthlyEmi(totalMoneyRemaining, yearlyLoanInterestPercent, loanDurationYears)

	for i := 0.0; i < loanDurationYears; i++ {
		totalMoneyRemaining = totalMoneyRemaining * yearlyInvestmentReturnRate
	}
	return Calculation{"Initial capital invested in market, EMI paid from pocket", totalMoneyRemaining - 12*loanDurationYears*monthlyEmi}
}

func calculate_boughtFromInitalCapital_emiAmountInvestedInSip(
	yearlyInvestmentReturnPercent float64,
	yearlyLoanInterestPercent float64,
	loanAmount float64,
	downPayment float64,
	loanDurationYears float64,
) Calculation {
	totalMoneyRemaining := loanAmount
	monthlyEmi := calculateMonthlyEmi(totalMoneyRemaining, yearlyLoanInterestPercent, loanDurationYears)
	return Calculation{"Bought from initial capital, EMI amount invested in SIP", calculateSipReturns(monthlyEmi, yearlyInvestmentReturnPercent, loanDurationYears) - totalMoneyRemaining}
}

func calculateStrategyReturns(
	yearlyInvestmentReturnPercent float64,
	yearlyLoanInterestPercent float64,
	loanAmount float64,
	downPayment float64,
	loanDurationYears float64,
) []Calculation {
	calculations := []Calculation{}

	calculations = append(
		calculations,
		calculate_intialCapitalInvested_liquidatedMonthlyToPayEmi(
			yearlyInvestmentReturnPercent,
			yearlyLoanInterestPercent,
			loanAmount,
			downPayment,
			loanDurationYears,
		),
		calculate_intialCapitalInvested_liquidatedYearlyToPayEmi(
			yearlyInvestmentReturnPercent,
			yearlyLoanInterestPercent,
			loanAmount,
			downPayment,
			loanDurationYears,
		),
		calculate_intialCapitalInvested_emiPaidFromPocket(
			yearlyInvestmentReturnPercent,
			yearlyLoanInterestPercent,
			loanAmount,
			downPayment,
			loanDurationYears,
		),
		calculate_boughtFromInitalCapital_emiAmountInvestedInSip(
			yearlyInvestmentReturnPercent,
			yearlyLoanInterestPercent,
			loanAmount,
			downPayment,
			loanDurationYears,
		),
	)
	return calculations
}

func printReturnsForFixedValues() {
	var (
		yearlyInvestmentReturnPercent = 11.0
		yearlyLoanInterestPercent     = 8.0
		loanAmount                    = 20_00_000.0
		downPayment                   = 0.0
		loanDurationYears             = 7.0
	)
	calculations := calculateStrategyReturns(
		yearlyInvestmentReturnPercent,
		yearlyLoanInterestPercent,
		loanAmount,
		downPayment,
		loanDurationYears,
	)

	for _, calculation := range calculations {
		fmt.Println(calculation)
	}
}

func optimizeReturnsForDurationAndDownPayment() {
	var (
		yearlyInvestmentReturnPercent = 11.0
		yearlyLoanInterestPercent     = 8.0
		loanAmount                    = 20_00_000.0
	)

	for loanDurationYears := 3.0; loanDurationYears <= 10.0; loanDurationYears++ {
		fmt.Printf("\nLoan Duration Years: %v\n", loanDurationYears)

		for downPayment := 0.0; downPayment <= 7_00_000; downPayment += 1_00_000 {
			fmt.Printf("\n\tDown Payment: %v\n", downPayment)

			calculations := calculateStrategyReturns(
				yearlyInvestmentReturnPercent,
				yearlyLoanInterestPercent,
				loanAmount,
				downPayment,
				loanDurationYears,
			)
			for _, calculation := range calculations {
				fmt.Printf("\t\t")
				fmt.Println(calculation)
			}
		}
	}
}

func main() {
	optimizeReturnsForDurationAndDownPayment()
}

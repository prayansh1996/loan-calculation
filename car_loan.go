package main

import (
	"fmt"
	"math"
	"sort"

	"github.com/go-playground/locales/currency"
	"github.com/go-playground/locales/en_IN"
)

const (
	CAPITAL_INVESTED_LIQUIDATED_MONTHLY = "Initial capital invested in market, liquidated monthly to pay EMI"
	CAPITAL_INVESTED_LIQUIDATED_YEARLY  = "Initial capital invested in market, liquidated yearly to pay EMI"
	CAPITAL_INVESTED_EMI_FROM_POCKET    = "Initial capital invested in market, EMI paid from pocket"
	CASH_PAYMENT_EMI_TO_SIP             = "Bought from initial capital"
)

type Calculation struct {
	strategy         string
	capitalRemaining float64
	capitalSpent     float64

	yearlyInvestmentReturnPercent float64
	yearlyLoanInterestPercent     float64
	loanAmount                    float64
	downPayment                   float64
	loanDurationYears             float64
}

type CalculationWithVariables struct {
	calculation       Calculation
	downPayment       float64
	loanDurationYears float64
}

func (c Calculation) String() string {
	l := en_IN.New()
	return fmt.Sprintf(
		"\nStrategy: %s\nFinal Capital: %v\nCapital Spent: %v\nNet Loss: %v",
		c.strategy,
		l.FmtCurrency(c.capitalRemaining, 0, currency.INR),
		l.FmtCurrency(c.capitalSpent, 0, currency.INR),
		l.FmtCurrency(c.capitalSpent-c.capitalRemaining, 0, currency.INR),
	)
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
	totalMoneySpent := downPayment

	monthlyInvestmentReturnRate := calculateMonthlyInterestRate(yearlyInvestmentReturnPercent)
	monthlyEmi := calculateMonthlyEmi(totalMoneyRemaining, yearlyLoanInterestPercent, loanDurationYears)

	for i := 0.0; i < 12.0*loanDurationYears; i++ {
		totalMoneyRemaining = monthlyInvestmentReturnRate*totalMoneyRemaining - monthlyEmi
		totalMoneySpent += monthlyEmi
	}

	return Calculation{
		strategy:                      CAPITAL_INVESTED_LIQUIDATED_MONTHLY,
		capitalRemaining:              totalMoneyRemaining,
		capitalSpent:                  totalMoneySpent,
		yearlyInvestmentReturnPercent: yearlyInvestmentReturnPercent,
		yearlyLoanInterestPercent:     yearlyLoanInterestPercent,
		loanAmount:                    loanAmount,
		downPayment:                   downPayment,
		loanDurationYears:             loanDurationYears,
	}
}

func calculate_intialCapitalInvested_liquidatedYearlyToPayEmi(
	yearlyInvestmentReturnPercent float64,
	yearlyLoanInterestPercent float64,
	loanAmount float64,
	downPayment float64,
	loanDurationYears float64,
) Calculation {
	totalMoneyRemaining := loanAmount - downPayment
	totalMoneySpent := downPayment

	yearlyInvestmentReturnRate := 1 + yearlyInvestmentReturnPercent/100.0
	monthlyEmi := calculateMonthlyEmi(totalMoneyRemaining, yearlyLoanInterestPercent, loanDurationYears)

	for i := 0.0; i < loanDurationYears; i++ {
		totalMoneyRemaining = totalMoneyRemaining*yearlyInvestmentReturnRate - monthlyEmi*12.0
		totalMoneySpent += monthlyEmi * 12
	}

	return Calculation{
		strategy:                      CAPITAL_INVESTED_LIQUIDATED_YEARLY,
		capitalRemaining:              totalMoneyRemaining,
		capitalSpent:                  totalMoneySpent,
		yearlyInvestmentReturnPercent: yearlyInvestmentReturnPercent,
		yearlyLoanInterestPercent:     yearlyLoanInterestPercent,
		loanAmount:                    loanAmount,
		downPayment:                   downPayment,
		loanDurationYears:             loanDurationYears,
	}
}

func calculate_intialCapitalInvested_emiPaidFromPocket(
	yearlyInvestmentReturnPercent float64,
	yearlyLoanInterestPercent float64,
	loanAmount float64,
	downPayment float64,
	loanDurationYears float64,
) Calculation {
	totalMoneyRemaining := loanAmount - downPayment
	totalMoneySpent := downPayment

	yearlyInvestmentReturnRate := 1 + yearlyInvestmentReturnPercent/100.0
	monthlyEmi := calculateMonthlyEmi(totalMoneyRemaining, yearlyLoanInterestPercent, loanDurationYears)

	for i := 0.0; i < loanDurationYears; i++ {
		totalMoneyRemaining = totalMoneyRemaining * yearlyInvestmentReturnRate
	}

	totalMoneyRemaining = totalMoneyRemaining - 12*loanDurationYears*monthlyEmi
	totalMoneySpent += 12 * loanDurationYears * monthlyEmi

	return Calculation{
		strategy:                      CAPITAL_INVESTED_EMI_FROM_POCKET,
		capitalRemaining:              totalMoneyRemaining,
		capitalSpent:                  totalMoneySpent,
		yearlyInvestmentReturnPercent: yearlyInvestmentReturnPercent,
		yearlyLoanInterestPercent:     yearlyLoanInterestPercent,
		loanAmount:                    loanAmount,
		downPayment:                   downPayment,
		loanDurationYears:             loanDurationYears,
	}
}

func calculate_boughtFromInitalCapital_emiAmountInvestedInSip(
	yearlyInvestmentReturnPercent float64,
	yearlyLoanInterestPercent float64,
	loanAmount float64,
	downPayment float64,
	loanDurationYears float64,
) Calculation {
	monthlyEmi := calculateMonthlyEmi(loanAmount, yearlyLoanInterestPercent, loanDurationYears)
	totalMoneyRemaining := calculateSipReturns(monthlyEmi, yearlyInvestmentReturnPercent, loanDurationYears) - loanAmount
	totalMoneySpent := loanAmount + monthlyEmi*12*loanDurationYears

	return Calculation{
		strategy:                      CASH_PAYMENT_EMI_TO_SIP,
		capitalRemaining:              totalMoneyRemaining,
		capitalSpent:                  totalMoneySpent,
		yearlyInvestmentReturnPercent: yearlyInvestmentReturnPercent,
		yearlyLoanInterestPercent:     yearlyLoanInterestPercent,
		loanAmount:                    loanAmount,
		downPayment:                   downPayment,
		loanDurationYears:             loanDurationYears,
	}
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

// --------------- STRATEGY TESTS -------------- //

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
		yearlyLoanInterestPercent     = 8.75
		loanAmount                    = 21_50_000.0
	)

	calculations := []Calculation{}
	for loanDurationYears := 1.0; loanDurationYears <= 7.0; loanDurationYears++ {
		for downPayment := 0.0; downPayment <= 10_00_000; downPayment += 1_00_000 {
			calculationForVariables := calculateStrategyReturns(
				yearlyInvestmentReturnPercent,
				yearlyLoanInterestPercent,
				loanAmount,
				downPayment,
				loanDurationYears,
			)
			calculations = append(calculations, calculationForVariables...)
		}
	}

	sort.Slice(calculations, func(l, r int) bool {
		return (calculations[l].capitalRemaining - calculations[l].capitalSpent) >
			(calculations[r].capitalRemaining - calculations[r].capitalSpent)
	})

	l := en_IN.New()
	for i := 0; i < 10; i++ {
		fmt.Printf("\n\t\tLoan Duration: %v, Down Payment: %v\n", calculations[i].loanDurationYears, l.FmtCurrency(calculations[i].downPayment, 0, currency.INR))
		fmt.Println(calculations[i])
	}
}

func printReturnsForFixedDurationAndDownPayment() {
	var (
		yearlyInvestmentReturnPercent = 12.0
		yearlyLoanInterestPercent     = 9.0
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

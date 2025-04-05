package domain

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type ExpenseCategory string

const (
	ExpenseCategoryCleaningSupplies ExpenseCategory = "Środki czystości"
	ExpenseCategoryGroceries        ExpenseCategory = "Spożywcze"
	ExpenseCategoryClothing         ExpenseCategory = "Odzież i obuwie"
	ExpenseCategoryEntertainment    ExpenseCategory = "Rozrywka"
	ExpenseCategoryDiningOut        ExpenseCategory = "Gastronomia"
	ExpenseCategoryCosmetics        ExpenseCategory = "Kosmetyki"
	ExpenseCategoryEducation        ExpenseCategory = "Edukacja"
	ExpenseCategoryTransport        ExpenseCategory = "Transport"
	ExpenseCategoryHealth           ExpenseCategory = "Zdrowie"
	ExpenseCategoryEquipment        ExpenseCategory = "Wyposażenie"
	ExpenseCategoryLeisureTravel    ExpenseCategory = "Wypoczynek i podróże"
	ExpenseCategoryRentAndUtilities ExpenseCategory = "Czynsz i media"
	ExpenseCategoryElectronics      ExpenseCategory = "Elektronika"
	ExpenseCategoryGifts            ExpenseCategory = "Prezenty"
	ExpenseCategoryInsurance        ExpenseCategory = "Ubezpieczenia"
	ExpenseCategoryCar              ExpenseCategory = "Samochód"
	ExpenseCategoryWedding          ExpenseCategory = "Ślub"
	ExpenseCategoryHygieneProducts  ExpenseCategory = "Artykuły higieniczne"
	ExpenseCategoryOtherServices    ExpenseCategory = "Inne usługi"
	// Add more categories as needed...
)

var ExpenseCategories = []ExpenseCategory{
	ExpenseCategoryCleaningSupplies,
	ExpenseCategoryGroceries,
	ExpenseCategoryClothing,
	ExpenseCategoryEntertainment,
	ExpenseCategoryDiningOut,
	ExpenseCategoryCosmetics,
	ExpenseCategoryEducation,
	ExpenseCategoryTransport,
	ExpenseCategoryHealth,
	ExpenseCategoryEquipment,
	ExpenseCategoryLeisureTravel,
	ExpenseCategoryRentAndUtilities,
	ExpenseCategoryElectronics,
	ExpenseCategoryGifts,
	ExpenseCategoryInsurance,
	ExpenseCategoryCar,
	ExpenseCategoryWedding,
	ExpenseCategoryHygieneProducts,
	ExpenseCategoryOtherServices,
}

type Expense struct {
	Timestamp time.Time
	Amount    float64
	Category  ExpenseCategory
}

func (e *Expense) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Timestamp string `json:"timestamp"`
		Amount    string `json:"amount"`
		Category  string `json:"category"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	timestamp, err := time.Parse(time.DateOnly, aux.Timestamp)
	if err != nil {
		return err
	}
	e.Timestamp = timestamp
	amount, err := strconv.ParseFloat(strings.TrimSpace(aux.Amount), 64)
	if err != nil {
		return err
	}
	e.Amount = amount
	e.Category = ExpenseCategory(aux.Category)
	return nil
}

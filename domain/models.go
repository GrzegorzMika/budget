package domain

import "time"

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
}

type Expense struct {
	Timestamp time.Time
	Amount    float64
	Category  ExpenseCategory
}

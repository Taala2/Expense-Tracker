package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type ExpenseTracker struct {
	Id          int       `json:"id"`
	Date        time.Time `json:"date"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
}

type ExpenseTrackersStorage struct {
	Trackers []ExpenseTracker `json:"expanse_tracker"`
	NextId   int              `json:"next_id"`
}

var expensesTrackersStorage = ExpenseTrackersStorage{NextId: 1}

func loadExpensesTracker(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&expensesTrackersStorage)
	if err != nil && err != os.ErrNotExist && err != io.EOF {
		return err
	}
	return nil
}

func saveExpense(f string) error {
	file, err := os.Create(f)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	return encoder.Encode(expensesTrackersStorage)
}

func addExpense(description string, amount float64) {
	exp := ExpenseTracker{
		Id:          expensesTrackersStorage.NextId,
		Date:        time.Now(),
		Amount:      amount,
		Description: description,
	}
	expensesTrackersStorage.NextId++
	expensesTrackersStorage.Trackers = append(expensesTrackersStorage.Trackers, exp)

	fmt.Printf("Expense added: %d", exp.Id)
}

func expensesList() {
	fmt.Println("ID	Date		Amount	Description")

	for _, exp := range expensesTrackersStorage.Trackers {
		fmt.Printf("%d	%s	%.2f	%s\n", exp.Id, exp.Date.Format("2006-01-02"), exp.Amount, exp.Description)
	}
}

func summary() {
	sum := 0.0
	for _, exp := range expensesTrackersStorage.Trackers {
		sum += exp.Amount
	}

	fmt.Printf("Summary: %.2f\n", sum)
}

func monthSummary(mouth string) {
	for _, exp := range expensesTrackersStorage.Trackers {
		if exp.Date.Format("2006-01") == mouth {
			fmt.Printf("%s	%.2f\n", exp.Date.Format("2006-01"), exp.Amount)
		}
	}
}

func deleteExpense(id string) {
	ida, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error parsing id: ", err)
		return
	}
	for i, exp := range expensesTrackersStorage.Trackers {
		if exp.Id == ida {
			expensesTrackersStorage.Trackers = append(expensesTrackersStorage.Trackers[:i], expensesTrackersStorage.Trackers[i+1:]...)
			fmt.Printf("Expense deleted: %d\n", ida)
			return
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: help")
	}

	arg := os.Args[1]

	filepath := "data.json"

	err := loadExpensesTracker(filepath)
	if err != nil {
		fmt.Println("Error loading expenses tracker: ", err)
	}

	switch arg {
	case "add":
		if len(os.Args) < 2 {
			fmt.Println("Usage: add <description> <amount>")
		}
		amount, err := strconv.ParseFloat(os.Args[3], 64)
		if err != nil {
			fmt.Println("Error parsing amount: ", err)
		}
		addExpense(os.Args[2], amount)
		saveExpense(filepath)
	case "list":
		expensesList()
	case "summary":
		summary()
	case "month-summary":
		if len(os.Args) > 2 {
			monthSummary(os.Args[2])
		} else if os.Args[2] == "--mouth" {
			fmt.Println("Usage: month-summary <mouth>")
		}
	case "delete":
		if len(os.Args) == 3 {
			deleteExpense(os.Args[2])
			saveExpense(filepath)
		} else {
			fmt.Println("Usage: delete <id>")
		}
	default:
		fmt.Println("Usage: help")
	}
}

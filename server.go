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

func addExpense(description string, amount float64) error{
	if amount < 0 {
		return fmt.Errorf("invalid amount")
	}

	exp := ExpenseTracker{
		Id:          expensesTrackersStorage.NextId,
		Date:        time.Now(),
		Amount:      amount,
		Description: description,
	}
	expensesTrackersStorage.NextId++
	expensesTrackersStorage.Trackers = append(expensesTrackersStorage.Trackers, exp)

	fmt.Printf("Expense added: %d", exp.Id)
	return nil
}

func expensesList() {
	fmt.Println("ID	Date		Amount	Description")

	for _, exp := range expensesTrackersStorage.Trackers {
		fmt.Printf("%d	%s	$%.2f	%s\n", exp.Id, exp.Date.Format("2006-01-02"), exp.Amount, exp.Description)
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

			for j := i; j < len(expensesTrackersStorage.Trackers); j++ {
				expensesTrackersStorage.Trackers[j].Id--
			}
	
			expensesTrackersStorage.NextId--

			return
		}
	}
}

func updateExpense(id string, arg string, change string) error {
	ida, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	if arg == "--d" {
		for i, exp := range expensesTrackersStorage.Trackers {
			if exp.Id == ida {
				expensesTrackersStorage.Trackers[i].Description = change
			} else {
				return fmt.Errorf("expense not found")
			}
		}
	} else if arg == "--a" {
		change, err := strconv.ParseFloat(change, 64)
		if err != nil {
			return err
		}

		for i, exp := range expensesTrackersStorage.Trackers {
			if exp.Id == ida {
				expensesTrackersStorage.Trackers[i].Amount = change
			} else {
				return fmt.Errorf("expense not found")
			}
		}
	} else {
		return err
	}

	fmt.Printf("Expense updated")

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: help")
		os.Exit(1)
	}
	arg := os.Args[1]

	filepath := "data.json"

	err := loadExpensesTracker(filepath)
	if err != nil {
		fmt.Println("Error loading expenses tracker: ", err)
	}

	switch arg {
	case "add":
		if len(os.Args) < 2 || len(os.Args) > 4 {
			fmt.Println("Usage: add <description> <amount>")
		}
		amount, err := strconv.ParseFloat(os.Args[3], 64)
		if err != nil {
			fmt.Println("Error parsing amount")
		}
		err =addExpense(os.Args[2], amount)
		if err != nil {
			fmt.Println("Error adding expense: ", err)
		}
		saveExpense(filepath)
	case "list":
		expensesList()
	case "summary":
		summary()
	case "month-summary":
		if len(os.Args) < 2 {
			monthSummary("")
			fmt.Println("Usage: month-summary --month <mouth> for example: month-summary --month 8")
		} else if os.Args[2] == "--mouth" {
			monthSummary(os.Args[3])
		}else{
			fmt.Println("Usage: month-summary --month <mouth> for example: month-summary --month 8")}
	case "delete":
		if len(os.Args) == 3 {
			deleteExpense(os.Args[2])
			saveExpense(filepath)
		} else {
			fmt.Println("Usage: delete <id>")
		}
	case "update":
		if len(os.Args) == 5 {
			err = updateExpense(os.Args[2], os.Args[3], os.Args[4])
			if err != nil {
				fmt.Println("Error updating expense: ", err)
			}
			saveExpense(filepath)
		} else {
			fmt.Println("Usage: update <id> --d or --a <argument>")
		}
	default:
		fmt.Println("Usage: help")
	}
}

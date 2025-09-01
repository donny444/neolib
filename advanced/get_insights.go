package advanced

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"neolib/database"
	"net/http"
	"time"
)

func TopFiveCategories(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TopFiveCategories function get called")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	username, ok := r.Context().Value("username").(string)
	if !ok || username == "" {
		http.Error(w, "Internal Server Wrror", http.StatusInternalServerError)
		log.Fatal("Username not found in context")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := database.SelectTopFive(ctx, username)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal("Unable to get the top five categories: ", err)
		return
	}
	defer rows.Close()

	type TopCategory struct {
		Category string `json:"category"`
		Count    int    `json:"count"`
	}
	var topCategory TopCategory
	var topCategories []TopCategory

	for rows.Next() {
		if err := rows.Scan(&topCategory.Category, &topCategory.Count); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Fatal("Unable to scan the row: ", err)
			return
		}

		topCategories = append(topCategories, topCategory)
	}

	jsonResponse, err := json.Marshal(topCategories)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal("Unable to marshal response to JSON: ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	fmt.Println("Succussfully retrieved top five categories in the bookshelf")
}

func ReadingStatusByCategory(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ReadingStatusByCategory function get called")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	username, ok := r.Context().Value("username").(string)
	if !ok || username == "" {
		http.Error(w, "Internal Server Wrror", http.StatusInternalServerError)
		log.Fatal("Username not found in context")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := database.SelectCategoryStatuses(ctx, username)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal("Unable to get counts of status by category: ", err)
		return
	}
	defer rows.Close()

	type CountOfStatus struct {
		Category string `json:"category"`
		Status   string `json:"status"`
		Count    int    `json:"count"`
	}
	var countOfStatus CountOfStatus
	var countsOfStatus []CountOfStatus

	for rows.Next() {
		if err := rows.Scan(&countOfStatus.Category, &countOfStatus.Status, &countOfStatus.Count); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Fatal("Unable to scan the row: ", err)
			return
		}

		countsOfStatus = append(countsOfStatus, countOfStatus)
	}

	jsonResponse, err := json.Marshal(countsOfStatus)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal("Unable to marshal response to JSON: ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	fmt.Println("Succussfully retrieved counts of status by category in the bookshelf")
}

func BooksByPages(w http.ResponseWriter, r *http.Request) {
	fmt.Println("BooksByPages function get called")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	username, ok := r.Context().Value("username").(string)
	if !ok || username == "" {
		http.Error(w, "Internal Server Wrror", http.StatusInternalServerError)
		log.Fatal("Username not found in context")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := database.BookGroupByPages(ctx, username)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal("Unable to get the books grouped by pages: ", err)
		return
	}
	defer rows.Close()

	type PageGroup struct {
		PageRange string `json:"page_range"`
		Count     int    `json:"count"`
	}
	var pageGroup PageGroup
	var pageGroups []PageGroup

	for rows.Next() {
		if err := rows.Scan(&pageGroup.PageRange, &pageGroup.Count); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Fatal("Unable to scan the row: ", err)
			return
		}

		pageGroups = append(pageGroups, pageGroup)
	}

	jsonResponse, err := json.Marshal(pageGroups)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal("Unable to marshal response to JSON: ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	fmt.Println("Succussfully retrieved the books grouped by pages in the bookshelf")
}

/*
func ReadBooksByMonth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ReadBooksByMonth function get called")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	username, ok := r.Context().Value("username").(string)
	if !ok || username == "" {
		http.Error(w, "Internal Server Wrror", http.StatusInternalServerError)
		log.Fatal("Username not found in context")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := database.SelectBooksByMonth(ctx, username)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal("Unable to get counts of read books by month: ", err)
		return
	}
	defer rows.Close()

	type CountByMonth struct {
		Month string `json:"month"`
		Count int    `json:"count"`
	}

	var countByMonth CountByMonth

	var countsByMonth []CountByMonth

	for rows.Next() {
		if err := rows.Scan(&countByMonth.Month, &countByMonth.Count); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Fatal("Unable to scan the row: ", err)
			return
		}

		countsByMonth = append(countsByMonth, countByMonth)
	}

	jsonResponse, err := json.Marshal(countsByMonth)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal("Unable to marshal response to JSON: ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	fmt.Println("Succussfully retrieved counts of read books by month in the bookshelf")
}
*/

package main

import (
	//"database/sql"
	//"encoding/json"

	"net/http"
	"strconv"

	//"time"

	_ "github.com/mattn/go-sqlite3"
)

// Обработчик DELETE-запроса
func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		respondWithError(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	err = deleteTaskFromDatabase(id)
	if err != nil {
		respondWithError(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, TaskResponse{Error: ""}, http.StatusOK)
}

func deleteTaskFromDatabase(id int) error {
	_, err := DB.Exec("DELETE FROM scheduler WHERE id=?", id)
	return err
}

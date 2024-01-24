package main

import (
	"PROJECT/models"
	"PROJECT/storage"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func main() {
	
	//endpoint to create new user
	http.HandleFunc("/user/create", createUser)
	//endpoint for getting all users
	http.HandleFunc("/user/all", GetAllUsers)
	http.HandleFunc("/user/update", updateUser)
	http.HandleFunc("/user/delete", deleteUser)
	log.Println("server starting...")
	err := http.ListenAndServe("localhost:8081", nil)

	if err != nil {
		fmt.Println("Error while running server: ", err)
	}

}

func createUser(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println("error while getting body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user models.User

	err = json.Unmarshal(bodyBytes, &user)

	id, err := uuid.NewUUID()

	if err != nil {
		log.Println("error while unmarshalling body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user.Id = id.String()
	respUser, err := storage.CreateUser(&user)

	if err != nil {
		log.Println("error while creating user", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBody, err := json.Marshal(respUser)

	if err != nil {
		log.Println("error while marshalling to response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(respBody)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {

	page := r.URL.Query().Get("page")

	intPage, err := strconv.Atoi(page)

	if err != nil {
		log.Println("error while converting page is not int", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	limit := r.URL.Query().Get("limit")

	intLimit, err := strconv.Atoi(limit)

	if err != nil {
		log.Println("error while converting limit is not int", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := storage.GetAll(intPage,intLimit)

	if err != nil {
		log.Println("error while getting all users smth went wrong", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respBody, err := json.Marshal(users)

	if err != nil {
		log.Println("error while marshalling to response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)

}


func updateUser(w http.ResponseWriter, r *http.Request) {
	// Parse the user ID from the URL query parameter
	userID := r.URL.Query().Get("id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error while reading request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var updatedUser models.User
	err = json.Unmarshal(bodyBytes, &updatedUser)
	if err != nil {
		log.Println("Error while unmarshalling JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Call the storage package to update the user
	respUser, err := storage.UpdateUser(userID, &updatedUser)
	if err != nil {
		log.Println("Error while updating user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Marshal the response user to JSON
	respBody, err := json.Marshal(respUser)
	if err != nil {
		log.Println("Error while marshalling response to JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	// Parse the user ID from the URL query parameter
	userID := r.URL.Query().Get("id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Call the storage package to delete the user
	respUser, err := storage.DeleteUser(userID)
	if err != nil {
		log.Println("Error while deleting user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Marshal the response user to JSON
	respBody, err := json.Marshal(respUser)
	if err != nil {
		log.Println("Error while marshalling response to JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

type Profile struct {
	ID         string `json:"id"`
	ColorTheme string `json:"color-theme"`
	User       *User  `json:"user"`
}

type User struct {
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
	Age       string `json:"age"`
	Gender    string `json:"gender"`
}

//-------------------------------------------------------------------------------------------

// функция , которая читает файл , где данный о профиле записаны в виде json

func getDBInfo() ([]Profile, error) {
	var profiles []Profile
	file, err := os.OpenFile("./static/profileInfo.txt", os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		var profile Profile
		err := json.Unmarshal([]byte(line), &profile)
		if err != nil {
			log.Println("Error unmarshaling JSON:", err)
			continue
		}
		profiles = append(profiles, profile)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return profiles, nil
}

//----------------------------------------------------------------------------------------

func setDBInfo(profiles []Profile) error {
	file, err := os.OpenFile("./static/profileInfo.txt", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, item := range profiles {
		itemJson, err := json.Marshal(item)
		if err != nil {
			return err
		}
		if _, err := writer.WriteString(string(itemJson) + "\n"); err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}

//-----------------------------------------------------------------------------------------

func getProfiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	profiles, err := getDBInfo()
	if err != nil {
		http.Error(w, "Error reading profiles", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(profiles)
}

//------------------------------------------------------------------------------------------

func getProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	profiles, err := getDBInfo()
	if err != nil {
		http.Error(w, "Error reading profiles", http.StatusInternalServerError)
		return
	}
	for _, value := range profiles {
		if value.ID == params["id"] {
			json.NewEncoder(w).Encode(value)
		}
	}
}

// ----------------------------------------------------------------------------------------------
func createProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text")
	var newProfile Profile

	profiles, err := getDBInfo()
	if err != nil {
		http.Error(w, "Error reading profiles", http.StatusInternalServerError)
		return
	}

	err1 := json.NewDecoder(r.Body).Decode(&newProfile)
	if err1 != nil {
		log.Println(err1)
	}
	newProfile.ID = strconv.Itoa(rand.Intn(1000000000000))
	profiles = append(profiles, newProfile)

	errDB := setDBInfo(profiles)
	if errDB != nil {
		log.Println(errDB)
	}
	str := "Создание прошло успешно"
	w.Write([]byte(str))

}

//----------------------------------------------------------------------------------------------

func deleteProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text")
	params := mux.Vars(r)
	profiles, err := getDBInfo()
	if err != nil {
		http.Error(w, "Error reading profiles", http.StatusInternalServerError)
		return
	}

	for index, item := range profiles {
		if item.ID == params["id"] {
			profiles = append(profiles[:index], profiles[index+1:]...)
			errDB := setDBInfo(profiles)
			if errDB != nil {
				log.Println(errDB)
			}
			str := "Удаление прошло успешно"
			w.Write([]byte(str))
			return
		}
	}

	str := "Удаление не может быть выполнено"
	w.Write([]byte(str))

}

//--------------------------------------------------------------------------------------------

func updateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text")
	params := mux.Vars(r)
	profiles, err := getDBInfo()
	if err != nil {
		http.Error(w, "Error reading profiles", http.StatusInternalServerError)
		return
	}

	for index, item := range profiles {
		if item.ID == params["id"] {
			profiles = append(profiles[:index], profiles[index+1:]...)
			var newProfile Profile

			err1 := json.NewDecoder(r.Body).Decode(&newProfile)
			if err1 != nil {
				log.Println(err1)
			}

			newProfile.ID = params["id"]
			profiles = append(profiles, newProfile)

			errDB := setDBInfo(profiles)
			if errDB != nil {
				log.Println(errDB)
			}
			str := "Обновление прошло успешно"
			w.Write([]byte(str))
			return
		}
	}

	str := "Обновление не может быть выполнено"
	w.Write([]byte(str))

}

//----------------------------------------------------------------------------------------------

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/profiles", getProfiles).Methods(http.MethodGet)
	r.HandleFunc("/profiles/{id}", getProfile).Methods(http.MethodGet)
	r.HandleFunc("/profiles", createProfile).Methods(http.MethodPost)
	r.HandleFunc("/profiles/{id}", deleteProfile).Methods(http.MethodDelete)
	r.HandleFunc("/profiles/{id}", updateProfile).Methods(http.MethodPut)

	fmt.Println("Server starting at port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

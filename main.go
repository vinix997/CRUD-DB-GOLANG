package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"ws/entity"
	"ws/service"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const PORT = ":8080"

func dbConn() *sql.DB {
	db, err := sql.Open("mysql", "root:admin123@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		panic(err)
	}
	return db
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", greet)
	r.HandleFunc("/users", UserHandler)
	r.HandleFunc("/users/{Id}", UserHandler)
	r.HandleFunc("/users/{Id}", UserHandler)

	http.Handle("/", r)
	http.ListenAndServe(PORT, nil)
}

func greet(w http.ResponseWriter, r *http.Request) {
	msg := "Hello world"
	fmt.Fprint(w, msg)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var user entity.User
		if err := decoder.Decode(&user); err != nil {
			w.Write([]byte("error decoding json body"))
			return
		}

		userSvc := service.NewUserService()
		userTemp := userSvc.Register(&user)

		jsonData, _ := json.Marshal(userTemp)

		w.Write(jsonData)
	}
}
func UserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	params := mux.Vars(r)
	Id := params["Id"]
	if r.Method == "GET" {
		if Id != "" {
			tempId, _ := strconv.Atoi(Id)
			GetUserById(w, r, tempId)
		} else {
			GetAllUser(w, r)
		}
	}
	if r.Method == "DELETE" {
		tempId, _ := strconv.Atoi(Id)
		DeleteUser(w, r, tempId)
	}
	if r.Method == "POST" {
		AddUser(w, r)
	}
	if r.Method == "PUT" {
		tempId, _ := strconv.Atoi(Id)
		UpdateUser(w, r, tempId)
	}
}

func GetUserById(w http.ResponseWriter, r *http.Request, id int) {
	db := dbConn()
	var user entity.User
	query := "SELECT ID, USERNAME, PASSWORD, EMAIL, AGE FROM USER WHERE ID = ?"
	row := db.QueryRow(query, id)
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Age)
	if err != nil {
		panic(err)
	}
	db.Close()
	jsonData, _ := json.Marshal(user)
	w.Write(jsonData)
}

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	results := []entity.User{}
	data, err := db.Query("SELECT id, username, password, email, age FROM USER")
	if err != nil {
		panic(err)
	}
	for data.Next() {
		var user entity.User
		err := data.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Age)
		if err != nil {
			panic(err)
		}
		results = append(results, user)
	}
	defer db.Close()
	test, _ := json.Marshal(results)
	w.Write(test)
}

func DeleteUser(w http.ResponseWriter, r *http.Request, id int) {
	db := dbConn()
	query := "DELETE FROM USER WHERE id = ?"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		panic(err)
	}
	res.LastInsertId()
	res.RowsAffected()
	defer db.Close()
	w.Write([]byte("User deleted successfully"))
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user entity.User
	if err := decoder.Decode(&user); err != nil {
		w.Write([]byte("error decoding json body"))
		return
	} else {
		db := dbConn()
		query := "INSERT INTO USER (username, password, email, age, createdat, updatedat) VALUES(?,?,?,?,?,?)"
		ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelfunc()
		stmt, err := db.PrepareContext(ctx, query)
		res, err := stmt.ExecContext(ctx, user.Username, user.Password, user.Email, user.Age, time.Now(), nil)
		if err != nil {
			log.Fatal(err)
		}
		res.LastInsertId()
		res.RowsAffected()
		defer db.Close()
		w.Write([]byte("User added successfully"))
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request, id int) {
	decoder := json.NewDecoder(r.Body)
	var temp entity.User
	if err := decoder.Decode(&temp); err != nil {
		w.Write([]byte("error decoding json body"))
		return
	}
	db := dbConn()
	query := "update user set Username = ?, Password = ?, Email = ?, Age = ?, UpdatedAt = ? where Id = ?"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		fmt.Println("Error 1")
	}
	res, err := stmt.ExecContext(ctx, temp.Username, temp.Password, temp.Email, temp.Age, time.Now(), id)
	if err != nil {
		// log.Fatal(err)
		fmt.Println("Error 2")
	}
	res.LastInsertId()
	res.RowsAffected()
	defer db.Close()
	w.Write([]byte("User updated successfuly"))

}

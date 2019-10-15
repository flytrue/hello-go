package main

//example for error (twice post)
//curl -X POST http://localhost:8081/user/test/mailtest@mail.ru/1234
//Error create = UNIQUE constraint failed: users.iduser

//Order not stored = userid =1, product=1 count=10 (типа пользователь 1 заказал 10 л молока)
//curl -X POST curl -X POST http://localhost:8081/order/1/1/10
//

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type User struct {
	gorm.Model
	Name   string
	Email  string
	Iduser string `gorm:"type:varchar(20);unique_index"`
	Orders Order
}

type Order struct {
	gorm.Model
	UserID       string
	DateDelivery *time.Time
	ProductId    int
	ProductCount int
	Price        float64
}

func addOrder(w http.ResponseWriter, r *http.Request) {
	//добавление заказа
	var productidint int
	var productcountint int

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// switch r.Method {
	// case "GET":
	// 	http.ServeFile(w, r, "form.html")
	// case "POST":
	// 	// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
	// 	if err := r.ParseForm(); err != nil {
	// 		fmt.Fprintf(w, "ParseForm() err: %v", err)
	// 		return
	// 	}
	// 	fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
	// 	name := r.FormValue("name")
	// 	address := r.FormValue("address")
	// 	fmt.Fprintf(w, "Name = %s\n", name)
	// 	fmt.Fprintf(w, "Address = %s\n", address)
	// default:
	// 	fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	// }
	// return

	fmt.Println("New order for User Endpoint Hit")

	vars := mux.Vars(r)
	userid := vars["userid"]

	productid := vars["productid"]
	productcount := vars["productcount"]

	if productidint, err := strconv.Atoi(productid); err == nil {
		fmt.Printf("i=%d, type: %T\n", productidint, productidint)
	}

	if productcountint, err := strconv.Atoi(productcount); err == nil {
		fmt.Printf("i=%d, type: %T\n", productcountint, productcountint)
	}

	fmt.Println(userid)
	fmt.Println(productid, " = ", productidint)
	fmt.Println(productcount, " = ", productcountint)

	dbc := db.Create(&Order{UserID: userid, DateDelivery: nil, ProductId: productidint, ProductCount: productcountint})

	if dbc.Error != nil {
		fmt.Fprintf(w, "Error create = "+dbc.Error.Error())
		return
	}
	fmt.Fprintf(w, "New User Successfully Created")

}

func allUsers(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var users []User
	db.Find(&users)
	fmt.Println("{}", users)

	json.NewEncoder(w).Encode(users)
}

func allOrders(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get all orders Hit")
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var orders []Order
	db.Find(&orders)
	fmt.Println("{}", orders)

	json.NewEncoder(w).Encode(orders)
}

func newUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New User Endpoint Hit")

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	vars := mux.Vars(r)
	name := vars["name"]
	email := vars["email"]
	iduser := vars["iduser"]

	fmt.Println(name)
	fmt.Println(email)
	fmt.Println(iduser)

	dbc := db.Create(&User{Name: name, Email: email, Iduser: iduser})
	fmt.Println("start create user name")
	if dbc.Error != nil {
		fmt.Fprintf(w, "Error create = "+dbc.Error.Error())
		return
	}
	fmt.Fprintf(w, "New User Successfully Created")
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]

	var user User
	db.Where("name = ?", name).Find(&user)
	db.Delete(&user)

	fmt.Fprintf(w, "Successfully Deleted User")
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	email := vars["email"]
	/*iduser := vars["iduser"] */

	var user User
	db.Where("name = ?", name).Find(&user)

	user.Email = email

	db.Save(&user)
	fmt.Fprintf(w, "Successfully Updated User")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/users", allUsers).Methods("GET")
	myRouter.HandleFunc("/user/{name}", deleteUser).Methods("DELETE")
	myRouter.HandleFunc("/user/{name}/{email}/{iduser}", updateUser).Methods("PUT")
	myRouter.HandleFunc("/user/{name}/{email}/{iduser}", newUser).Methods("POST")
	//myRouter.HandleFunc("/orders", allOrders).Methods("GET")
	myRouter.HandleFunc("/orders", allOrders).Methods("POST")
	myRouter.HandleFunc("/order/{userid}/{productid}/{productcount}", addOrder)

	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func initialMigration() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Order{})
}

func main() {
	fmt.Println("Go ORM Tutorial")

	initialMigration()
	// Handle Subsequent requests
	handleRequests()
}

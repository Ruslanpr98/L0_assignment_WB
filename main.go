package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "l0_user"
	password = "l0_password"
	dbname   = "l0_database"
)

type orders_model struct {
	Order_uid          string `json:"order_uid"`
	Track_number       string `json:"track_number"`
	Entry_name         string `json:"entry_name"`
	Locale             string `json:"locale"`
	Internal_signature string `json:"internal_signature"`
	Customer_id        string `json:"customer_id"`
	Delivery_service   string `json:"delivery_service"`
	Shardkey           string `json:"shardkey"`
	Sm_id              int    `json:"sm_id"`
	Date_created       string `json:"date_created"`
	Oof_shard          string `json:"oof_shard"`
}

var memcache map[string]orders_model // map for memcache of Data

//func for showing all orders
func show_all_orders(w http.ResponseWriter, req *http.Request) {
	for _, value := range memcache {
		result, err := json.Marshal(value)
		if err != nil {
			log.Println("No data exists", err)
			http.Error(w, http.StatusText(500), 500)
		}

		io.WriteString(w, string(result))
	}
}

//func for showing order by id
func get_order_by_uid(w http.ResponseWriter, req *http.Request) {
	v := req.FormValue("id")
	for key, value := range memcache {
		if key == v {
			result, err := json.Marshal(value)
			if err != nil {
				log.Panicln("Error marshaling JSON", err)
				http.Error(w, http.StatusText(500), 500)
			}
			io.WriteString(w, string(result))
			return
		}

	}
	io.WriteString(w, "No result")
}

func main() {

	// Database connection and checking if connection is ok
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalln("Attempt to connect to Database is failed", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Panicln("Database is not responding", err)
		return
	}

	fmt.Println("Successfully connected!")

	// Initializing memcache
	memcache = make(map[string]orders_model)

	// Getting Data from Postgres to our memcache
	rows, err := db.Query("SELECT * FROM orders")
	if err != nil {
		log.Fatalln("No SQL data", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		order := orders_model{}
		err = rows.Scan(&order.Order_uid,
			&order.Track_number,
			&order.Entry_name,
			&order.Locale,
			&order.Internal_signature,
			&order.Customer_id,
			&order.Delivery_service,
			&order.Shardkey,
			&order.Sm_id,
			&order.Date_created,
			&order.Oof_shard,
		)
		if err != nil {
			log.Panic(err)
			return
		}
		memcache[order.Order_uid] = order
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return
	}

	// checking if Data from Postgres is received
	fmt.Println(memcache)

	// Connection to NATS and error check
	sc, err := stan.Connect("test-cluster", "subscriber", stan.NatsURL(stan.DefaultNatsURL))

	if err != nil {
		log.Panicln("Error connectiong to NATS", err.Error())
	}

	var new_orders orders_model

	//Subscribing to channel to get data and insert it to Database
	// Also checking for errors
	sub, err := sc.Subscribe("orders_model", func(m *stan.Msg) {
		if err != nil {
			log.Fatalln("Subscription failed", err.Error())
		}

		err = json.Unmarshal(m.Data, &new_orders)
		if err != nil {
			log.Panicln("Could not unmarshal content", err.Error())
			return
		}
		fmt.Println(new_orders)

		if new_orders.Order_uid == "" || len(new_orders.Order_uid) < 18 {
			log.Fatalln("Order_uid is incorrect, adding stopped", err.Error())
			return
		}

		insert_statement := "INSERT INTO orders VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);"
		_, err = db.Exec(insert_statement, new_orders.Order_uid,
			new_orders.Track_number,
			new_orders.Entry_name,
			new_orders.Locale,
			new_orders.Internal_signature,
			new_orders.Customer_id,
			new_orders.Delivery_service,
			new_orders.Shardkey,
			new_orders.Sm_id,
			new_orders.Date_created,
			new_orders.Oof_shard,
		)
		if err != nil {
			log.Panicln("Insert failed", err.Error())
			return
		}
		log.Println("Adding sucessful")
		memcache[new_orders.Order_uid] = new_orders

	})

	// Deferring subscription and connection close
	defer sc.Close()
	defer sub.Close()

	http.HandleFunc("/orders/", get_order_by_uid)
	http.HandleFunc("/orders", show_all_orders)

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatalln(err)
	}

}

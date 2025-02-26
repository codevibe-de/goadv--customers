package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"slices"
	//"fmt"
)

var db = make(map[string]string)

type Customer struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func CreateJsonFile(c []Customer){
	b, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create("/tmp/customers.json")
	if err != nil {
			log.Fatal(err)
	}
	defer f.Close()
	f.Write(b)
}

func GetAllCustomers() []Customer {
	emptyReturnArray := make([]Customer, 0, 0)
	content, err := ioutil.ReadFile("/tmp/customers.json")
	if err != nil {
			//log.Fatal("Error when opening file: ", err)
			CreateJsonFile(emptyReturnArray)
			return emptyReturnArray
	}

	// Now let's unmarshall the data into `payload`
	var customers []Customer
	err = json.Unmarshal(content, &customers)
	if err != nil {
			//log.Fatal("Error during Unmarshal(): ", err)
			CreateJsonFile(emptyReturnArray)
			return emptyReturnArray
	}
	// fancy business logic
	//var return_list =
	return customers
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	r.GET("/customers", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, GetAllCustomers())
	})

	r.GET("/customer/:phone", func(c *gin.Context) {
		phone := c.Params.ByName("phone")
		customers := GetAllCustomers()
		for _, customer := range customers {
			if customer.Phone == phone {
				c.JSON(http.StatusOK, customer)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"status": "no customer found"})
	})

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/customers", func(c *gin.Context) {
		phone := c.Params.ByName("phone")
		name := c.Params.ByName("name")

		customers := GetAllCustomers()
		newCustomers := make([]Customer, 0, 0)
		for _, customer := range customers {
			if customer.Phone != phone {
				slices.Insert(newCustomers, 0, Customer{Name: customer.Name, Phone: customer.Phone})
			}
		}
		slices.Insert(newCustomers, 0, Customer{Name: name, Phone: phone})
		//TODO i am to dumb this array is still empty
		CreateJsonFile(newCustomers)
	})


	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

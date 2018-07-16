package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"encoding/csv"
	"log"
	"path/filepath"
	"io"
	"bufio"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

)

// our main function
func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/restapp")
	if err != nil {
		fmt.Print(err.Error())
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
            log.Fatal(err)
    }

	defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}
	type User struct {
		Id          string
		Firstname   string `json:"first_name"`
		Lastname    string `json:"last_name"`
		Email       string `json:"email"`
		Phonenumber string `json:"phone_number"`
	}
	router := gin.Default()
	// Add API handlers here

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello world!",
		})
	})

	router.POST("/users", func(c *gin.Context) {
		var user User
    	c.BindJSON(&user)
		stmt, err := db.Prepare("insert into users (first_name, last_name,email,phone_number) values(?,?,?,?);")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(user.Firstname, user.Lastname, user.Email, user.Phonenumber)

		if err != nil {
			fmt.Print(err.Error())
		}
		defer stmt.Close()
		c.JSON(http.StatusOK, gin.H{
			"message": "User successfully created",
		})
	})

	router.GET("/users", func(c *gin.Context) {
		var (
			user  User
			users []User
		)
		rows, err := db.Query("select id, first_name, last_name, email, phone_number from users;")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Email, &user.Phonenumber)
			users = append(users, user)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"result": users,
			"count":  len(users),
		})
	})

	router.GET("/users/:id", func(c *gin.Context) {
		var (
			user User
			result gin.H
		)
		id := c.Param("id")
		row := db.QueryRow("select id, first_name, last_name, email, phone_number from users where id = ?;", id)
		err = row.Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Email, &user.Phonenumber)
		if err != nil {
			// If no results send null
			result = gin.H{
				"result": nil,
				"count":  0,
			}
		} else {
			result = gin.H{
				"result": user,
				"count":  1,
			}
		}
		c.JSON(http.StatusOK, result)
	})

	router.PUT("/users/:id", func(c *gin.Context) {
		var user User
    	c.BindJSON(&user)
		Id := c.Param("id")
		stmt, err := db.Prepare("update users set first_name= ?, last_name= ?, email= ?, phone_number= ? where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(user.Firstname, user.Lastname, user.Email, user.Phonenumber, Id)
		if err != nil {
			fmt.Print(err.Error())
		}

		defer stmt.Close()
		c.JSON(http.StatusOK, gin.H{
			"message": "User Details successfully updated",
		})
	})

	router.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		stmt, err := db.Prepare("delete from users where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(id)
		if err != nil {
			fmt.Print(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully deleted user",
		})
	})

	router.GET("/export", func(c *gin.Context) {
		var (
			user  User
		)
		file, err := os.Create("users.csv")
		defer file.Close()
		if err != nil {
			fmt.Print(err.Error())
		}

		writer := csv.NewWriter(file)
		defer writer.Flush()

		rows, err := db.Query("select id, first_name, last_name, email, phone_number from users;")
		for rows.Next() {
			err = rows.Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Email, &user.Phonenumber)
			u := []string{ user.Firstname, user.Lastname, user.Email, user.Phonenumber }
			writer.Write(u)
		}

		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("File exported to %s", filepath.Join(dir, "/users.csv")),
		})
	})

	router.POST("/import", func(c *gin.Context) {
		csvFile, _ := os.Open(c.Query("path"))
		reader := csv.NewReader(bufio.NewReader(csvFile))
		for {
			line, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				log.Fatal(error)
			}
			stmt, err := db.Prepare("insert into users (first_name, last_name,email,phone_number) values(?,?,?,?);")
			if err != nil {
				fmt.Print(err.Error())
			}
			_, err = stmt.Exec(line[0], line[1], line[2], line[3])

			if err != nil {
				fmt.Print(err.Error())
			}
			defer stmt.Close()
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "File successfully imported",
		})
	})
	router.Run(":3000")
}


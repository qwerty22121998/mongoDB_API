package main

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

const DB_NAME string = "school"
const C_NAME string = "student"

type Student struct {
	ID   string `bson:"id"`
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

type Request struct {
	Status int
	Msg    string
	Data   interface{}
}

func getSession() (session *mgo.Session, err error) {
	session, err = mgo.Dial("localhost")
	if err != nil {
		fmt.Println("Can't connect to db")
		return nil, err
	}
	return session, nil
}
func getCollections(session *mgo.Session) *mgo.Collection {
	return session.DB(DB_NAME).C(C_NAME)
}

// CRUD - create / retrieve / update / delete
// Create
func insertStudent(collection *mgo.Collection, student Student) Request {

	find := retrieveStudent(collection, student.ID)
	if find.ID != "" {
		return Request{400, "Duplicated", nil}
	}
	err := collection.Insert(student)
	if err != nil {
		return Request{400, err.Error(), nil}
	}
	return Request{200, "Success", student}
}

// retrieve
func retrieveStudent(collection *mgo.Collection, id string) Student {
	var thisStudent Student
	err := collection.Find(bson.M{"id": id}).One(&thisStudent)

	if err != nil {
		return Student{}
	}
	return thisStudent
}

func retrieveAllStudent(collection *mgo.Collection) []Student {
	var students []Student
	err := collection.Find(bson.M{}).All(&students)
	if err != nil {
		return nil
	}
	return students
}

// update
func updateStudent(collection *mgo.Collection, student Student) Request {

	find := retrieveStudent(collection, student.ID)

	if find.ID == "" {
		return Request{404, "Not found", nil}
	}


	err := collection.Update(bson.M{"id": student.ID}, student)
	if err != nil {
		fmt.Println(err)
		return Request{400, err.Error(), nil}
	}
	return Request{200, "Success", student}
}

// delete
func deleteStudent(collection *mgo.Collection, id string) Request {
	currentStudent := retrieveStudent(collection, id)

	if currentStudent.ID == "" {
		return Request{404, "Not found", nil}
	}

	err := collection.Remove(bson.M{"id": id})
	if err != nil {
		return Request{404, err.Error(), nil}
	}
	return Request{200, "Success", currentStudent}
}
func deleteAll(collection *mgo.Collection) Request {
	change, err :=  collection.RemoveAll(bson.M{})
	if err != nil {
		return Request{400, err.Error(), nil}
	}
	return Request{200, "Success", change}
}


func serverStart(collection *mgo.Collection) {
	route := gin.Default()

	route.POST("/insert", func(context *gin.Context) {
		var currentStudent Student
		err := context.BindJSON(&currentStudent)
		if err != nil {
			context.JSON(400, Request{400, "Bad request", nil})
			return
		}
		res := insertStudent(collection, currentStudent)
		context.JSON(res.Status, res)
	})





	route.GET("/student", func(context *gin.Context) {
		students := retrieveAllStudent(collection)
		if len(students) == 0 {
			context.JSON(404, Request{404, "No data", nil})
			return
		}
		context.JSON(200, Request{200, "Success", students})
	})



	route.GET("/student/:id", func(context *gin.Context) {
		id := context.Param("id")
		currentStudent := retrieveStudent(collection, id)
		if currentStudent.ID == "" {
			context.JSON(404, Request{404, "Not found", nil})
			return
		}
		context.JSON(200, Request{200, "Success", currentStudent})

	})
	route.PATCH("/update", func(context *gin.Context) {
		var currentStudent Student
		err := context.BindJSON(&currentStudent)
		if err != nil {
			context.JSON(400, Request{400, err.Error(), nil})
			return
		}
		resp := updateStudent(collection, currentStudent)

		context.JSON(resp.Status, Request{resp.Status, resp.Msg, resp.Data})

	})
	route.DELETE("/delete/:id", func(context *gin.Context) {
		id := context.Param("id")
		resp := deleteStudent(collection, id)
		context.JSON(resp.Status, resp)

	})
	route.DELETE("/student", func(context *gin.Context) {
		resp := deleteAll(collection)
		context.JSON(resp.Status, resp)
	})


	route.Run(":8080")
}

func main() {
	s, _ := getSession()
	c := getCollections(s)
	serverStart(c)
}

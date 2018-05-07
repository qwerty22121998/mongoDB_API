package main

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"github.com/gin-gonic/gin/json"
)

const DB_NAME string = "school"
const C_NAME string = "student"

type Student struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
	Age  int    `bson:"age"`
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
func newStudent(id, name string, age int) Student {
	return Student{ID: id, Name: name, Age: age}
}

// CRUD - create / retrieve / update / delete
// Create
func createStudent(collection *mgo.Collection, student Student) error {
	err := collection.Insert(student)
	if err != nil {
		return fmt.Errorf("Cannot insert %v to DB", student)
	}
	return nil
}

// retrieve
func (s Student) ToJson() string {
	j, err := json.Marshal(s)
	if err != nil {
		return "Err"
	}
	return fmt.Sprintf("%v", string(j))
}
func retrieveStudent(collection *mgo.Collection, id string) string {
	var thisStudent Student
	err := collection.Find(bson.M{"_id": id}).One(&thisStudent)
	fmt.Println(thisStudent)
	if err != nil {
		return Student{}.ToJson()
	}
	return thisStudent.ToJson()
}
func retrieveAllStudent(collection *mgo.Collection) string {
	var students []Student
	err := collection.Find(bson.M{}).All(&students)
	if err != nil {
		return "Err"
	}
	j, err := json.Marshal(students)
	if err != nil {
		return "Err"
	}
	return fmt.Sprintf("%v", string(j))
}

// update
func updateStudent(collection *mgo.Collection, student Student) error {
	err := collection.Update(bson.M{"_id": student.ID}, student)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Can't update")
	}
	return nil
}

// delete
func deleteStudent(collection *mgo.Collection, id string) error {
	_, err := collection.RemoveAll(bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("Can't delete")
	}
	return nil
}
func serverStart(collection *mgo.Collection) {
	route := gin.Default()
	route.POST("/insert", func(context *gin.Context) {
		var stu Student

		err := context.BindJSON(&stu)
		if err != nil {
			context.String(http.StatusBadRequest, "400 Bad Request")
			return
		}
		err = createStudent(collection, stu)
		if err != nil {
			context.String(http.StatusBadRequest, err.Error())
			return;
		}
		context.String(http.StatusOK, "")
	})
	route.GET("/students", func(context *gin.Context) {
		context.String(http.StatusOK, retrieveAllStudent(collection))

	})
	route.GET("/student/:id", func(context *gin.Context) {
		context.String(http.StatusOK, retrieveStudent(collection, context.Param("id")))
	})
	route.PATCH("/update", func(context *gin.Context) {
		var student Student
		err := context.BindJSON(&student)
		fmt.Println(student)
		if err != nil {
			context.String(http.StatusBadRequest, "400 Bad Request")
			return
		}
		err = updateStudent(collection, student)
		if err != nil {
			context.String(http.StatusBadRequest, err.Error())
			return
		}
		context.String(http.StatusOK, "")
	})
	route.DELETE("/delete/:id", func(context *gin.Context) {
		err := deleteStudent(collection, context.Param("id"))
		if err != nil {
			context.String(http.StatusBadRequest, err.Error())
			return
		}
		context.String(http.StatusOK, "")
	})
	route.Run(":8080")
}

func main() {
	s, _ := getSession()
	c := getCollections(s)
	serverStart(c)
}

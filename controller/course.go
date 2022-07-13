package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pilinux/gorest/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"reflect"
	"time"

	"github.com/pilinux/gorest/lib/renderer"
)

const (
	DbName         = "sillage"
	CollectionName = "course"
)

type Situation struct {
	Teacher string   `json:"teacher" bson:"teacher"`
	Room    string   `json:"room" bson:"room"`
	Groups  []string `json:"groups" bson:"groups"`
}

type CourseInfo struct {
	Name string `json:"name" bson:"name"`
	Code string `json:"code" bson:"code"`
	Bgc  string `json:"bgc" bson:"bgc"`
}

type Course struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Grade      string             `json:"grade,omitempty" bson:"grade,omitempty"`
	Dates      []string           `json:"dates,omitempty" bson:"dates"`
	LessonNum  int                `json:"lessonNum" bson:"lessonNum"`
	Info       CourseInfo         `json:"info" bson:"info"`
	Situations []Situation        `json:"situations" bson:"situations"`
	Method     string             `json:"method" bson:"method"`
	Note       string             `json:"note" bson:"note"`
}

func (c Course) isEmpty() bool {
	return reflect.DeepEqual(c, Course{})
}

func CourseCreate(c *gin.Context) {
	data := Course{}
	if err := c.ShouldBindJSON(&data); err != nil {
		renderer.Render(c, gin.H{"msg": "bad request"}, http.StatusBadRequest)
		return
	}

	// return bad request if it's empty
	if data.isEmpty() {
		renderer.Render(c, gin.H{"msg": "empty body"}, http.StatusBadRequest)
		return
	}

	// generate a new ObjectID
	data.ID = primitive.NewObjectID()

	client := database.GetMongo()                                            // connect MongoDB
	db := client.Database(DbName)                                            // set database name
	collection := db.Collection(CollectionName)                              // set collection name
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // set max TTL
	defer cancel()

	// insert one document
	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		renderer.Render(c, gin.H{"msg": "internal server error"}, http.StatusInternalServerError)
		return
	}

	renderer.Render(c, data, http.StatusCreated)
}

func CourseRetrieveList(c *gin.Context) {
	client := database.GetMongo()                                            // connect MongoDB
	db := client.Database(DbName)                                            // set database name
	collection := db.Collection(CollectionName)                              // set collection name
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // set max TTL
	defer cancel()

	var data []Course
	err := collection.Find(ctx, bson.M{}).All(&data)
	if err != nil {
		renderer.Render(c, gin.H{"msg": "internal server error"}, http.StatusInternalServerError)
		return
	}
	if len(data) == 0 {
		renderer.Render(c, gin.H{"msg": "no record found"}, http.StatusNotFound)
		return
	}

	renderer.Render(c, data, http.StatusOK)
}

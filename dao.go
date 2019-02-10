package dao

import (
  "fmt"
  "github.com/mongodb/mongo-go-driver/mongo"
)

type DAO struct {
  DB *mongo.Database
  Collection *mongo.Collection
}

func New(db *mongo.Database, collection *mongo.Collection) (DAO) {
  return DAO{DB: db, Collection: collection}
}

func init() {
  fmt.Println("Hello world!")
}

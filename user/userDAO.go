package user

import (
  //NATIVE
  "context"
  "fmt"

  //THIRD PARTY
  "github.com/mongodb/mongo-go-driver/bson"
  "github.com/mongodb/mongo-go-driver/mongo"
  "github.com/mongodb/mongo-go-driver/mongo/options"
)

type UserRecord struct {
  ID        string `bson:"_id" json:"id"`
  FirstName string `bson:"firstName" json:"firstName,omitempty"`
  LastName  string `bson:"lastName" json:"lastName,omitempty"`
  Birthday  string `bson:"birthday" json:"birthday,omitempty"`
  ImageURL  string `bson:"imageURL" json:"imageURL,omitempty"`
  Friends   []string `bson:"friends" json:"friends,omitempty"`
}

type Friend struct {
  ID        string `bson:"_id" json:"id"`
  FirstName string `bson:"firstName" json:"firstName,omitempty"`
  LastName  string `bson:"lastName" json:"lastName,omitempty"`
  Birthday  string `bson:"birthday" json:"birthday,omitempty"`
  ImageURL  string `bson:"imageURL" json:"imageURL,omitempty"`
}

type DAO struct {
  DB *mongo.Database
  Collection *mongo.Collection
}

func (dao *DAO) FindById(id string) (UserRecord, error) {
	var u UserRecord
  IDFilter := bson.M{"_id": id}
  fmt.Println(IDFilter)
	err := dao.Collection.FindOne(context.Background(), IDFilter).Decode(&u)
	return u, err
}

func (dao *DAO) Upsert(user UserRecord) (error) {
  fmt.Println(user)
  IDFilter := bson.M{"_id": user.ID}
  update := bson.D{{"$set", user}}//mongo.NewUpdateOneModel().SetUpdate(user) ID filter works properly because it's same as FindById. The update model is incorrect.
  opts := options.Update().SetUpsert(true)
  fmt.Println(update)
  _, err := dao.Collection.UpdateOne(context.Background(), IDFilter, update, opts)
  fmt.Println(err)
  return err
}

func (dao *DAO) Delete(id string) (UserRecord, error) {
  var u UserRecord
  IDFilter := bson.M{"_id": id}
  err := dao.Collection.FindOneAndDelete(context.Background(), IDFilter).Decode(&u)
  return u, err
}

func (dao *DAO) Append (filterId string, bodyId string, field string) (error) {
  IDFilter := bson.M{"_id": filterId}
  update := bson.D{{"$addToSet", bson.M{(field): bodyId}}}
  _, err := dao.Collection.UpdateOne(context.Background(), IDFilter, update)
  return err
}

func (dao *DAO) Remove(filterId string, bodyId string, field string) (error) {
  IDFilter := bson.M{"_id": filterId}
  update := bson.D{{"$pull", bson.M{(field): bodyId}}}
  _, err := dao.Collection.UpdateOne(context.Background(), IDFilter, update)
  return err
}

func (dao *DAO) FindWhereArrayContains(filterItem string, filterField string) ([]Friend, error) {
  ctx := context.Background()
  var friendList []Friend
  filter := bson.M{(filterField): bson.D{{"$all", bson.A{filterItem}}}}
  cursor, err := dao.Collection.Find(ctx, filter)
  if err != nil {
    return friendList, err
  }
  defer cursor.Close(ctx)
  for cursor.Next(ctx) {
    var f Friend
    err = cursor.Decode(&f)
    if err != nil {
      return friendList, err
    }
    friendList = append(friendList, f)
  }
  return friendList, err
}

func (dao *DAO) SymmetricRemove(firstId string, secondId string, fields []string) (error) {
  deleteField := fields[0]
  err := dao.Remove(firstId, secondId, deleteField)
  if err != nil {
    return err
  }
  if (len(fields) > 1) {
    deleteField = fields[1]
  }
  err = dao.Remove(secondId, firstId, deleteField)
  return err
}

func (dao *DAO) SymmetricAppend(firstId string, secondId string, fields []string) (error) {
  appendField := fields[0]
  err := dao.Append(firstId, secondId, appendField)
  if err != nil {
    return err
  }
  if (len(fields) > 1) {
    appendField = fields[1]
  }
  err = dao.Append(secondId, firstId, appendField)
  return err
}

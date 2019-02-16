package user

import (
	//NATIVE
	"context"
	"dao/place"
	"fmt"
	//  "dao/event"

	//THIRD PARTY
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

type UserRecord struct {
	ID        string        `bson:"_id" json:"id"`
	FirstName string        `bson:"firstName,omitempty" json:"firstName,omitempty"`
	LastName  string        `bson:"lastName,omitempty" json:"lastName,omitempty"`
	Birthday  string        `bson:"birthday,omitempty" json:"birthday,omitempty"`
	ImageURL  string        `bson:"imageURL,omitempty" json:"imageURL,omitempty"`
	Friends   []string      `bson:"friends,omitempty" json:"friends,omitempty"`
	Places    []string 			`bson:"places,omitempty" json:"places,omitempty"`
}

type UserFriends struct {
	Friends   []Friend      `bson:"friends" json:"friends,omitempty"`
}

type UserPlaces struct {
	Places    []place.Place `bson:"places" json:"places,omitempty"`
}

type Friend struct {
	ID        string `bson:"_id" json:"id"`
	FirstName string `bson:"firstName,omitempty" json:"firstName,omitempty"`
	LastName  string `bson:"lastName,omitempty" json:"lastName,omitempty"`
	Birthday  string `bson:"birthday,omitempty" json:"birthday,omitempty"`
	ImageURL  string `bson:"imageURL,omitempty" json:"imageURL,omitempty"`
}

type DAO struct {
	DB         *mongo.Database
	Collection *mongo.Collection
}

func New(db *mongo.Database, collection *mongo.Collection) DAO {
	return DAO{DB: db, Collection: collection}
}

func (dao *DAO) FindById(id string) (UserRecord, error) {
	var u UserRecord
	IDFilter := bson.M{"_id": id}
	fmt.Println(IDFilter)
	err := dao.Collection.FindOne(context.Background(), IDFilter).Decode(&u)
	return u, err
}

func (dao *DAO) Upsert(user UserRecord) error {
	fmt.Println(user)
	IDFilter := bson.M{"_id": user.ID}
	update := bson.D{{"$set", user}} //mongo.NewUpdateOneModel().SetUpdate(user) ID filter works properly because it's same as FindById. The update model is incorrect.
	opts := options.Update().SetUpsert(true)
	_, err := dao.Collection.UpdateOne(context.Background(), IDFilter, update, opts)
	return err
}

func (dao *DAO) Delete(id string) (UserRecord, error) {
	var u UserRecord
	IDFilter := bson.M{"_id": id}
	err := dao.Collection.FindOneAndDelete(context.Background(), IDFilter).Decode(&u)
	return u, err
}

func (dao *DAO) Append(filterId string, bodyId string, field string) error {
	IDFilter := bson.M{"_id": filterId}
	update := bson.D{{"$addToSet", bson.M{(field): bodyId}}}
	_, err := dao.Collection.UpdateOne(context.Background(), IDFilter, update)
	return err
}

func (dao *DAO) Remove(filterId string, bodyId string, field string) error {
	IDFilter := bson.M{"_id": filterId}
	update := bson.D{{"$pull", bson.M{(field): bodyId}}}
	_, err := dao.Collection.UpdateOne(context.Background(), IDFilter, update)
	return err
}

func (dao *DAO) FriendLookup(filterId string, filterField string) (UserFriends, error) {
	ctx := context.Background()
	var itemList []UserFriends
	var uf UserFriends
	pipeline := bson.A{
		bson.M{"$match": bson.M{"_id": filterId}},
		bson.M{"$lookup": bson.M{
			"from":         "User",
			"localField":   filterField,
			"foreignField": "_id",
			"as":           filterField,
		},
		},
	}
	cursor, err := dao.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return uf, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var u UserFriends
		err = cursor.Decode(&u)
		if err != nil {
			return u, err
		}
		itemList = append(itemList, u)
	}
	if (len(itemList) > 0) {
		friends := itemList[0]
		return friends, err
	}
	return uf, err
}

func (dao *DAO) PlaceLookup(filterId string, filterField string) (UserPlaces, error) {
	ctx := context.Background()
	var up UserPlaces
	var itemList []UserPlaces
	pipeline := bson.A{
		bson.M{"$match": bson.M{"_id": filterId}},
		bson.M{"$lookup": bson.M{
			"from":         "Place",
			"localField":   filterField,
			"foreignField": "_id",
			"as":           "places",
		},
		},
	}
	cursor, err := dao.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return up, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var u UserPlaces
		err = cursor.Decode(&u)
		if err != nil {
			return u, err
		}
		itemList = append(itemList, u)
	}
	if (len(itemList) > 0) {
		places := itemList[0]
		return places, err
	}
	return up, err
}

func (dao *DAO) SymmetricRemove(firstId string, secondId string, fields []string) error {
	deleteField := fields[0]
	err := dao.Remove(firstId, secondId, deleteField)
	if err != nil {
		return err
	}
	if len(fields) > 1 {
		deleteField = fields[1]
	}
	err = dao.Remove(secondId, firstId, deleteField)
	return err
}

func (dao *DAO) SymmetricAppend(firstId string, secondId string, fields []string) error {
	appendField := fields[0]
	err := dao.Append(firstId, secondId, appendField)
	if err != nil {
		return err
	}
	if len(fields) > 1 {
		appendField = fields[1]
	}
	err = dao.Append(secondId, firstId, appendField)
	return err
}

func PrefixField(field string) string {
	return fmt.Sprintf("$%v", field)
}

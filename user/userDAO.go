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
	FirstName string        `bson:"firstName" json:"firstName,omitempty"`
	LastName  string        `bson:"lastName" json:"lastName,omitempty"`
	Birthday  string        `bson:"birthday" json:"birthday,omitempty"`
	ImageURL  string        `bson:"imageURL" json:"imageURL,omitempty"`
	Friends   []string      `bson:"friends" json:"friends,omitempty"`
	Places    []place.Place `bson:"places" json:"places,omitempty"`
}

type Friend struct {
	ID        string `bson:"_id" json:"id"`
	FirstName string `bson:"firstName" json:"firstName,omitempty"`
	LastName  string `bson:"lastName" json:"lastName,omitempty"`
	Birthday  string `bson:"birthday" json:"birthday,omitempty"`
	ImageURL  string `bson:"imageURL" json:"imageURL,omitempty"`
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

func (dao *DAO) FriendLookup(filterItem string, filterField string, target map[string]string) ([]Friend, error) {
	ctx := context.Background()
	var itemList []Friend
	targetCollection, targetField := target["collection"], target["field"]
	pipeline := bson.A{
		bson.M{"$match": bson.M{"_id": filterItem}},
		bson.M{"$unwind": filterField},
		bson.M{"$lookup": bson.M{
			"from":         targetCollection,
			"localField":   filterField,
			"foreignField": targetField,
			"as":           filterField,
		},
		},
	}
	cursor, err := dao.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return itemList, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var t Friend
		err = cursor.Decode(&t)
		if err != nil {
			return itemList, err
		}
		itemList = append(itemList, t)
	}
	return itemList, err
}

func (dao *DAO) PlaceLookup(filterItem string, filterField string, target map[string]string) ([]place.Place, error) {
	ctx := context.Background()
	var itemList []place.Place
	targetCollection, targetField := target["collection"], target["field"]
	pipeline := bson.A{
		bson.M{"$match": bson.M{"_id": filterItem}},
		bson.M{"$unwind": filterField},
		bson.M{"$lookup": bson.M{
			"from":         targetCollection,
			"localField":   filterField,
			"foreignField": targetField,
			"as":           filterField,
		},
		},
	}
	cursor, err := dao.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return itemList, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var p place.Place
		err = cursor.Decode(&p)
		if err != nil {
			return itemList, err
		}
		itemList = append(itemList, p)
	}
	return itemList, err
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

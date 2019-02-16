package event

import (
	//NATIVE
	"context"
	"dao/place"
	"fmt"
	"github.com/oleiade/reflections"
	"reflect"

	//THIRD PARTY
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

type TimeOption struct {
	ID    string `bson:"timeId,omitempty" json:"timeId,omitempty"`
	StartTime string `bson:"startTime,omitempty" json:"startTime,omitempty"`
	EndTime   string `bson:"endTime,omitempty" json:"endTime,omitempty"`
}

type Member struct {
	ID        string `bson:"id,omitempty" json:"id,omitempty"`
	FirstName string `bson:"firstName,omitempty" json:"firstName,omitempty"`
	LastName  string `bson:"lastName,omitempty" json:"lastName,omitempty"`
	ImageURL  string `bson:"imageURL,omitempty" json:"imageURL,omitempty"`
}

type Votes struct {
	Times  map[string]int `bson:"times,omitempty" json:"times,omitempty"`
	Places map[string]int `bson:"places,omitempty" json:"places,omitempty"`
}

type Event struct {
	ID      primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Title   string             `bson:"title,omitempty" json:"title"`
	Times   []TimeOption       `bson:"times,omitempty" json:"times"`
	Members []Member           `bson:"members,omitempty" json:"members"`
	Places  []place.Place      `bson:"places,omitempty" json:"places"`
	Votes   Votes              `bson:"votes,omitempty" json:"votes,omitempty"`
}

type UserVote struct {
	ID string `json:"id,omitempty"`
}

type DAO struct {
	DB         *mongo.Database
	Collection *mongo.Collection
}

func New(db *mongo.Database, collection *mongo.Collection) DAO {
	return DAO{DB: db, Collection: collection}
}

func (dao *DAO) FindById(id string) (Event, error) { //DONE
	var e Event
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println(err)
		return e, err
	}
	e.ID = objectId
	IDFilter := bson.M{"_id": e.ID}
	err = dao.Collection.FindOne(context.Background(), IDFilter).Decode(&e)
	fmt.Println(err)
	return e, err
}

func (dao *DAO) Upsert(event Event) (Event, error) { //DONE U+C
	var e Event
	if event.ID == primitive.NilObjectID {
		event.ID = primitive.NewObjectID()
	}
	IDFilter := bson.M{"_id": event.ID}
	fmt.Println(event)
	updateObj, err := BuildUpdate(event)
	fmt.Println(updateObj)
	if err != nil {
		return e, err
	}
	update := bson.D{{"$set", updateObj}} //mongo.NewUpdateOneModel().SetUpdate(user) ID filter works properly because it's same as FindById. The update model is incorrect.
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	err = dao.Collection.FindOneAndUpdate(context.Background(), IDFilter, update, opts).Decode(&e)
	return e, err
}

func (dao *DAO) Delete(id string) (Event, error) {
	var e Event
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println(err)
		return e, err
	}
	e.ID = objectId
	IDFilter := bson.M{"_id": e.ID}
	err = dao.Collection.FindOneAndDelete(context.Background(), IDFilter).Decode(&e)
	return e, err
}

func (dao *DAO) Append(filterId string, item interface{}, field string) error {
	objectId, err := primitive.ObjectIDFromHex(filterId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	IDFilter := bson.M{"_id": objectId}
	update := bson.D{{"$push", bson.M{(field): item}}}
	_, err = dao.Collection.UpdateOne(context.Background(), IDFilter, update)
	return err
}

func (dao *DAO) Remove(filterId string, bodyId string, field string) error {
	eventId, err := primitive.ObjectIDFromHex(filterId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	IDFilter := bson.M{"_id": eventId}
	update := bson.M{
		"$pull": bson.M{
			(field): bson.M{
				"id": bson.M{
					"$eq": bodyId,
				},
			},
		},
	}
	_, err = dao.Collection.UpdateOne(context.Background(), IDFilter, update)
	return err
}

func (dao *DAO) IncrementField(collectionFilterId string, docFilterId string, field string, nestedField string) error {
	eventId, err := primitive.ObjectIDFromHex(collectionFilterId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	compoundField := fmt.Sprintf("%v.%v", field, docFilterId)
	IDFilter := bson.M{"_id": eventId}
	update := bson.M{"$inc": bson.M{
		(compoundField): 1,
	},
	}
	_, err = dao.Collection.UpdateOne(context.Background(), IDFilter, update)
	return err
}

func BuildUpdate(event Event) (Event, error) {
	var e Event
	fieldList, err := reflections.Fields(event)
	if err != nil {
		return e, err
	}
	for _,v := range fieldList {
		fieldVal, err := reflections.GetField(event, v)
		if err != nil {
			return e, err
		}
		if FieldNotNil(fieldVal) {
			fmt.Println(fieldVal)
			fmt.Println(reflect.TypeOf(fieldVal))
			err = reflections.SetField(&e, v, fieldVal)
			if err != nil {
				return e, err
			}
		}
	}
	return e, err
}

func FieldNotNil (fieldVal interface{}) (bool) {
	var emptyMembers []Member
	var emptyTimes []TimeOption
	var emptyPlaces []place.Place
	if reflect.DeepEqual(fieldVal, emptyTimes) {
		return false
	}
	if reflect.DeepEqual(fieldVal, emptyPlaces) {
		return false
	}
	if reflect.DeepEqual(fieldVal, emptyMembers) {
		return false
	}
	return true
}

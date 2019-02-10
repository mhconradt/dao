package event

import (
	//NATIVE
	"context"
	"dao/place"
	"fmt"

	//THIRD PARTY
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

type TimeOption struct {
	TimeID    string `bson:"timeId" json:"timeId,omitempty"`
	StartTime string `bson:"startTime" json:"startTime,omitempty"`
	EndTime   string `bson:"endTime" json:"endTime,omitempty"`
}

type Member struct {
	ID        string `bson:"id" json:"id,omitempty"`
	FirstName string `bson:"firstName" json:"firstName,omitempty"`
	LastName  string `bson:"lastName" json:"lastName,omitempty"`
	ImageURL  string `bson:"imageURL" json:"imageURL,omitempty"`
}

type Votes struct {
	Times  map[string]int
	Places map[string]int
}

type Event struct {
	ID      primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Title   string             `bson:"title" json:"title"`
	Times   []TimeOption       `bson:"times" json:"times"`
	Members []Member           `bson:"members" json:"members"`
	Places  []place.Place      `bson:"places" json:"places"`
	Votes   Votes              `bson:"votes" json:"votes,omitempty"`
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
	if event.ID == primitive.NilObjectID {
		event.ID = primitive.NewObjectID()
	}
	IDFilter := bson.M{"_id": event.ID}
	update := bson.D{{"$set", event}} //mongo.NewUpdateOneModel().SetUpdate(user) ID filter works properly because it's same as FindById. The update model is incorrect.
	opts := options.Update().SetUpsert(true)
	result, err := dao.Collection.UpdateOne(context.Background(), IDFilter, update, opts)
	if result.UpsertedID != nil {
		UpsertedID := result.UpsertedID.(primitive.ObjectID)
		event.ID = UpsertedID
	}
	return event, err
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
	update := bson.D{{"$addToSet", bson.M{(field): item}}}
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

func (dao *DAO) IncrementField(collectionFilterId string, docFilterId string, field string) error {
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

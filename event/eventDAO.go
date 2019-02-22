package event

import (
	//NATIVE
	"context"
	"dao/place"

	//THIRD PARTY
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

type TimeOption struct {
	ID    string `bson:"id,omitempty" json:"id,omitempty"`
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
	Times  *map[string][]string `bson:"times,omitempty" json:"times,omitempty"`
	Places *map[string][]string `bson:"places,omitempty" json:"places,omitempty"`
}

type Event struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	Title   string             `bson:"title,omitempty" json:"title,omitempty"`
	Times   []TimeOption       `bson:"times,omitempty" json:"times,omitempty"`
	Members []Member           `bson:"members,omitempty" json:"members,omitempty"`
	Places  []place.Place      `bson:"places,omitempty" json:"places,omitempty"`
	PlaceVotes   *map[string][]string              `bson:"placesVotes,omitempty" json:"placesVotes,omitempty"`
	TimeVotes   *map[string][]string              `bson:"timesVotes,omitempty" json:"timesVotes,omitempty"`
}

type DAO struct {
	DB         *mongo.Database
	Collection *mongo.Collection
}

func New(db *mongo.Database, collection *mongo.Collection) DAO {
	return DAO{DB: db, Collection: collection}
}

func (dao *DAO) FindById(filterId string) (Event, error) { //DONE
	var e Event
	IDFilter, err := GetIDFilter(filterId)
	if err != nil {
		return e, err
	}
	err = dao.Collection.FindOne(context.Background(), IDFilter).Decode(&e)
	return e, err
}

func (dao *DAO) Upsert(event Event) (Event, error) { //DONE U+C
	var e Event
	if event.ID == primitive.NilObjectID {
		event.ID = primitive.NewObjectID()
	}
	IDFilter := bson.M{"_id": event.ID}

	update := bson.M{"$set": event} //mongo.NewUpdateOneModel().SetUpdate(user) ID filter works properly because it's same as FindById. The update model is incorrect.
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	err := dao.Collection.FindOneAndUpdate(context.Background(), IDFilter, update, opts).Decode(&e)
	return e, err
}

func (dao *DAO) Delete(filterId string) (Event, error) {
	var e Event
	IDFilter, err := GetIDFilter(filterId)
	if err != nil {
		return e, err
	}
	err = dao.Collection.FindOneAndDelete(context.Background(), IDFilter).Decode(&e)
	return e, err
}

func (dao *DAO) Append(filterId string, item interface{}, field string) (Event, error) {
	update := bson.M{"$addToSet": bson.M{(field): item}}
	event, err := dao.ExecuteUpdate(filterId, update)
	return event, err
}

func (dao *DAO) Remove(filterId string, bodyId string, field string) (Event, error) {
	update := bson.M{
		"$pull": bson.M{
			(field): bson.M{
				"id": bson.M{
					"$eq": bodyId,
				},
			},
		},
	}
	event, err := dao.ExecuteUpdate(filterId, update)
	return event, err
}

func (dao *DAO) IncrementField(filterId string, field string) (Event, error) {
	update := bson.M{"$inc": bson.M{
		(field): 1,
		},
	}
	event, err := dao.ExecuteUpdate(filterId, update)
	return event, err
}

func GetIDFilter(filterId string) (bson.M, error) {
	objectId, err := primitive.ObjectIDFromHex(filterId)
	if err != nil {
		return bson.M{}, err
	}
	return bson.M{"_id": objectId}, err
}

func (dao *DAO) ExecuteUpdate(id string, update bson.M) (Event, error) {
	var e Event
	IDFilter, err := GetIDFilter(id)
	if err != nil {
		return e, err
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	err = dao.Collection.FindOneAndUpdate(context.Background(), IDFilter, update, opts).Decode(&e)
	return e, err
}

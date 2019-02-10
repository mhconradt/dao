package place

import (
	//NATIVE
	"context"

	//THIRD PARTY
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"googlemaps.github.io/maps"
)

type Location struct {
	GeoPoint maps.LatLng `bson:"coordinates" json:"coordinates,omitempty"`
	Address  Address     `bson:"address" json:"address,omitempty"`
}

type Address struct {
	Address1 string `bson:"address1" json:"address1"`
	City     string `bson:"city" json:"city"`
	State    string `bson:"state" json:"state"`
	ZipCode  string `bson:"zipCode" json:"zipCode"`
}

type Place struct {
	ID       string   `bson:"_id" json:"id"`
	Name     string   `bson:"name" json:"name,omitempty"`
	Location Location `bson:"location" json:"location,omitempty"`
	ImageURL string   `bson:"imageURL" json:"imageURL,omitempty"`
	Rating   float64  `bson:"rating" json:"rating,omitempty"`
}

type DAO struct {
	DB         *mongo.Database
	Collection *mongo.Collection
}

func New(db *mongo.Database, collection *mongo.Collection) DAO {
	return DAO{DB: db, Collection: collection}
}

func (dao *DAO) FindById(id string) (Place, error) { //DONE
	var p Place
	IDFilter := bson.M{"_id": id}
	err := dao.Collection.FindOne(context.Background(), IDFilter).Decode(&p)
	return p, err
}

func (dao *DAO) Upsert(place Place) (Place, error) {
	IDFilter := bson.M{"_id": place.ID}
	update := bson.D{{"$set", place}} //mongo.NewUpdateOneModel().SetUpdate(user) ID filter works properly because it's same as FindById. The update model is incorrect.
	opts := options.Update().SetUpsert(true)
	_, err := dao.Collection.UpdateOne(context.Background(), IDFilter, update, opts)
	return place, err
}

func (dao *DAO) Delete(id string) (Place, error) {
	var p Place
	IDFilter := bson.M{"_id": id}
	err := dao.Collection.FindOneAndDelete(context.Background(), IDFilter).Decode(&p)
	return p, err
}

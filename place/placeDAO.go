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
	GeoPoint *maps.LatLng `bson:"coordinates,omitempty" json:"coordinates,omitempty"`
	Address  *Address     `bson:"address,omitempty" json:"address,omitempty"`
}

type Address struct {
	Address1 string `bson:"address1,omitempty" json:"address1,omitempty"`
	City     string `bson:"city,omitempty" json:"city,omitempty"`
	State    string `bson:"state,omitempty" json:"state,omitempty"`
	ZipCode  string `bson:"zipCode,omitempty" json:"zipCode,omitempty"`
}

type Place struct {
	ID       string   `bson:"_id" json:"id"`
	Name     *string   `bson:"name,omitempty" json:"name,omitempty"`
	Location *Location `bson:"location,omitempty" json:"location,omitempty"`
	ImageURL *string   `bson:"imageURL,omitempty" json:"imageURL,omitempty"`
	Rating   *float64  `bson:"rating,omitempty" json:"rating,omitempty"`
	Categories *[]string `bson:"categories,omitempty" json:"categories,omitempty"`
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
	var p Place
	IDFilter := bson.M{"_id": place.ID}
	update := bson.D{{"$set", place}} //mongo.NewUpdateOneModel().SetUpdate(user) ID filter works properly because it's same as FindById. The update model is incorrect.
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	err := dao.Collection.FindOneAndUpdate(context.Background(), IDFilter, update, opts).Decode(&p)
	if err != nil {
		return p, err
	}
	return p, err
}

func (dao *DAO) Delete(id string) (Place, error) {
	var p Place
	IDFilter := bson.M{"_id": id}
	err := dao.Collection.FindOneAndDelete(context.Background(), IDFilter).Decode(&p)
	return p, err
}

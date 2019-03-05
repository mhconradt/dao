package place

import (
	//NATIVE
	"context"
	"fmt"

	//THIRD PARTY
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GeoJSON struct {
	Type        string    `bson:"type,omitempty" json:"type,omitempty"`
	Coordinates []float64 `bson:"coordinates,omitempty" json:"coordinates,omitempty"`
}

type Location struct {
	Address  Address `bson:"address,omitempty" json:"address,omitempty"`
	GeoPoint GeoJSON `bson:"geopoint,omitempty" json:"geopoint,omitempty"`
}

// Closed or all day

type TimeRange struct {
	Start int `bson:"start,omitempty" json:"start,omitempty"`
	End   int `bson:"end,omitempty" json:"end,omitempty"`
}

type Hours struct {
	Sunday    []TimeRange `bson:"sunday,omitempty" json:"sunday,omitempty"`
	Monday    []TimeRange `bson:"monday,omitempty" json:"monday,omitempty"`
	Tuesday   []TimeRange `bson:"tuesday,omitempty" json:"tuesday,omitempty"`
	Wednesday []TimeRange `bson:"wednesday,omitempty" json:"wednesday,omitempty"`
	Thursday  []TimeRange `bson:"thursday,omitempty" json:"thursday,omitempty"`
	Friday    []TimeRange `bson:"friday,omitempty" json:"friday,omitempty"`
	Saturday  []TimeRange `bson:"saturday,omitempty" json:"saturday,omitempty"`
}

type Address struct {
	Address1 string `bson:"address1,omitempty" json:"address1,omitempty"`
	City     string `bson:"city,omitempty" json:"city,omitempty"`
	State    string `bson:"state,omitempty" json:"state,omitempty"`
	ZipCode  int    `bson:"zipCode,omitempty" json:"zipCode,omitempty"`
}

type Place struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Name       string             `bson:"name,omitempty" json:"name,omitempty"`
	Location   Location           `bson:"location,omitempty" json:"location,omitempty"`
	ImageURL   string             `bson:"imageURL,omitempty" json:"imageURL,omitempty"`
	URL        string             `bson:"url,omitempty" json:"url,omitempty"`
	Rating     float64            `bson:"rating,omitempty" json:"rating,omitempty"`
	Price      float32            `bson:"price,omitempty" json:"price,omitempty"`
	Categories []string           `bson:"categories,omitempty" json:"categories,omitempty"`
	Hours      Hours              `bson:"hours,omitempty" json:"hours,omitempty"`
	Type       string             `bson:"type,omitempty" json:"type,omitempty"`
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

func (dao *DAO) FindNear(lat, lng float64, radius int, filters []bson.M) (Place, error) { //DONE
	var p Place
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

func (dao *DAO) FilterEmpty() (*mongo.DeleteResult, error) {
	bigFilter := bson.M{"type": ""}
	return dao.DB.Collection("PlaceBeta").DeleteMany(context.Background(), bigFilter)
}

func (dao *DAO) BulkWrite(p []Place, collection string) (*mongo.BulkWriteResult, error) {
	numPlaces := len(p)
	inputChannel := make(chan Place, numPlaces)
	outputChannel := make(chan *mongo.InsertOneModel)
	signalChannel := make(chan bool)
	for i := 0; i < 200; i++ {
		go MakeModel(inputChannel, outputChannel, signalChannel)
	}
	for _, place := range p {
		inputChannel <- place
	}
	close(inputChannel)
	var modelList []mongo.WriteModel
	runLoop := true
	for runLoop {
		model := <-outputChannel
		modelList = append(modelList, model)
		if len(modelList) == numPlaces {
			runLoop = false
		}
	}
	opts := options.BulkWrite().SetOrdered(false)
	writeResult, err := dao.DB.Collection(collection).BulkWrite(context.Background(), modelList, opts)
	if err != nil {
		fmt.Println(err)
		return writeResult, err
	}
	return writeResult, err
}

func (dao *DAO) GetAll() ([]Place, error) {
	var p []Place
	collection := dao.DB.Collection("PlaceBeta")
	bigFilter := bson.M{"type": "Restaurant"}
	ctx := context.Background()
	cur, err := collection.Find(ctx, bigFilter)
	if err != nil {
		return p, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var plc Place
		cur.Decode(&plc)
		p = append(p, plc)
	}
	return p, nil
}

func (dao *DAO) DeleteAll() (*mongo.DeleteResult, error) {
	bigFilter := bson.M{"type": "Restaurant"}
	return dao.DB.Collection("PlaceBeta").DeleteMany(context.Background(), bigFilter)
}

func (dao *DAO) Delete(id string) (Place, error) {
	var p Place
	IDFilter := bson.M{"_id": id}
	err := dao.Collection.FindOneAndDelete(context.Background(), IDFilter).Decode(&p)
	return p, err
}

func MakeModel(inputChannel <-chan Place, outputChannel chan<- *mongo.InsertOneModel, done chan<- bool) {
	for place := range inputChannel {
		doc := bson.M{
			"name":       place.Name,
			"location":   place.Location,
			"imageURL":   place.ImageURL,
			"url":        place.URL,
			"rating":     place.Rating,
			"price":      place.Price,
			"categories": place.Categories,
			"hours":      place.Hours,
			"type":       "Restaurant",
		}
		newModel := mongo.NewInsertOneModel().SetDocument(doc)
		outputChannel <- newModel
	}
	return
}

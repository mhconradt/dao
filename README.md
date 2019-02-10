package dao

Language: Golang

Can be used by importing "github.com/mhconradt/dao"
Developers can alter this code, but importing it should be the single source to prevent divergence in the code base.

Current Version: v0.1

Version 0.1:

Contains event and user data access objects.

Common Functions:
  FindById: Finds document by ID
  Upsert: Updates document and adds if the ID does not exist
  Delete: Deletes a document that matches the ID condition
  Append: Adds an item to an array
  Remove: Removes an item from an array

Event:
Contains types DAO, Event, TimeOption, Place, Address and Member. All are structs.

DAO: contains DB client and collection to be used.

Event: contains list of Places, TimeOptions and Members, as well as a title.
  Member: a limited UserRecord.
  TimeOption: contains StartTime, EndTime and Votes.
  (Considering adding ID)
  Place: Contains placeId, the ID of a document in the place collection, votes, GeoPoint, Address, Name and ImageURL.
    Address: Contains address1, city, state and zipCode.
  I am considering changing votes from an integer to list of userIds

Unique Functions:
  IncrementField: Increments a field in an array that matches a condition in a document that matches a condition.

User:
Contains types DAO, UserRecord and Friend.
  UserRecord:
    Stores a user's friends list, an object with sent and received friend requests, name, image and birthday.
  Friend: Stores same as UserRecord but excludes friends and requests.

Unique Functions:
  SymmetricAppend: Given two IDs and a list of fields, it will add the ID of each to the field in the other's document. If multiple fields, the first user is added to the second field.
  SymmetricRemove: Same as Append, but Removes.
  FindWhereArrayContains: Finds documents that contain a specified item in a specified field.

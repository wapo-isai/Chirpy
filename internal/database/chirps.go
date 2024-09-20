package database

import (
	"fmt"
	"strconv"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
	AuthorId int `json:"author_id"` 
}

func (db *DB) CreateChirp(body string, authorId string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	userId, err := strconv.Atoi(authorId)

	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:   id,
		Body: body,
		AuthorId: userId,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, ErrNotExist
	}

	return chirp, nil
}

// 1. ensure user is authenticated by validating token
// 2. get the user information from the token
// 3. ensure token user id is the chirp authorid

func (db *DB) DeleteChirp(id int, userId string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]

	if !ok {
		return Chirp{}, ErrNotExist
	}

	digitUserId, err := strconv.Atoi(userId)

	if err != nil {
		return Chirp{}, err
	}

	if chirp.AuthorId == digitUserId {
		delete(dbStructure.Chirps, id)
		err = db.writeDB(dbStructure)
		if err != nil {
			return Chirp{}, err
		}
	} else {
		return Chirp{}, fmt.Errorf("not allowed")
	}

	return chirp, nil
}
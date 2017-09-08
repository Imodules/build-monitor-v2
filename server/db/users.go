package db

import (
	"errors"

	"time"

	"github.com/elithrar/simple-scrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	DbObject    `bson:",inline"`
	Username    string    `bson:"username" json:"username"`
	Email       string    `bson:"email" json:"email"`
	Password    string    `bson:"password" json:"-"`
	Token       string    `bson:"-" json:"token"`
	LastLoginAt time.Time `bson:"lastLoginAt" json:"lastLoginAt"`
}

var DuplicateUser = errors.New("A user already exists with this username or email")
var MissingUserField = errors.New("The user is missing a required parameter")
var UserNotFound = errors.New("A user matching those credentials could not be found")

func Users(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("users")
}

func (appDb *AppDb) CreateUser(username, email, password string) (*User, error) {
	if err := validateUser(username, email, password); err != nil {
		return nil, err
	}

	hashedPassword, hashErr := hash(password, appDb.PasswordSalt)
	if hashErr != nil {
		return nil, hashErr
	}

	user := User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	setCreated(&user.DbObject)
	user.LastLoginAt = user.CreatedAt

	err := Users(appDb.Session).Insert(user)
	if err != nil {
		if mgo.IsDup(err) {
			return nil, DuplicateUser
		}

		return nil, err
	}

	return &user, nil
}

func (appDb *AppDb) FindUserByLogin(usernameOrEmail string, password string) (*User, error) {
	var user User

	err := Users(appDb.Session).Find(bson.M{"$or": []bson.M{{"username": usernameOrEmail}, {"email": usernameOrEmail}}}).One(&user)
	if err != nil {
		return nil, UserNotFound
	}

	if err := validateHash(password, appDb.PasswordSalt, user.Password); err != nil {
		return nil, UserNotFound
	}

	return &user, nil
}

func (appDb *AppDb) FindUserById(id string) (*User, error) {
	var user User

	err := Users(appDb.Session).FindId(bson.ObjectIdHex(id)).One(&user)
	if err != nil {
		return nil, UserNotFound
	}

	return &user, nil
}

func (appDb *AppDb) LogUserLogin(user *User) {
	err := Users(appDb.Session).UpdateId(user.Id, bson.M{"$set": bson.M{"lastLoginAt": getNow()}})
	if err != nil {
		appDb.Log.Errorf("Failed to update LastLoginAt for userId: %s", user.Id.Hex())
	}
}

func validateUser(username, email, password string) error {
	if username == "" || email == "" || password == "" {
		return MissingUserField
	}

	return nil
}

//region Hashing
func hash(password, salt string) (string, error) {
	saltedPassword := saltPassword(password, salt)

	hash, err := scrypt.GenerateFromPassword([]byte(saltedPassword), scrypt.DefaultParams)
	if err != nil {
		return "", err
	}

	s := string(hash)

	return s, nil
}

var validateHash = func(password, salt, hash string) error {
	saltedPassword := saltPassword(password, salt)

	return scrypt.CompareHashAndPassword([]byte(hash), []byte(saltedPassword))
}

func saltPassword(password, salt string) string {
	return salt + password + "-secret-sauce"
}

//endregion

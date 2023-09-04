package database

import (
	"app-auth/pkg/models"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var (
	mongoURI     = "mongodb+srv://oguzhancart1:8En6lGLMS3PZzqsI@cluster0.holjrh9.mongodb.net/"
	mongoClient  *mongo.Client
	databaseName = "app-db"
	secretKey    = "983yıthejh78f4ı5uehfgkey"
)

func InitMongoClient() {
	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB ping error: %v", err)
	}

	mongoClient = client
	log.Println("Connected to MongoDB")
	initialUser := models.User{
		Username: "admin",
		Password: "123456789",
		Roles:    []string{"admin"},
	}
	existingUser, _ := FindUserByUsername(initialUser.Username)
	if existingUser == nil {
		err := InsertUser(&initialUser)
		if err != nil {
			log.Fatalf("Başlangıçta kullanıcı eklerken hata: %v", err)
		}
	}
}

func GetMongoClient() *mongo.Client {
	return mongoClient
}

func AuthenticateUser(user *models.User) (bool, error) {
	collection := mongoClient.Database(databaseName).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"username": user.Username}
	var result models.User
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
	if err != nil {
		return false, nil
	}
	return true, nil
}

func InsertUser(user *models.User) error {
	collection := mongoClient.Database(databaseName).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"username": user.Username}
	var existingUser models.User
	err := collection.FindOne(ctx, filter).Decode(&existingUser)
	if err == nil {
		return fmt.Errorf("Username already exists")
	} else if err != mongo.ErrNoDocuments {
		return err
	}

	user.ID = primitive.NewObjectID()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Error adding user: %v\n", err)
		return err
	}

	return nil
}

func GenerateToken(user *models.User) (string, error) {
	if user.Roles == nil || len(user.Roles) == 0 {
		user.Roles = []string{"admin"}
	}

	claims := jwt.MapClaims{
		"username": user.Username,
		"roles":    user.Roles,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func FindUserByUsername(username string) (*models.User, error) {
	collection := mongoClient.Database(databaseName).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"username": username}
	var user models.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetSecretKey() string {
	return secretKey
}

package users

import (
	"C2S/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Store) CreateUser(user *models.User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	_, err := s.collection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	return nil
}


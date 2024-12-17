package user

import (
	"C2S/internal/models"
	"C2S/internal/services/auth"
	"C2S/internal/types"
	"C2S/internal/utils"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)


func (h *Handler) HandleRegister(c *fiber.Ctx) error{
	var payload types.RegisterUserPayload

	if err := utils.ParseJSON(c, &payload); err != nil {
		log.Println("hi5")
		return utils.WriteError(c, fiber.StatusBadRequest, err)
	}
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return utils.WriteError(c, fiber.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		
	}
	existingUser, err := h.store.GetUserByUserName(payload.UserName)
	if err == nil && existingUser != nil {
		log.Println("h3i")
		return utils.WriteError(c, fiber.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.UserName))
	}
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println("h2i")
		return utils.WriteError(c, fiber.StatusInternalServerError, err)
	}


	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		log.Println("hi1")
		return utils.WriteError(c, fiber.StatusInternalServerError, err)
	}


	user := models.User{
		UserName:  payload.UserName,
		Password:  hashedPassword,
		IsAdmin:   false,
	}


	if err := h.store.CreateUser(&user); err != nil {
		log.Println("hi8")
		return utils.WriteError(c, fiber.StatusInternalServerError, err)
	}

	if err := h.store.SeedQuestionsForUser(c.Context(),user.ID); err != nil {
		log.Println("hi95")
		return utils.WriteError(c, fiber.StatusInternalServerError, fmt.Errorf("failed to seed questions for user: %v", err))
	}

	return utils.WriteJSON(c, fiber.StatusCreated, "Success")
}
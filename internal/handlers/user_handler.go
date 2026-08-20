package handlers

import (
	"net/http"
	"time"
	"track-your-game-api/internal/models"       //modeles
	"track-your-game-api/internal/repositories" //dépôt

	"github.com/gin-gonic/gin"   //requête http
	"github.com/google/uuid"     //uuid
	"golang.org/x/crypto/bcrypt" //cryptage
)

type UserHandler struct {
	userRepo *repositories.UserRepository
}

func newUserHandler(userRepo *repositories.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest

	//vérification données sous format JSON de la requête
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	//hash du mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "Erreur lors du cryptage du mot de passe"})
		return
	}

	//création de l'utilisateur temporaire
	user := &models.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.userRepo.CreateUser(c.Request.Context(), user); err != nil { //insertion en bdd par la fonction du dépôt
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "Erreur lors de la création de l'utilisateur"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{ //retour json dans la requête
		"id":        user.ID,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "Donnée invalides"})
		return
	}

	//récupérer l'utilisateur par email
	user, err := h.userRepo.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "Erreur serveur"})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mdp incorrect"})
		return
	}

	//vérification du mdp
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"erreur": "Email ou mdp incorrect"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
	})
}

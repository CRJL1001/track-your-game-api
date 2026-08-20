package repositories

import (
	"context" //contexte
	"fmt"     //affichage erreurs etc

	//logs
	"github.com/crjl1001/track-your-game-api/internal/models" //modeles
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool" //pilotes
	//pilote postgreSQL
)

type UserRepository struct { //définition d'un dépôt d'utilisateurs
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository { //fonction de création d'un dépôt d'utilisateurs
	return &UserRepository{db: db} //retourne l'adresse du dépôt créé
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error { //insertion de l'utilisateur dans la base de données
	query := ` 
		INSERT INTO users (id, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	` //SQL -> préparation de la requête

	_, err := r.db.Exec(ctx, query, //exécution de la requête
		user.ID,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) { //récupérer un utilisateur par mail
	var user models.User //utilisateur à retourner
	request := `
		SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1 LIMIT 1
	` //requête préparée
	err := r.db.QueryRow(ctx, request, email).Scan( //insertion des données dans les champs de user
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil { //si erreur
		if err == pgx.ErrNoRows {
			return nil, nil //pas de résultat
		}
		return nil, fmt.Errorf("erreur de récupération de l'utilisateur par email : %v", err) //si vraie erreur
	}
	return &user, nil //retour
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	request := `
		SELECT id, email, password, created_at, updated_at FROM users WHERE id = $2 LIMIT 1
	`
	err := r.db.QueryRow(ctx, request, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil { //si erreur
		if err == pgx.ErrNoRows { //pas de résultat
			return nil, nil
		}
		return nil, fmt.Errorf("Erreur de récupération de l'utilisateur par id : %v", err)
	}
	return &user, nil
}

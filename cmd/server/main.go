package main

import (
	"context"   //contexte
	"log"       //log
	"net/http"  //requête http
	"os"        //os
	"os/signal" //gestion des signaux système
	"syscall"   //constante des signaux système
	"time"      //temps

	"track-your-game-api/internal/config"
	"track-your-game-api/internal/handlers"
	"track-your-game-api/internal/repositories"

	"github.com/crjl1001/track-your-game-api/internal/config"
	"github.com/gin-gonic/gin"        // requête http
	"github.com/jackc/pgx/v5/pgxpool" //pilote POSTEGRESQL
)

func main() {
	//chargement de la configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erreur lors du chargement de la configuration : %v", err)
	}

	//connexion à PostgreSQL
	dbPool, err := pgxpool.New(context.Background(), buildDSN(cfg))
	if err != nil {
		log.Fatalf("impossible de se connecter à la base de donnée : %v", err)
	}
	defer dbPool.Close() //la connexion se fermera quand la fonction se terminera

	//test de connexion
	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Echec du ping à la base de données : %v", err)
	}
	
	log.Println("Connecté à PostegreSQL")

	//initialisation des dépôts et handlers
	userRepo := repositories.newUserRepository(dbPool)
	userHandler := handlers.newUserHandler(userRepo)

	//initialisation de Gin (framework HTTP)
	router := gin.Default()

	//Routes
	api := router.Group("/api/v1"){
		api.POST("/users", userHandler.CreateUser)
		api.POST("/users/login", UserHandler.Login)
	}

	//Démarrer le serveur
	srv := &http.Server{
		Addr: ":" + cfg.Port,
		Handler: router,
	}

	//Gestion des signaux pour un arrêt propre
	go func(){ //goroutine -> ne bloque pas l'exécution du main
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { //écoute les erreurs du serveur
			log.Fatalf("Erreur de serveur : %v", err) //si erreur autre que fermeture intentionnele on arrête tout
		}
	}()

	log.Printf("Serveur démarré sur le port %s", cfg.Port) //message

	//Attente de signaux SUGINT ou SIGTERM pour arrpêter le serveur
	quit := make(chan os.Signal, 1) //on créer un canal quit 
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //on envoie un signal dans le canal quit lors de la réception des signaux SIGINT et SIGTERM
	<-quit //bloque l'éxécution sur ce canal jusqu'a recevoir un signal
	log.Println("Arrêt du serveur...") //continuer l'éxécution 

	//arrêt avec timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //donne un context qui donne 10 secondes max pour terminer les requêtes
	defer cancel() //Annule le contexte à la fin de l'exécution, attend la fin des requêtes
	if err := srv.Shutdown(ctx); err != nil { //ferme le serveur proprement
		log.Fatalf("Erreur lors de l'arrêt du serveur : %v", err) 
	}
	log.Println("Serveur arrêté")
}

func buildDSN(cfg *config.Config) string {
	return "postegres://" + cfg.DBUser + ":" + cfg.DBPassword + "@" + cfg.DBHost + ":" + cfg.DBPort + "/" + cfg.DBName + "?sslmode=" + cfg.DBSSLMode
}







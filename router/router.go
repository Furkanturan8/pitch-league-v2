package router

import (
	"time"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/personal-project/pitch-league/handlers"
	"github.com/personal-project/pitch-league/middleware"
	"github.com/personal-project/pitch-league/repository"
	"github.com/uptrace/bun"
)

type Config struct {
	JWTSecret              string
	AccessTokenExpireTime  int
	RefreshTokenExpireTime int
}

func Setup(app fiber.Router, db *bun.DB, cfg Config) {
	app.Use(logger.New())
	app.Use(cors.New())

	// Swagger konfigürasyonu
	app.Get("/swagger/*", swagger.New(swagger.Config{
		Title:        "Pitch League API",
		DeepLinking:  true,
		DocExpansion: "list",
		OAuth: &swagger.OAuthConfig{
			AppName:  "Pitch League",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		OAuth2RedirectUrl: "/swagger/oauth2-redirect.html",
	}))

	api := app.Group("/api")

	// Repository'leri oluştur
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db, cfg.JWTSecret, time.Duration(cfg.AccessTokenExpireTime)*time.Hour, time.Duration(cfg.RefreshTokenExpireTime)*time.Hour)
	teamRepo := repository.NewTeamRepository(db)
	fieldRepo := repository.NewFieldRepository(db)
	gameRepo := repository.NewGameRepository(db)
	gamePartRepo := repository.NewGameParticipantsRepository(db)
	leagueRepo := repository.NewLeagueRepository(db)
	leagueTeamRepo := repository.NewLeagueTeamRepository(db)
	matchRepo := repository.NewMatchRepository(db)

	// Handler'ları oluştur
	authHandler := handlers.NewAuthHandler(authRepo, userRepo, time.Duration(cfg.RefreshTokenExpireTime)*time.Hour, cfg.JWTSecret)
	userHandler := handlers.NewUserHandler(userRepo)
	teamHandler := handlers.NewTeamHandler(teamRepo)
	fieldHandler := handlers.NewFieldHandler(fieldRepo)
	gameHandler := handlers.NewGameHandler(gameRepo)
	gamePartHandler := handlers.NewGameParticipantsHandler(gamePartRepo)
	leagueHandler := handlers.NewLeagueHandler(leagueRepo)
	leagueTeamHandler := handlers.NewLeagueTeamHandler(leagueTeamRepo)
	matchHandler := handlers.NewMatchHandler(matchRepo)

	// Public routes
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Post("/logout", authHandler.Logout)

	// Protected routes
	api.Use(middleware.JWTMiddleware(cfg.JWTSecret))

	// Game Participants routes
	gameParts := api.Group("/gameParts")
	gameParts.Get("/", gamePartHandler.GetAllGameParticipants)     // takımların oyun ile ilişkilerini getirir
	gameParts.Get("/:id", gamePartHandler.GetByGameParticipantsID) // kullanıcı idsine göre oyun ilişkilerini getirir

	// League routes
	leagues := api.Group("/leagues")
	leagues.Get("/", leagueHandler.GetAllLeagues)    // tüm ligleri getirir
	leagues.Get("/:id", leagueHandler.GetByLeagueID) // id ye göre belli bir ligi getirir

	// League Team routes
	leagueTeams := api.Group("/leaguesTeam")
	leagueTeams.Get("/", leagueTeamHandler.GetAllLeagueTeams)       // tüm liglerdeki takımları puan sıralamasına göre getirir
	leagueTeams.Get("/:id", leagueTeamHandler.GetByLeagueTeamID)    // team id sine göre takımı getirir
	leagueTeams.Get("/league/:id", leagueTeamHandler.GetByLeagueID) // lig id sine göre ligdeki tüm takımları getirir

	// Match routes
	matches := api.Group("/matches")
	matches.Get("/", matchHandler.GetAllMatches)   // tüm maçları getirir
	matches.Get("/:id", matchHandler.GetByMatchID) // gameID ye göre o maçı getirir

	// Team routes
	teams := api.Group("/teams")
	teams.Get("/", teamHandler.GetAllTeams)
	teams.Get("/:id", teamHandler.GetByTeamID)
	teams.Post("/:id/join/:userID", teamHandler.JoinTeam)

	// Field routes
	fields := api.Group("/fields")
	fields.Get("/", fieldHandler.GetAllFields)
	fields.Get("/:id", fieldHandler.GetByFieldID)

	// Game routes
	games := api.Group("/games")
	games.Get("/", gameHandler.GetAllGames)
	games.Get("/:id", gameHandler.GetByGameID)

	// Admin routes
	adminRoutes := api.Group("/admin")
	adminRoutes.Use(middleware.AdminControl)

	// Admin User routes
	adminUsers := adminRoutes.Group("/users")
	adminUsers.Post("/", userHandler.CreateUser)
	adminUsers.Get("/", userHandler.GetAllUsers)
	adminUsers.Get("/:id", userHandler.GetByUserID)
	adminUsers.Delete("/:id", userHandler.DeleteByUserID)
	adminUsers.Put("/:id", userHandler.UpdateUserByID)

	// Admin Match routes
	adminMatches := adminRoutes.Group("/matches")
	adminMatches.Post("/", matchHandler.CreateMatch)          // maç oluşturur
	adminMatches.Delete("/:id", matchHandler.DeleteByMatchID) // maçı iptal eder

	// Admin League routes
	adminLeagues := adminRoutes.Group("/leagues")
	adminLeagues.Post("/", leagueHandler.CreateLeague)          // yeni bir yerel lig oluşturur
	adminLeagues.Delete("/:id", leagueHandler.DeleteByLeagueID) // ligi siler

	// Admin League Team routes
	adminLeagueTeams := adminRoutes.Group("/leagueTeams")
	adminLeagueTeams.Post("/", leagueTeamHandler.CreateLeagueTeam)          // yeni oluşturulan takımı belirli bir lige kaydeder
	adminLeagueTeams.Delete("/:id", leagueTeamHandler.DeleteByLeagueTeamID) // takımı ligden siler

	// Admin Game Participants routes
	adminGameParts := adminRoutes.Group("/gameParts")
	adminGameParts.Get("/users/:id", gamePartHandler.GetGameParticipantsUsers) // game idsine göre maça katılan tüm kullanıcıları getirir
	adminGameParts.Post("/", gamePartHandler.CreateGameParticipants)           // gamePart ekler
	adminGameParts.Delete("/:id", gamePartHandler.DeleteByGameParticipantsID)  // gamePart siler

	// Admin Team routes
	adminTeams := adminRoutes.Group("/teams")
	adminTeams.Post("/", teamHandler.CreateTeam)
	adminTeams.Delete("/:id", teamHandler.DeleteByTeamID)
	adminTeams.Put("/:id", teamHandler.UpdateTeamByID)

	// Admin Field routes
	adminFields := adminRoutes.Group("/fields")
	adminFields.Post("/", fieldHandler.CreateField)
	adminFields.Delete("/:id", fieldHandler.DeleteByFieldID)
	adminFields.Put("/:id", fieldHandler.UpdateFieldByID)

	// Admin Game routes
	adminGames := adminRoutes.Group("/games")
	adminGames.Post("/", gameHandler.CreateGame)
	adminGames.Delete("/:id", gameHandler.DeleteByGameID)
	adminGames.Put("/:id", gameHandler.UpdateGameByID)
}

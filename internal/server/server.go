package server

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"spy_cat_agency/internal/database"
	"spy_cat_agency/internal/server/controllers"
	"spy_cat_agency/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Server struct {
	router            *gin.Engine
	catController     controllers.CatController
	missionController controllers.MissionController
	infoLog           *log.Logger
	errorLog          *log.Logger
}

func NewServer() *Server {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db := initDB(errorLog)

	catRepo := database.NewCatRepository(db)
	catService := services.NewCatService(catRepo)
	catController := controllers.NewCatController(*catService, errorLog)

	missionRepo := database.NewMissionRepository(db)
	missionService := services.NewMissionService(missionRepo)
	missinController := controllers.NewMissionController(*missionService, errorLog)

	server := &Server{
		router:            gin.Default(),
		catController:     *catController,
		missionController: *missinController,
		infoLog:           infoLog,
		errorLog:          errorLog,
	}

	server.runDBMigration()
	server.setupRoutes()
	server.AddBreedValidator()

	return server
}

func (s *Server) setupRoutes() {
	catRoutes := s.router.Group("/cat")
	catRoutes.POST("/add", s.catController.HireCat)
	catRoutes.DELETE("/delete", s.catController.FireCat)
	catRoutes.GET("/list", s.catController.ListCats)
	catRoutes.GET("/get", s.catController.GetCat)
	catRoutes.PATCH("/updateSalary", s.catController.UpdateSalary)

	missionRoutes := s.router.Group("/mission")
	missionRoutes.POST("/add", s.missionController.AddMission)
	missionRoutes.PATCH("/assign", s.missionController.Assign)
	missionRoutes.GET("/get", s.missionController.GetMission)
	missionRoutes.DELETE("/delete", s.missionController.DeleteMission)
	missionRoutes.GET("/list", s.missionController.ListMissions)
	missionRoutes.PATCH("/update", s.missionController.UpdateMission)

	targetRoutes := s.router.Group("target")
	targetRoutes.GET("/get", s.missionController.GetTarget)
	targetRoutes.DELETE("/delete", s.missionController.DeleteTarget)
	targetRoutes.POST("/add", s.missionController.AddTarget)
	targetRoutes.PATCH("/complete", s.missionController.CompleteTarget)
	targetRoutes.PATCH("/updateNotes", s.missionController.UpdateTargetNotes)

}

func initDB(errorLog *log.Logger) *sql.DB {
	conn, err := sql.Open("postgres", os.Getenv("DB_SOURCE"))
	if err != nil {
		errorLog.Fatal(err)

	}
	return conn
}

func (s *Server) runDBMigration() {
	migration, err := migrate.New(os.Getenv("MIGRATION_PATH"), os.Getenv("DB_SOURCE"))
	if err != nil {
		s.errorLog.Fatal("cannot create migration:", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up:", err)
	}

	s.infoLog.Println("DB migrated successfully")

}

func (s *Server) Run() {
	port := os.Getenv("SERVER_PORT")
	if err := s.router.Run(":" + port); err != nil {
		s.errorLog.Fatalf("Failed to run server %v", err.Error())
	}
}

type BreedName struct {
	Name string `json:"name"`
}

func (s *Server) AddBreedValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// registering validation for nontoneof
		v.RegisterValidation("breed", func(fl validator.FieldLevel) bool {
			var breedList []BreedName
			url := "https://api.thecatapi.com/v1/breeds"

			queryRes, err := http.Get(url)
			if err != nil {
				s.errorLog.Println("Failed to fetch breed info.")
				return false
			}

			body, err := io.ReadAll(queryRes.Body)
			if err != nil {
				s.errorLog.Println("Failed to fetch breed info.")
				return false
			}

			err = json.Unmarshal(body, &breedList)
			if err != nil {
				s.errorLog.Println("Failed to fetch breed info.")
				return false
			}

			if breed, ok := fl.Field().Interface().(string); ok {
				for _, v := range breedList {
					if v.Name == breed {
						return true
					}
				}
			}

			return false
		})
	}
}

package main

import (
	"log"
	"net/http"

	"github.com/deividroger/api-go/configs"
	"github.com/deividroger/api-go/internal/entity"
	"github.com/deividroger/api-go/internal/infra/database"
	handlers "github.com/deividroger/api-go/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"

	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/deividroger/api-go/docs"
)

// @title           Swagger Example API (Product Service)
// @version         1.0
// @description     Product API with authentication.
// @termsOfService

// @contact.name   Deivid Roger Oliveira Santos
// @contact.url
// @contact.email

// @license.name
// @license.url

// @host      localhost:8000
// @BasePath  /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	configs, err := configs.LoadConfig(".")

	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.User{}, &entity.Product{})

	productDb := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDb)

	userDb := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDb)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(LogRequest)
	r.Use(middleware.Recoverer)

	r.Use(middleware.WithValue("jwt", configs.TokenAuth))
	r.Use(middleware.WithValue("expiresIn", configs.JWTExpiresIn))

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.TokenAuth)) //valida se o token é válido
		r.Use(jwtauth.Authenticator)               //verifica se o token é valido

		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.GetProducts)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Post("/generateToken", userHandler.GetJWT)
	})

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))

	http.ListenAndServe(":8000", r)

}

func LogRequest(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

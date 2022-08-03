package main

import (
	"github.com/Adebayobenjamin/clean_arch/config"
	"github.com/Adebayobenjamin/clean_arch/controller"
	"github.com/Adebayobenjamin/clean_arch/middleware"
	"github.com/Adebayobenjamin/clean_arch/repository"
	"github.com/Adebayobenjamin/clean_arch/service"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB                  = config.SetupDatabaseConnection()
	userRepository repository.UserRepository = repository.NewUserRepository(db)
	bookRepository repository.BookRepository = repository.NewBookRepository(db)
	jwtService     service.JWTService        = service.NewJWTService()
	authService    service.AuthService       = service.NewAuthService(userRepository)
	userService    service.UserService       = service.NewUserService(userRepository)
	bookService    service.BookService       = service.NewBookService(bookRepository)
	authController controller.AuthController = controller.NewAuthController(authService, jwtService)
	userController controller.UserController = controller.NewUserController(userService, jwtService)
	bookController controller.BookController = controller.NewBookController(bookService, jwtService)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()
	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/update", userController.Update)
	}

	bookRoutes := r.Group("api/books", middleware.AuthorizeJWT(jwtService))
	{
		bookRoutes.GET("/", bookController.AllBooks)
		bookRoutes.POST("/", bookController.Insert)
		bookRoutes.GET("/:id", bookController.FindByID)
		bookRoutes.PUT("/:id", bookController.Update)
		bookRoutes.DELETE("/:id", bookController.Delete)
	}

	r.Run(":2202")
}

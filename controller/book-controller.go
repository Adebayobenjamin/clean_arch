package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Adebayobenjamin/clean_arch/dto"
	"github.com/Adebayobenjamin/clean_arch/entity"
	"github.com/Adebayobenjamin/clean_arch/helper"
	"github.com/Adebayobenjamin/clean_arch/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// BookController is a contract that explains what this controller can do
type BookController interface {
	AllBooks(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	Insert(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(context *gin.Context)
}

type bookController struct {
	bookService service.BookService
	jwtService  service.JWTService
}

// NewBookController creates a new instance of book controller
func NewBookController(bookserv service.BookService, jwtServ service.JWTService) BookController {
	return &bookController{
		bookService: bookserv,
		jwtService:  jwtServ,
	}
}

func (c *bookController) AllBooks(ctx *gin.Context) {
	var books []entity.Book = c.bookService.AllBooks()
	res := helper.BuildResponse(true, "OK!", books)
	ctx.JSON(http.StatusOK, res)
}

func (c *bookController) FindByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param(("id")), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("No param id found", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	var book entity.Book = c.bookService.FindByID(id)
	if (book == entity.Book{}) {
		res := helper.BuildErrorResponse("Data not found", "no data with the given id", helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	} else {
		res := helper.BuildResponse(true, "OK!", book)
		ctx.JSON(http.StatusOK, res)
	}
}

func (c *bookController) Insert(ctx *gin.Context) {
	var bookCreateDTO dto.BookCreateDTO
	errDTO := ctx.ShouldBind(&bookCreateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := ctx.GetHeader("Authorization")
		userID := c.getUserIDByToken(authHeader)
		convertedUserId, err := strconv.ParseUint(userID, 10, 64)
		if err == nil {
			bookCreateDTO.UserID = convertedUserId
			result := c.bookService.Insert(bookCreateDTO)
			res := helper.BuildResponse(true, "OK!", result)
			ctx.JSON(http.StatusCreated, res)
		} else {
			res := helper.BuildErrorResponse("Could not parse to uint", err.Error(), helper.EmptyObj{})
			ctx.JSON(http.StatusCreated, res)
		}
	}

}

func (c *bookController) Update(ctx *gin.Context) {
	var bookUpdateDTO dto.BookUpdateDTO
	errDTO := ctx.ShouldBind(&bookUpdateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to proceess request", errDTO.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.bookService.IsAllowedToEdit(userID, bookUpdateDTO.ID) {
		id, errID := strconv.ParseUint(userID, 10, 64)
		if errID == nil {
			bookUpdateDTO.UserID = id
		}
		result := c.bookService.Update(bookUpdateDTO)
		res := helper.BuildResponse(true, "OK!", result)
		ctx.JSON(http.StatusOK, res)
	} else {
		res := helper.BuildErrorResponse("You do not have permission", "you are not the owner", helper.EmptyObj{})
		ctx.JSON(http.StatusForbidden, res)
	}
}

func (c *bookController) Delete(ctx *gin.Context) {
	var book entity.Book
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("could not get id", "no param id was found", helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
	}
	book.ID = id
	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.bookService.IsAllowedToEdit(userID, book.ID) {
		c.bookService.Delete(book)
		res := helper.BuildResponse(true, "Deleted", helper.EmptyObj{})
		ctx.JSON(http.StatusOK, res)
	} else {
		res := helper.BuildErrorResponse("You do not have permission", "you are not the owner", helper.EmptyObj{})
		ctx.JSON(http.StatusForbidden, res)
	}
}

func (c *bookController) getUserIDByToken(token string) string {
	aToken, err := c.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}
	claims := aToken.Claims.(jwt.MapClaims)
	return fmt.Sprintf("%v", claims["user_id"])
}

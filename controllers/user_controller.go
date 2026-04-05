package controllers

import (
	"errandify/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func(u *UserController) Login(ctx *gin.Context) {
	user := models.User{}
	errBindJson := ctx.ShouldBindJSON(&user)
	if errBindJson != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errBindJson.Error()})
		return
	}

	password := user.Password

	errDB := u.DB.Where("email=?", user.Email).Take(&user).Error
	if errDB != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "email not found, contact admoon"})
		return
	}
	
	bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)

	errHash := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)

	if errHash != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error":"email or password is wrong"})
		return
	}

	ctx.JSON(http.StatusOK, user)
	

}

func(u *UserController) CreateAccount(ctx *gin.Context) {
	user := models.User{}
	errBindJson := ctx.ShouldBindJSON(&user)
	if errBindJson != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errBindJson.Error()})
		return
	}

	emailExists := u.DB.Where("email=?", user.Email).First(&user).RowsAffected != 0
	if emailExists == true {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}


	hashedPasswordBytes, errHash := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if errHash != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":errHash.Error()})
		return
	}

	user.Password = string(hashedPasswordBytes)
	user.Role = "Employee"

	errDB := u.DB.Create(&user).Error
	if errDB != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":errDB.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)

}

func(u *UserController) DeleteAccount(ctx *gin.Context) {
	id := ctx.Param("id")

	errDB := u.DB.Delete(&models.User{},id).Error
	if errDB != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":errDB.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "user deleted successfully"})

}
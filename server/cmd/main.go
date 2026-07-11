package main

import (
	"context"
	"fmt"
	"log"
	"server/core"
	"server/internal/application/auth"
	"server/internal/application/category"
	"server/internal/application/document"
	"server/internal/infracstructure"
	"server/internal/interfaces"
	"server/internal/security"

	"github.com/gin-gonic/gin"
)

func main() {
	config := core.LoadConfig()
	fmt.Println(config.MinioEndpoint, config.MinioAccessKey, config.MinioSecretKey)

	// connect gorm
	db, err := core.ConnectDB(config.DBURL)
	if err != nil {
		panic(err)
	}

	// connect minIO
	client, err := core.ConnectMinIO(*config)
	if err != nil {
		panic(err)
	}

	// repository
	userGormRepository := infracstructure.NewUserRepository(db)
	categoryGormRepository := infracstructure.NewCategoryRepository(db)
	documentStorageRepository := infracstructure.NewMinIOStorageDocument(client, "documents")
	documentRepository := infracstructure.NewGormDocumentRepository(db)
	storageRepository := infracstructure.NewStorageRepository(db)
	ragRepository := infracstructure.NewPythonIndexClientRepository(config.RagURL)

	// create document bucker
	err = documentStorageRepository.CreateBucket(context.Background())
	if err != nil {
		panic(err)
	}

	// security
	passwordHasher := security.NewBcryptHaser()
	jwtGenerator := security.NewJwtGenerate(config.Secretkey)

	// application
	authApplication := auth.NewAuthService(userGormRepository, storageRepository, *passwordHasher, *jwtGenerator)
	categoryApplication := category.NewCategoryService(categoryGormRepository, userGormRepository)
	documentApplication := document.NewDocumentService(
		documentStorageRepository,
		userGormRepository,
		categoryGormRepository,
		storageRepository,
		ragRepository,
		documentRepository,
	)

	// handler
	authHandler := interfaces.NewAuthHandler(authApplication)
	categoryHandler := interfaces.NewCategoryHandler(categoryApplication)
	documentHandler := interfaces.NewDocumentHandler(documentApplication)

	r := gin.New()

	auth := r.Group("/api/auth")
	{
		// đăng ký
		fmt.Println("Request: Đăng ký")
		auth.POST("/register", authHandler.Register)
		// đăng nhập
		fmt.Println("Request: Đăng nhập")
		auth.POST("/sign-in", authHandler.SignIn)
	}

	category := r.Group("/api/categories")
	{
		// tạo category
		fmt.Println("Request: Tạo category")
		category.POST("", categoryHandler.CreateCategory)

		// lấy danh sách category của bản thân
		fmt.Println("Request: Lấy danh sách category")
		category.GET("", categoryHandler.GetCategoryList)
	}

	document := r.Group("/api/documents")
	{
		// upload file
		fmt.Println("Request: Upload file")
		document.POST("", documentHandler.UploadFile)
		// xoá file
		fmt.Println("Request: Xoá file")
		document.DELETE("", documentHandler.DeleteDocument)
		// update status
		fmt.Println("Request: Upadate Status")
		document.PUT("/:id", documentHandler.UpdateStatus)
		// lấy danh sách document qua category
		fmt.Println("Request: Lấy danh sách document")
		document.GET("/:cateName", documentHandler.GetListFileByUserID)
	}

	if err := r.Run(":" + config.Port); err != nil {
		log.Fatal("Server down")
	}
}

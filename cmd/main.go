package main

import (
	auth_controller "main/internal/controllers/http/v1/auth"

	departament_controller "main/internal/controllers/http/v1/departament"
	qr_controller "main/internal/controllers/http/v1/qr_code"
	room_controller "main/internal/controllers/http/v1/room"
	room_type_controller "main/internal/controllers/http/v1/room_type"
	user_controller "main/internal/controllers/http/v1/user"
	auth_middleware "main/internal/middleware/auth"
	"main/internal/pkg/config"
	"main/internal/pkg/postgres"
	"main/internal/repository/postgres/departament"
	"main/internal/repository/postgres/qr"
	"main/internal/repository/postgres/room"
	roomtype "main/internal/repository/postgres/room_type"
	"main/internal/repository/postgres/user"
	"main/internal/services/auth"
	file_service "main/internal/services/file"

	auth_use_case "main/internal/usecase/auth"
	departament_use_case "main/internal/usecase/departament"
	qr_use_case "main/internal/usecase/qr_code"
	room_use_case "main/internal/usecase/room"
	room_type_use_case "main/internal/usecase/room_type"
	user_use_case "main/internal/usecase/user"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	serverPost := ":" + config.GetConfig().Port

	r := gin.Default()

	//databases
	postgresDB := postgres.NewDB()

	r.Static("/media", "./media")

	//repositories
	userRepository := user.NewRepository(postgresDB)
	roomRepository := room.NewRepository(postgresDB)
	qrCodeRepository := qr.NewRepository(postgresDB)
	departamentRepository := departament.NewRepository(postgresDB)
	roomTypeRepository := roomtype.NewRepository(postgresDB)

	//services
	authService := auth.NewService(userRepository)
	fileService := file_service.NewService()
	// videoService := video_service.NewService()
	// audioService := audio_service.NewService()

	//cache
	// newCache := cache.NewCache(config.GetConfig().RedisHost, config.GetConfig().RedisDB, time.Duration(config.GetConfig().RedisExpires)*time.Second)

	//usecase
	authUseCase := auth_use_case.NewUseCase(authService, userRepository)
	userUseCase := user_use_case.NewUseCase(userRepository, authService, fileService)
	roomUseCase := room_use_case.NewUseCase(roomRepository, authService)
	qrCodeUseCase := qr_use_case.NewUseCase(qrCodeRepository, fileService, authService)
	departamentUseCase := departament_use_case.NewUseCase(departamentRepository, authService)
	roomTypeUseCase := room_type_use_case.NewUseCase(roomTypeRepository, authService)

	//controller
	authController := auth_controller.NewController(authUseCase)
	userController := user_controller.NewController(userUseCase)
	roomController := room_controller.NewController(roomUseCase)
	qrCodeController := qr_controller.NewController(qrCodeUseCase)
	departamentController := departament_controller.NewController(departamentUseCase)
	roomTypeController := room_type_controller.NewController(roomTypeUseCase)

	//middleware
	authMiddleware := auth_middleware.NewMiddleware(authService)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://172.20.10.2:5173", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api")
	{
		v1 := api.Group("v1")

		// #auth

		//sign-in
		v1.POST("/sign-in", authController.SignIn)

		//	#user
		//list
		v1.GET("/admin/user/list", authMiddleware.AuthMiddleware(), userController.AdminGetUserList)
		// udpate
		v1.PATCH("/admin/user/update/:id", authMiddleware.AuthMiddleware(), userController.AdminUpdateUser)
		//	create
		v1.POST("/admin/user/create", authMiddleware.AuthMiddleware(), userController.AdminCreateUser)
		// //delete
		v1.DELETE("/admin/user/delete/:id", authMiddleware.AuthMiddleware(), userController.AdminDeleteUser)
		//	get by id
		v1.GET("/admin/user/:id", authMiddleware.AuthMiddleware(), userController.AdminGetById)

		// #room
		// create
		v1.POST("/admin/room/create", authMiddleware.AuthMiddleware(), roomController.AdminCreateRoom)
		// delete
		v1.DELETE("/admin/room/delete/:id", authMiddleware.AuthMiddleware(), roomController.AdminDeleteRoom)
		// update
		v1.GET("/admin/room/list", authMiddleware.AuthMiddleware(), roomController.AdminGetList)
		// get by id
		v1.GET("/admin/room/:id", authMiddleware.AuthMiddleware(), roomController.AdminGetById)

		// #qr_code
		// create
		v1.POST("/admin/generate/qr", authMiddleware.AuthMiddleware(), qrCodeController.AdminGenerateQRCode)

		// #departament
		// create
		v1.POST("/admin/departament/create", authMiddleware.AuthMiddleware(), departamentController.AdminCreateDepartament)
		// delete
		v1.DELETE("/admin/departament/delete/:id", authMiddleware.AuthMiddleware(), departamentController.AdminDeleteDepartament)
		// list
		v1.GET("/admin/departament/list", authMiddleware.AuthMiddleware(), departamentController.AdminGetDepartamentList)
		// get by id
		v1.GET("/admin/departament/:id", authMiddleware.AuthMiddleware(), departamentController.AdminGetDepartamentById)
		// update
		v1.PATCH("/admin/departament/update/:id", authMiddleware.AuthMiddleware(), departamentController.AdminUpdateDepartament)

		// #room_type
		// create
		v1.POST("/admin/room_type/create", authMiddleware.AuthMiddleware(), roomTypeController.AdminCreateRoomType)
		// delete
		v1.DELETE("/admin/room_type/delete/:id", authMiddleware.AuthMiddleware(), roomTypeController.AdminDeleteRoomType)
		// list
		v1.GET("/admin/room_type/list", authMiddleware.AuthMiddleware(), roomTypeController.AdminGetRoomTypeList)
		// get by id
		v1.GET("/admin/room_type/:id", authMiddleware.AuthMiddleware(), roomTypeController.AdminGetRoomTypeById)
		// update
		v1.PATCH("/admin/room_type/update/:id", authMiddleware.AuthMiddleware(), roomTypeController.AdminUpdateRoomType)

	}

	r.Run(serverPost)

}

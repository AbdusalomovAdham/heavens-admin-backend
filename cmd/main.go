package main

import (
	auth_controller "main/internal/controllers/http/v1/auth"

	department_controller "main/internal/controllers/http/v1/department"
	qr_controller "main/internal/controllers/http/v1/qr_code"
	role_controller "main/internal/controllers/http/v1/role"
	room_controller "main/internal/controllers/http/v1/room"
	room_type_controller "main/internal/controllers/http/v1/room_type"
	user_controller "main/internal/controllers/http/v1/user"
	work_status_controller "main/internal/controllers/http/v1/work_status"
	auth_middleware "main/internal/middleware/auth"
	"main/internal/pkg/config"
	"main/internal/pkg/postgres"
	department "main/internal/repository/postgres/department"
	"main/internal/repository/postgres/qr"
	"main/internal/repository/postgres/role"
	"main/internal/repository/postgres/room"
	roomtype "main/internal/repository/postgres/room_type"
	"main/internal/repository/postgres/user"
	work_status "main/internal/repository/postgres/work_status"
	"main/internal/services/auth"
	file_service "main/internal/services/file"

	auth_use_case "main/internal/usecase/auth"
	department_use_case "main/internal/usecase/department"
	qr_use_case "main/internal/usecase/qr_code"
	role_use_case "main/internal/usecase/role"
	room_use_case "main/internal/usecase/room"
	room_type_use_case "main/internal/usecase/room_type"
	user_use_case "main/internal/usecase/user"
	work_status_use_case "main/internal/usecase/work_status"
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
	departmentRepository := department.NewRepository(postgresDB)
	roomTypeRepository := roomtype.NewRepository(postgresDB)
	roleRepository := role.NewRepository(postgresDB)
	workStatusRepository := work_status.NewRepository(postgresDB)

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
	departmentUseCase := department_use_case.NewUseCase(departmentRepository, authService)
	roomTypeUseCase := room_type_use_case.NewUseCase(roomTypeRepository, authService)
	roleUseCase := role_use_case.NewUseCase(roleRepository, authService)
	workStatusUseCase := work_status_use_case.NewUseCase(workStatusRepository, authService)

	//controller
	authController := auth_controller.NewController(authUseCase)
	userController := user_controller.NewController(userUseCase)
	roomController := room_controller.NewController(roomUseCase)
	qrCodeController := qr_controller.NewController(qrCodeUseCase)
	departmentController := department_controller.NewController(departmentUseCase)
	roomTypeController := room_type_controller.NewController(roomTypeUseCase)
	roleController := role_controller.NewController(roleUseCase)
	workStatusController := work_status_controller.NewController(workStatusUseCase)

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

		// #department
		// create
		v1.POST("/admin/department/create", authMiddleware.AuthMiddleware(), departmentController.AdminCreateDepartment)
		// delete
		v1.DELETE("/admin/department/delete/:id", authMiddleware.AuthMiddleware(), departmentController.AdminDeleteDepartment)
		// list
		v1.GET("/admin/department/list", authMiddleware.AuthMiddleware(), departmentController.AdminGetDepartmentList)
		// get by id
		v1.GET("/admin/department/:id", authMiddleware.AuthMiddleware(), departmentController.AdminGetDepartmentById)
		// update
		v1.PATCH("/admin/department/update/:id", authMiddleware.AuthMiddleware(), departmentController.AdminUpdateDepartment)

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

		// #role
		// create
		v1.POST("/admin/role/create", authMiddleware.AuthMiddleware(), roleController.AdminCreateRole)
		// delete
		v1.DELETE("/admin/role/delete/:id", authMiddleware.AuthMiddleware(), roleController.AdminDeleteRole)
		// list
		v1.GET("/admin/role/list", authMiddleware.AuthMiddleware(), roleController.AdminGetRoleList)
		// get by id
		v1.GET("/admin/role/:id", authMiddleware.AuthMiddleware(), roleController.AdminGetRoleById)
		// update
		v1.PATCH("/admin/role/update/:id", authMiddleware.AuthMiddleware(), roleController.AdminUpdateRole)

		// #work_status
		// create
		v1.POST("/admin/work/status/create", authMiddleware.AuthMiddleware(), workStatusController.AdminCreateWorkStatus)
		// delete
		v1.DELETE("/admin/work/status/delete/:id", authMiddleware.AuthMiddleware(), workStatusController.AdminDeleteWorkStatus)
		// list
		v1.GET("/admin/work/status/list", authMiddleware.AuthMiddleware(), workStatusController.AdminGetWorkStatusList)
		// get by id
		v1.GET("/admin/work/status/:id", authMiddleware.AuthMiddleware(), workStatusController.AdminGetWorkStatusById)
		// update
		v1.PATCH("/admin/work/status/update/:id", authMiddleware.AuthMiddleware(), workStatusController.AdminUpdateWorkStatus)

	}

	r.Run(serverPost)

}

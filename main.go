package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"todolist/controller"
	"todolist/repository"
	"todolist/service"
)

func init() {
	// Загрузка .env файла
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// Чтение переменных окружения
	user := os.Getenv("db_user")
	host := os.Getenv("db_host")
	port := os.Getenv("db_port")
	dbname := os.Getenv("db_name")
	password := os.Getenv("db_password")

	// Подключение к базе данных
	db, err := connectToDatabase(user, host, port, dbname, password)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создание таблицы
	if err := createTable(db); err != nil {
		log.Fatal(err)
	}

	// Наполнение таблицы данными
	if err := insertData(db); err != nil {
		log.Fatal(err)
	}

	// Инициализация репозитория и сервиса
	todoRepository := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepository)

	// Инициализация контроллера
	todoController := controller.NewTodoController(todoService)

	// Конфигурация сервера
	if err := startServer(todoController); err != nil {
		log.Fatal(err)
	}
}

func connectToDatabase(user, host, port, dbname, password string) (*sql.DB, error) {
	dsn := fmt.Sprintf("user=%s host=%s port=%s dbname=%s sslmode=disable password=%s",
		user, host, port, dbname, password)
	return sql.Open("postgres", dsn)
}

func createTable(db *sql.DB) error {
	createTableSQL, err := os.ReadFile("sql/create_table.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(createTableSQL))
	return err
}

func insertData(db *sql.DB) error {
	insertDataSQL, err := os.ReadFile("sql/insert_data.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(insertDataSQL))
	return err
}

func startServer(todoController *controller.TodoController) error {
	router := httprouter.New()
	router.Handler("GET", "/static/swagger/docs/*filepath", http.StripPrefix("/static/swagger/docs/", http.FileServer(http.Dir("./docs"))))
	router.GET("/swagger/*filepath", swaggerHandler)
	router.POST("/create", todoController.Create)
	router.GET("/get/:id", todoController.Get)
	router.POST("/update/:id", todoController.Update)
	router.DELETE("/delete/:id", todoController.Delete)
	router.GET("/list", todoController.List)

	return http.ListenAndServe(":8080", router)
}

func swaggerHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	swaggerFileUrl := "http://localhost:8080/static/swagger/docs/swagger.json"
	handler := httpSwagger.Handler(httpSwagger.URL(swaggerFileUrl))
	handler.ServeHTTP(w, r)
}

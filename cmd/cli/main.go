package main

import (
	"log"
	"os"

	"github.com/LiFeAiR/crud-ai/internal/repository"
)

func main() {
	// Получаем строку подключения к БД из переменной окружения или используем значение по умолчанию
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "host=localhost port=5432 user=postgres password=password dbname=httpserverdb sslmode=disable"
	}

	// Создаем подключение к БД
	db, err := repository.NewDB(dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Создаем репозиторий пользователей
	userRepo := repository.NewUserRepository(db)
	orgRepo := repository.NewOrganizationRepository(db)
	permRepo := repository.NewPermissionRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	tariffRepo := repository.NewTariffRepository(db)

	// Инициализируем таблицы в БД
	err = userRepo.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Инициализируем таблицы в БД
	err = orgRepo.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Инициализируем таблицы в БД
	err = permRepo.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Инициализируем таблицы в БД
	err = roleRepo.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Инициализируем таблицы в БД
	err = tariffRepo.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	log.Println("Сli executed successfully")
}

package users_generation

import (
	"bufio"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"log"
	"os"
	"strings"
	"time"

	mrand "math/rand"

	"github.com/doug-martin/goqu/v9"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"socialNetworkOtus/internal/config"
	"socialNetworkOtus/internal/db"
)

type User struct {
	FirstName    string
	SecondName   string
	BirthDate    string
	Biography    string
	City         string
	PasswordHash string
}

func generateRandomPassword(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:length], nil
}

func GenerateUsers() {
	_ = godotenv.Load()
	cfg := config.Load()
	dbConn, err := db.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	// Генерируем один bcrypt-хеш для всех пользователей
	password := "password"
	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}
	passwordHash := string(passwordHashBytes)

	file, err := os.Open("internal/utils/users_generation/people.v2.csv")
	if err != nil {
		log.Fatalf("failed to open csv: %v", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Fatalf("failed to close csv: %v", err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	batchSize := 1000
	var users []User
	total := 0
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			log.Printf("skip invalid line: %s", line)
			continue
		}
		nameParts := strings.SplitN(parts[0], " ", 2)
		if len(nameParts) != 2 {
			log.Printf("skip invalid name: %s", parts[0])
			continue
		}
		secondName := nameParts[0]
		firstName := nameParts[1]
		birthDate := parts[1]
		city := parts[2]

		user := User{
			FirstName:    firstName,
			SecondName:   secondName,
			BirthDate:    birthDate,
			Biography:    "",
			City:         city,
			PasswordHash: passwordHash,
		}
		users = append(users, user)
		total++

		if len(users) == batchSize {
			insertBatch(dbConn, users)
			users = users[:0]
			log.Printf("Inserted %d users", total)
		}
	}
	if len(users) > 0 {
		insertBatch(dbConn, users)
		log.Printf("Inserted %d users (final batch)", total)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("scanner error: %v", err)
	}
	log.Printf("Done! Inserted total %d users", total)
}

func insertBatch(dbConn *goqu.Database, users []User) {
	records := make([]goqu.Record, len(users))
	for i, u := range users {
		records[i] = goqu.Record{
			"first_name":    u.FirstName,
			"second_name":   u.SecondName,
			"birth_date":    u.BirthDate,
			"biography":     u.Biography,
			"city":          u.City,
			"password_hash": u.PasswordHash,
		}
	}
	// исправление: преобразуем []goqu.Record в []interface{} для передачи в Rows(...)
	recordsIface := make([]interface{}, len(records))
	for i, rec := range records {
		recordsIface[i] = rec
	}
	_, err := dbConn.Insert("users").Rows(recordsIface...).Executor().Exec()
	if err != nil {
		log.Fatalf("failed to insert batch: %v", err)
	}
}

// GenerateUsersFillUp создает недостающих пользователей до 1_000_000
func GenerateUsersFillUp() {
	_ = godotenv.Load()
	cfg := config.Load()
	dbConn, err := db.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	// Считаем текущее количество пользователей
	var count int
	if sqlDB, ok := dbConn.Db.(*sql.DB); ok {
		err := sqlDB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
		if err != nil {
			return
		}
	} else {
		log.Fatalf("dbConn.Db не является *sql.DB")
	}
	missing := 1_000_000 - count
	if missing <= 0 {
		log.Printf("Пользователей уже достаточно: %d", count)
		return
	}
	log.Printf("Нужно добавить %d пользователей", missing)

	// Генерируем один bcrypt-хеш для всех пользователей
	password := "password"
	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}
	passwordHash := string(passwordHashBytes)

	// Читаем people.v2.csv в память
	file, err := os.Open("internal/utils/users_generation/people.v2.csv")
	if err != nil {
		log.Fatalf("failed to open csv: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("scanner error: %v", err)
	}
	if len(lines) == 0 {
		log.Fatalf("csv пустой")
	}

	batchSize := 1000
	var users []User
	total := 0
	mrand.Seed(time.Now().UnixNano())
	for i := 0; i < missing; i++ {
		line := lines[mrand.Intn(len(lines))]
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			log.Printf("skip invalid line: %s", line)
			continue
		}
		nameParts := strings.SplitN(parts[0], " ", 2)
		if len(nameParts) != 2 {
			log.Printf("skip invalid name: %s", parts[0])
			continue
		}
		secondName := nameParts[0]
		firstName := nameParts[1]
		birthDate := parts[1]
		city := parts[2]

		user := User{
			FirstName:    firstName,
			SecondName:   secondName,
			BirthDate:    birthDate,
			Biography:    "",
			City:         city,
			PasswordHash: passwordHash,
		}
		users = append(users, user)
		total++

		if len(users) == batchSize {
			insertBatch(dbConn, users)
			users = users[:0]
			log.Printf("Inserted %d users", total)
		}
	}
	if len(users) > 0 {
		insertBatch(dbConn, users)
		log.Printf("Inserted %d users (final batch)", total)
	}

	log.Printf("Done! Досоздано %d пользователей", total)
}

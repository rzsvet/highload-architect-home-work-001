package database

import (
	"api/internal/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Database struct {
	WriteDB *sql.DB
	ReadDB  *sql.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	writeDB, err := connectDB(
		cfg.WriteDBHost,
		cfg.WriteDBPort,
		cfg.WriteDBName,
		cfg.WriteDBUser,
		cfg.WriteDBPassword,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to write database: %w", err)
	}

	readDB, err := connectDB(
		cfg.ReadDBHost,
		cfg.ReadDBPort,
		cfg.ReadDBName,
		cfg.ReadDBUser,
		cfg.ReadDBPassword,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to read database: %w", err)
	}

	return &Database{
		WriteDB: writeDB,
		ReadDB:  readDB,
	}, nil
}

func connectDB(host, port, dbname, user, password string) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Printf("Connected to database at %s:%s", host, port)
	return db, nil
}

func (d *Database) Close() {
	if d.WriteDB != nil {
		d.WriteDB.Close()
	}
	if d.ReadDB != nil {
		d.ReadDB.Close()
	}
}

func InitTables(db *Database) error {
	// Create users table on write database
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(100) NOT NULL UNIQUE,
        email VARCHAR(255) NOT NULL UNIQUE,
        password VARCHAR(255) NOT NULL,
        first_name VARCHAR(100) NOT NULL,
        last_name VARCHAR(100) NOT NULL,
        birth_date DATE NOT NULL,
        gender VARCHAR(20) NOT NULL CHECK (gender IN ('male', 'female', 'unknown')),
        interests TEXT,
        city VARCHAR(100),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    
    CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
    CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
    CREATE INDEX IF NOT EXISTS idx_users_city ON users(city);
    CREATE INDEX IF NOT EXISTS idx_users_gender ON users(gender);
    `

	_, err := db.WriteDB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("Database tables initialized successfully")
	return nil
}

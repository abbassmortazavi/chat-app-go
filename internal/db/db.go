package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init(dbpath, dbname string) {
	err := os.MkdirAll(dbpath, os.ModePerm)
	if err != nil {
		log.Fatalf("Could not create database directory: %v", err)
	}
	dbfile := filepath.Join(dbpath, dbname)
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}
	//db.SetMaxIdleConns(30)
	//db.SetConnMaxIdleTime(5)
	//db.SetConnMaxLifetime(30)

	pragmas := []string{
		"PRAGMA foreign_keys = ON;",
		"PRAGMA journal_mode = WAL;",
		"PRAGMA busy_timeout = 5000;",
		"PRAGMA synchronous = NORMAL;",
	}
	for _, pragmas := range pragmas {
		_, err = db.Exec(pragmas)
		if err != nil {
			log.Fatalf("Could not create PRAGMA table: %v", err)
		}
	}

	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			password TEXT NOT NULL,
			refresh_token web TEXT,
			refresh_token_web_at DATETIME,
 			refresh_token_mobile TEXT,
 			refresh_token_mobile_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS privates ( 
    		id INTEGER PRIMARY KEY AUTOINCREMENT,
    		user1_id INTEGER NOT NULL,
    		user2_id INTEGER NOT NULL,
    		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    		UNIQUE (user1_id, user2_id),
			CHECK (user1_id < user2_id),
			FOREIGN KEY(user1_id) REFERENCES users(id) ON DELETE CASCADE,
    		FOREIGN KEY(user2_id) REFERENCES users(id) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS messages ( 
    		id INTEGER PRIMARY KEY AUTOINCREMENT,
    		from_id INTEGER NOT NULL,
    		private_id INTEGER,
			message_type TEXT NOT NULL,
			content TEXT NOT NULL,
			delivered INTEGER NOT NULL DEFAULT 0,
			read INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (from_id) REFERENCES users (id) ON DELETE CASCADE, 
    		FOREIGN KEY (private_id) REFERENCES privates(id) ON DELETE CASCADE
		);`,
	}
	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			log.Fatalf("Could not create table %s: %v", table, err)
		}
	}

	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_messages_private_id ON messages(private_id);`,
		`CREATE INDEX IF NOT EXISTS idx_messages_from_id ON messages(from_id);`,
		`CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);`,
		`CREATE INDEX IF NOT EXISTS idx_privates_user1_id ON privates(user1_id);`,
		`CREATE INDEX IF NOT EXISTS idx_privates_user2_id ON privates(user2_id);`,
	}
	for _, index := range indexes {
		_, err := db.Exec(index)
		if err != nil {
			log.Fatalf("Could not create index %s: %v", index, err)
		}
	}
	log.Printf("Database is ready")
	DB = db
}

func CloseDB() {
	if DB == nil {
		return
	}
	err := DB.Close()
	if err != nil {
		log.Println("Error closing DB", err)
	} else {
		log.Println("DB closed")
	}
}

//video 5====> 4min seen

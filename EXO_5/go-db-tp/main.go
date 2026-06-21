package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID    int    `db:"id" gorm:"primaryKey;autoIncrement"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

func createSchema(db *sql.DB) error {
	schema := `CREATE TABLE IF NOT EXISTS users (
        id    INTEGER PRIMARY KEY AUTOINCREMENT,
        name  TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL
    );`
	_, err := db.Exec(schema)
	return err
}

func CreateUserSQL(db *sql.DB, user User) (int64, error) {
	res, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetUsersSQL(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func GetUserByIDSQL(db *sql.DB, id int) (*User, error) {
	var u User
	err := db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func UpdateUserSQL(db *sql.DB, user User) error {
	_, err := db.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", user.Name, user.Email, user.ID)
	return err
}

func DeleteUserSQL(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

func CreateUserSQLX(db *sqlx.DB, user User) (int64, error) {
	res, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetUsersSQLX(db *sqlx.DB) ([]User, error) {
	var users []User
	err := db.Select(&users, "SELECT id, name, email FROM users")
	return users, err
}

func GetUserByIDSQLX(db *sqlx.DB, id int) (*User, error) {
	var u User
	err := db.Get(&u, "SELECT id, name, email FROM users WHERE id = ?", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func UpdateUserSQLX(db *sqlx.DB, user User) error {
	_, err := db.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", user.Name, user.Email, user.ID)
	return err
}

func DeleteUserSQLX(db *sqlx.DB, id int) error {
	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

func CreateUserGORM(db *gorm.DB, user User) error {
	return db.Create(&user).Error
}

func GetUsersGORM(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Find(&users).Error
	return users, err
}

func GetUserByIDGORM(db *gorm.DB, id int) (*User, error) {
	var u User
	err := db.First(&u, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func UpdateUserGORM(db *gorm.DB, user User) error {
	return db.Save(&user).Error
}

func DeleteUserGORM(db *gorm.DB, id int) error {
	return db.Delete(&User{}, id).Error
}

func runSQL() {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Connecté à la base de données SQLite.")

	if err := createSchema(db); err != nil {
		log.Fatal(err)
	}

	log.Println("=== database/sql ===")

	id, err := CreateUserSQL(db, User{Name: "Alice", Email: "alice@example.com"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Utilisateur créé, id = %d", id)

	if _, err := CreateUserSQL(db, User{Name: "Bob", Email: "bob@example.com"}); err != nil {
		log.Fatal(err)
	}

	users, err := GetUsersSQL(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Tous les utilisateurs : %+v", users)

	u, err := GetUserByIDSQL(db, int(id))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Utilisateur %d : %+v", id, u)

	if err := UpdateUserSQL(db, User{ID: int(id), Name: "Alice Martin", Email: "alice.martin@example.com"}); err != nil {
		log.Fatal(err)
	}
	u, _ = GetUserByIDSQL(db, int(id))
	log.Printf("Après mise à jour : %+v", u)

	if err := DeleteUserSQL(db, int(id)); err != nil {
		log.Fatal(err)
	}
	users, _ = GetUsersSQL(db)
	log.Printf("Après suppression : %+v", users)
}

func runSQLX() {
	db, err := sqlx.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("=== sqlx ===")

	id, err := CreateUserSQLX(db, User{Name: "Carol", Email: "carol@example.com"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Utilisateur créé, id = %d", id)

	users, err := GetUsersSQLX(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Tous les utilisateurs : %+v", users)

	u, err := GetUserByIDSQLX(db, int(id))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Utilisateur %d : %+v", id, u)

	if err := UpdateUserSQLX(db, User{ID: int(id), Name: "Carol Dupont", Email: "carol.dupont@example.com"}); err != nil {
		log.Fatal(err)
	}
	u, _ = GetUserByIDSQLX(db, int(id))
	log.Printf("Après mise à jour : %+v", u)

	if err := DeleteUserSQLX(db, int(id)); err != nil {
		log.Fatal(err)
	}
	users, _ = GetUsersSQLX(db)
	log.Printf("Après suppression : %+v", users)
}

func runGORM() {
	db, err := gorm.Open(sqlite.Open("gorm_test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Échec de la connexion à la base de données GORM:", err)
	}
	log.Println("Connecté à la base de données SQLite avec GORM.")

	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatal("Échec de l'auto-migration GORM:", err)
	}
	log.Println("Auto-migration GORM effectuée.")

	log.Println("=== GORM ===")

	user := User{Name: "David", Email: "david@example.com"}
	if err := CreateUserGORM(db, user); err != nil {
		log.Fatal(err)
	}

	users, err := GetUsersGORM(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Tous les utilisateurs : %+v", users)

	created := users[len(users)-1]

	u, err := GetUserByIDGORM(db, created.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Utilisateur %d : %+v", created.ID, u)

	if err := UpdateUserGORM(db, User{ID: created.ID, Name: "David Bernard", Email: "david.bernard@example.com"}); err != nil {
		log.Fatal(err)
	}
	u, _ = GetUserByIDGORM(db, created.ID)
	log.Printf("Après mise à jour : %+v", u)

	if err := DeleteUserGORM(db, created.ID); err != nil {
		log.Fatal(err)
	}
	users, _ = GetUsersGORM(db)
	log.Printf("Après suppression : %+v", users)
}

func main() {
	os.Remove("test.db")
	os.Remove("gorm_test.db")

	runSQL()
	runSQLX()
	runGORM()
}

package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Mutter0815/marketplace/internal/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type SQLRepository struct {
	db  *sql.DB
	log *zap.SugaredLogger
}

func NewSQLRepository(db *sql.DB, log *zap.SugaredLogger) *SQLRepository {
	return &SQLRepository{db: db, log: log}
}

func (r *SQLRepository) CreateUser(username string, password string) (*models.User, error) {
	r.log.Debugw("repo create user", "username", username)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		r.log.Errorw("password hash failed", "err", err)
		return nil, fmt.Errorf("hash password: %w", err)
	}
	var id int
	err = r.db.QueryRow("INSERT INTO users(username,password_hash) VALUES ($1,$2) RETURNING id", username, string(hash)).Scan(&id)
	if err != nil {
		r.log.Errorw("insert user failed", "err", err)
		return nil, err
	}
	r.log.Infow("user created", "id", id)
	return &models.User{ID: uint(id), Username: username}, nil
}

func (r *SQLRepository) AuthenticateUser(username string, password string) (*models.User, error) {
	r.log.Debugw("repo authenticate", "username", username)
	var user models.User
	err := r.db.QueryRow("SELECT id, password_hash,created_at FROM users WHERE username=$1", username).Scan(&user.ID, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		r.log.Warnw("user not found", "err", err)
		return nil, errors.New("Неверные данные")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		r.log.Warnw("password mismatch", "err", err)
		return nil, errors.New("Неверный данные")
	}
	user.Username = username
	return &user, nil
}

func (r *SQLRepository) CreateAd(userID uint, title, desc string, imageURL *string, price uint) (*models.Ad, error) {
	if price == 0 {
		r.log.Warn("create ad zero price")
		return nil, errors.New("price должен быть больше нуля")
	}

	var id int
	var createdAt time.Time
	var username string
	err := r.db.QueryRow(
		`INSERT INTO ads(user_id, title, description, image_url, price)
         VALUES($1, $2, $3, $4, $5)
       RETURNING id, created_at, (SELECT username FROM users WHERE id = $1)`,
		userID, title, desc, imageURL, price,
	).Scan(&id, &createdAt, &username)
	if err != nil {
		r.log.Errorw("insert ad failed", "err", err)
		return nil, err
	}

	ad := &models.Ad{
		ID:          uint(id),
		UserID:      userID,
		Title:       title,
		Description: desc,
		ImageURL:    imageURL,
		Price:       price,
		CreatedAt:   createdAt,
		Author:      username,
	}
	r.log.Infow("ad created", "id", ad.ID)
	return ad, nil
}

func (r *SQLRepository) ListAds(priceMin, priceMax *uint, sortBy, order string, page, pageSize int) ([]*models.Ad, int, error) {
	r.log.Debugw("repo list ads", "page", page, "page_size", pageSize)
	var parts []string
	var args []interface{}
	idx := 1

	if priceMin != nil {
		parts = append(parts, fmt.Sprintf("price >= $%d", idx))
		args = append(args, *priceMin)
		idx++
	}
	if priceMax != nil {
		parts = append(parts, fmt.Sprintf("price <= $%d", idx))
		args = append(args, *priceMax)
		idx++
	}
	where := ""
	if len(parts) > 0 {
		where = "WHERE " + strings.Join(parts, " AND ")
	}

	field := "created_at"
	if sortBy == "price" {
		field = "price"
	}

	offset := (page - 1) * pageSize
	args = append(args, pageSize, offset)
	limitIdx := idx
	offsetIdx := idx + 1

	query := fmt.Sprintf(
		`SELECT a.id, a.user_id, u.username, a.title, a.description, a.image_url, a.price, a.created_at
         FROM ads a
         JOIN users u ON a.user_id = u.id
         %s
         ORDER BY a.%s %s
         LIMIT $%d OFFSET $%d`,
		where, field, order, limitIdx, offsetIdx,
	)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		r.log.Errorw("query ads failed", "err", err)
		return nil, 0, fmt.Errorf("query ads: %w", err)
	}
	defer rows.Close()

	var ads []*models.Ad
	for rows.Next() {
		a := new(models.Ad)
		if err := rows.Scan(
			&a.ID, &a.UserID, &a.Author, &a.Title,
			&a.Description, &a.ImageURL,
			&a.Price, &a.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan ad: %w", err)
		}
		ads = append(ads, a)
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM ads %s", where)
	var total int
	if err := r.db.QueryRow(countQuery, args[:idx-1]...).Scan(&total); err != nil {
		r.log.Errorw("count ads failed", "err", err)
		return nil, 0, fmt.Errorf("count ads: %w", err)
	}
	r.log.Debugw("ads listed", "total", total)

	return ads, total, nil
}

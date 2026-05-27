package app

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"

	"goapi/backend/internal/model"
	"goapi/backend/internal/repository"
	"goapi/backend/internal/security"
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type Services struct {
	Auth  *AuthService
	Items *ItemService
}

func NewServices(users *repository.UserRepository, items *repository.ItemRepository, passwords *security.PasswordHasher, tokens *security.TokenService) *Services {
	return &Services{
		Auth:  NewAuthService(users, passwords, tokens),
		Items: NewItemService(items),
	}
}

type AuthService struct {
	users     *repository.UserRepository
	passwords *security.PasswordHasher
	tokens    *security.TokenService
}

type AuthResult struct {
	User      model.User `json:"user"`
	Token     string     `json:"token"`
	ExpiresAt time.Time  `json:"expires_at"`
}

func NewAuthService(users *repository.UserRepository, passwords *security.PasswordHasher, tokens *security.TokenService, rdb *redis.Client) *AuthService {
	return &AuthService{users: users, passwords: passwords, tokens: tokens, redis: rdb}
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (AuthResult, error) {
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))

	if name == "" || len(password) < 8 || !validEmail(email) {
		return AuthResult{}, ErrInvalidInput
	}

	hash, err := s.passwords.Hash(password)
	if err != nil {
		return AuthResult{}, err
	}

	user, err := s.users.Create(ctx, name, email, hash)
	if errors.Is(err, repository.ErrDuplicateEmail) {
		return AuthResult{}, ErrEmailAlreadyExists
	}
	if err != nil {
		return AuthResult{}, err
	}

	token, expiresAt, err := s.tokens.Generate(user.ID)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{User: user, Token: token, ExpiresAt: expiresAt}, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (AuthResult, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if !validEmail(email) || password == "" {
		return AuthResult{}, ErrInvalidCredentials
	}

	user, err := s.users.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrNotFound) {
		return AuthResult{}, ErrInvalidCredentials
	}
	if err != nil {
		return AuthResult{}, err
	}
	if !s.passwords.Compare(user.PasswordHash, password) {
		return AuthResult{}, ErrInvalidCredentials
	}

	token, expiresAt, err := s.tokens.Generate(user.ID)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{User: user, Token: token, ExpiresAt: expiresAt}, nil
}

func (s *AuthService) Logout(ctx context.Context, tokenID string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return nil
	}
	return s.redis.Set(ctx, "blacklist:"+tokenID, "true", ttl).Err()
}

func (s *AuthService) IsTokenBlacklisted(ctx context.Context, tokenID string) bool {
	exists, err := s.redis.Exists(ctx, "blacklist:"+tokenID).Result()
	return err == nil && exists > 0
}

func (s *AuthService) AllowRequest(ctx context.Context, ip string, maxReqs int64, window time.Duration) bool {
	key := "ratelimit:auth:" + ip
	
	count, err := s.redis.Incr(ctx, key).Result()
	if err != nil {
		return false // On Redis failure, default to block or handle gracefully. Let's block to be safe.
	}
	
	if count == 1 {
		s.redis.Expire(ctx, key, window)
	}

	return count <= maxReqs
}

type ItemService struct {
	items *repository.ItemRepository
}

func NewItemService(items *repository.ItemRepository) *ItemService {
	return &ItemService{items: items}
}

func (s *ItemService) Create(ctx context.Context, userID, title, body string) (model.Item, error) {
	title = strings.TrimSpace(title)
	body = strings.TrimSpace(body)

	if userID == "" || title == "" || len(title) > 160 || len(body) > 5000 {
		return model.Item{}, ErrInvalidInput
	}

	return s.items.Create(ctx, userID, title, body)
}

func (s *ItemService) List(ctx context.Context, userID string, limit int32) ([]model.Item, error) {
	if userID == "" {
		return nil, ErrInvalidInput
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	return s.items.ListByUser(ctx, userID, limit)
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
turn err == nil
}

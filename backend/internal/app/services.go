package app

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"goapi/backend/internal/mailer"
	"goapi/backend/internal/model"
	"goapi/backend/internal/repository"
	"goapi/backend/internal/security"
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUnauthorized       = errors.New("unauthorized")
)

type Services struct {
	Auth   *AuthService
	Items  *ItemService
	Stores *StoreService
}

func NewServices(
	users *repository.UserRepository,
	items *repository.ItemRepository,
	stores *repository.StoreRepository,
	passwords *security.PasswordHasher,
	tokens *security.TokenService,
	rdb *redis.Client,
	mail *mailer.Mailer,
) *Services {
	return &Services{
		Auth:   NewAuthService(users, passwords, tokens, rdb, mail),
		Items:  NewItemService(items),
		Stores: NewStoreService(stores),
	}
}

type AuthService struct {
	users     *repository.UserRepository
	passwords *security.PasswordHasher
	tokens    *security.TokenService
	redis     *redis.Client
	mail      *mailer.Mailer
}

type AuthResult struct {
	User      model.User `json:"user"`
	Token     string     `json:"token"`
	ExpiresAt time.Time  `json:"expires_at"`
}

func NewAuthService(users *repository.UserRepository, passwords *security.PasswordHasher, tokens *security.TokenService, rdb *redis.Client, mail *mailer.Mailer) *AuthService {
	return &AuthService{users: users, passwords: passwords, tokens: tokens, redis: rdb, mail: mail}
}

func (s *AuthService) Register(ctx context.Context, name, email, password, role string) (AuthResult, error) {
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))

	if role == "" {
		role = model.RoleClient
	}

	if name == "" || len(password) < 8 || !validEmail(email) {
		return AuthResult{}, ErrInvalidInput
	}

	hash, err := s.passwords.Hash(password)
	if err != nil {
		return AuthResult{}, err
	}

	user, err := s.users.Create(ctx, name, email, hash, role)
	if errors.Is(err, repository.ErrDuplicateEmail) {
		return AuthResult{}, ErrEmailAlreadyExists
	}
	if err != nil {
		return AuthResult{}, err
	}

	token, expiresAt, err := s.tokens.Generate(user.ID, user.Role, user.Email, user.Name)
	if err != nil {
		return AuthResult{}, err
	}

	// Dispara o e-mail de boas-vindas assincronamente (em background)
	go func() {
		_ = s.mail.SendWelcomeEmail(user.Email, user.Name)
	}()

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

	token, expiresAt, err := s.tokens.Generate(user.ID, user.Role, user.Email, user.Name)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{User: user, Token: token, ExpiresAt: expiresAt}, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id string) (model.User, error) {
	// Implementar se necessário buscar por ID, por enquanto usamos claims
	return model.User{}, nil
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
		return false
	}
	
	if count == 1 {
		s.redis.Expire(ctx, key, window)
	}

	return count <= maxReqs
}

type StoreService struct {
	stores *repository.StoreRepository
}

func NewStoreService(stores *repository.StoreRepository) *StoreService {
	return &StoreService{stores: stores}
}

func (s *StoreService) Create(ctx context.Context, ownerID, name, slug string) (model.Store, error) {
	name = strings.TrimSpace(name)
	slug = strings.ToLower(strings.TrimSpace(slug))

	if ownerID == "" || name == "" || slug == "" {
		return model.Store{}, ErrInvalidInput
	}

	return s.stores.Create(ctx, ownerID, name, slug)
}

func (s *StoreService) GetBySlug(ctx context.Context, slug string) (model.Store, error) {
	return s.stores.FindBySlug(ctx, slug)
}

func (s *StoreService) GetByOwner(ctx context.Context, ownerID string) (model.Store, error) {
	return s.stores.FindByOwner(ctx, ownerID)
}

type ItemService struct {
	items *repository.ItemRepository
}

func NewItemService(items *repository.ItemRepository) *ItemService {
	return &ItemService{items: items}
}

func (s *ItemService) Create(ctx context.Context, item model.Item) (model.Item, error) {
	item.UserID = strings.TrimSpace(item.UserID)
	item.StoreID = strings.TrimSpace(item.StoreID)
	item.Name = strings.TrimSpace(item.Name)
	item.Description = strings.TrimSpace(item.Description)
	item.Variant = strings.TrimSpace(item.Variant)
	item.ShippingType = strings.TrimSpace(item.ShippingType)

	if item.UserID == "" || item.StoreID == "" || item.Name == "" || len(item.Name) > 160 || len(item.Description) > 5000 {
		return model.Item{}, ErrInvalidInput
	}
	if item.Price < 0 || item.VariantPrice < 0 {
		return model.Item{}, ErrInvalidInput
	}
	if item.ShippingType != "free" && item.ShippingType != "consult" {
		return model.Item{}, ErrInvalidInput
	}

	return s.items.Create(ctx, item)
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

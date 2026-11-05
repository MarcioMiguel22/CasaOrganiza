package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sysadminsmedia/homebox/backend/internal/data/ent"
)

var ErrPantryProductNotFound = errors.New("pantry product not found")

type PantryRepository struct {
	db *ent.Client
}

type PantryProduct struct {
	ID              uuid.UUID `json:"id"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Name            string    `json:"name"`
	Quantity        float64   `json:"quantity"`
	MinQuantity     float64   `json:"minQuantity"`
	Unit            string    `json:"unit"`
	ExpiryDate      string    `json:"expiryDate"`
	PurchasePrice   float64   `json:"purchasePrice"`
	OnShoppingList  bool      `json:"onShoppingList"`
}

type PantryProductCreate struct {
	Name          string  `json:"name" validate:"required,min=1,max=255"`
	Quantity      float64 `json:"quantity"`
	MinQuantity   float64 `json:"minQuantity"`
	Unit          string  `json:"unit" validate:"max=40"`
	ExpiryDate    string  `json:"expiryDate"`
	PurchasePrice float64 `json:"purchasePrice"`
}

type PantryProductUpdate struct {
	Name            string  `json:"name" validate:"required,min=1,max=255"`
	Quantity        float64 `json:"quantity"`
	MinQuantity     float64 `json:"minQuantity"`
	Unit            string  `json:"unit" validate:"max=40"`
	ExpiryDate      string  `json:"expiryDate"`
	PurchasePrice   float64 `json:"purchasePrice"`
	OnShoppingList  bool    `json:"onShoppingList"`
}

type PantryConsume struct {
	Amount float64 `json:"amount" validate:"gt=0"`
}

type PantryOverview struct {
	Products      []PantryProduct `json:"products"`
	LowStock      []PantryProduct `json:"lowStock"`
	Expiring      []PantryProduct `json:"expiring"`
	ShoppingList  []PantryProduct `json:"shoppingList"`
}

func (p PantryProductCreate) Validate() error {
	return validatePantryValues(p.Quantity, p.MinQuantity, p.ExpiryDate, p.PurchasePrice)
}

func (p PantryProductUpdate) Validate() error {
	return validatePantryValues(p.Quantity, p.MinQuantity, p.ExpiryDate, p.PurchasePrice)
}

func validatePantryValues(quantity, minQuantity float64, expiryDate string, purchasePrice float64) error {
	if math.IsNaN(quantity) || math.IsInf(quantity, 0) || quantity < 0 {
		return errors.New("quantity must be a finite number and cannot be negative")
	}
	if math.IsNaN(minQuantity) || math.IsInf(minQuantity, 0) || minQuantity < 0 {
		return errors.New("minimum quantity must be a finite number and cannot be negative")
	}
	if math.IsNaN(purchasePrice) || math.IsInf(purchasePrice, 0) || purchasePrice < 0 {
		return errors.New("purchase price must be a finite number and cannot be negative")
	}
	if expiryDate != "" {
		if _, err := time.Parse("2006-01-02", expiryDate); err != nil {
			return errors.New("expiry date must use YYYY-MM-DD")
		}
	}
	return nil
}

func (r *PantryRepository) List(ctx context.Context, gid uuid.UUID) ([]PantryProduct, error) {
	rows, err := r.db.Sql().QueryContext(ctx, `
		SELECT id, created_at, updated_at, name, quantity, min_quantity, unit,
		       expiry_date, purchase_price, on_shopping_list
		FROM pantry_products
		WHERE group_pantry_products = $1
		ORDER BY name ASC`, gid)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	products := make([]PantryProduct, 0)
	for rows.Next() {
		var product PantryProduct
		if err := rows.Scan(
			&product.ID,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.Name,
			&product.Quantity,
			&product.MinQuantity,
			&product.Unit,
			&product.ExpiryDate,
			&product.PurchasePrice,
			&product.OnShoppingList,
		); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *PantryRepository) Get(ctx context.Context, gid, id uuid.UUID) (PantryProduct, error) {
	var product PantryProduct
	err := r.db.Sql().QueryRowContext(ctx, `
		SELECT id, created_at, updated_at, name, quantity, min_quantity, unit,
		       expiry_date, purchase_price, on_shopping_list
		FROM pantry_products
		WHERE id = $1 AND group_pantry_products = $2`, id, gid).Scan(
		&product.ID,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.Name,
		&product.Quantity,
		&product.MinQuantity,
		&product.Unit,
		&product.ExpiryDate,
		&product.PurchasePrice,
		&product.OnShoppingList,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return PantryProduct{}, ErrPantryProductNotFound
	}
	if err != nil {
		return PantryProduct{}, err
	}
	return product, nil
}

func (r *PantryRepository) Create(ctx context.Context, gid uuid.UUID, data PantryProductCreate) (PantryProduct, error) {
	id := uuid.New()
	now := time.Now().UTC()
	_, err := r.db.Sql().ExecContext(ctx, `
		INSERT INTO pantry_products (
			id, created_at, updated_at, name, quantity, min_quantity, unit,
			expiry_date, purchase_price, on_shopping_list, group_pantry_products
		) VALUES ($1, $2, $2, $3, $4, $5, $6, $7, $8, false, $9)`,
		id,
		now,
		strings.TrimSpace(data.Name),
		data.Quantity,
		data.MinQuantity,
		strings.TrimSpace(data.Unit),
		strings.TrimSpace(data.ExpiryDate),
		data.PurchasePrice,
		gid,
	)
	if err != nil {
		return PantryProduct{}, err
	}
	return r.Get(ctx, gid, id)
}

func (r *PantryRepository) Update(ctx context.Context, gid, id uuid.UUID, data PantryProductUpdate) (PantryProduct, error) {
	result, err := r.db.Sql().ExecContext(ctx, `
		UPDATE pantry_products
		SET updated_at = $1, name = $2, quantity = $3, min_quantity = $4,
		    unit = $5, expiry_date = $6, purchase_price = $7, on_shopping_list = $8
		WHERE id = $9 AND group_pantry_products = $10`,
		time.Now().UTC(),
		strings.TrimSpace(data.Name),
		data.Quantity,
		data.MinQuantity,
		strings.TrimSpace(data.Unit),
		strings.TrimSpace(data.ExpiryDate),
		data.PurchasePrice,
		data.OnShoppingList,
		id,
		gid,
	)
	if err != nil {
		return PantryProduct{}, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return PantryProduct{}, err
	}
	if count == 0 {
		return PantryProduct{}, ErrPantryProductNotFound
	}
	return r.Get(ctx, gid, id)
}

func (r *PantryRepository) Consume(ctx context.Context, gid, id uuid.UUID, amount float64) (PantryProduct, error) {
	result, err := r.db.Sql().ExecContext(ctx, `
		UPDATE pantry_products
		SET updated_at = $1,
		    quantity = CASE WHEN quantity > $2 THEN quantity - $2 ELSE 0 END
		WHERE id = $3 AND group_pantry_products = $4`,
		time.Now().UTC(), amount, id, gid)
	if err != nil {
		return PantryProduct{}, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return PantryProduct{}, err
	}
	if count == 0 {
		return PantryProduct{}, ErrPantryProductNotFound
	}
	return r.Get(ctx, gid, id)
}

func (r *PantryRepository) ToggleShopping(ctx context.Context, gid, id uuid.UUID) (PantryProduct, error) {
	result, err := r.db.Sql().ExecContext(ctx, `
		UPDATE pantry_products
		SET updated_at = $1,
		    on_shopping_list = CASE WHEN on_shopping_list THEN false ELSE true END
		WHERE id = $2 AND group_pantry_products = $3`,
		time.Now().UTC(), id, gid)
	if err != nil {
		return PantryProduct{}, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return PantryProduct{}, err
	}
	if count == 0 {
		return PantryProduct{}, ErrPantryProductNotFound
	}
	return r.Get(ctx, gid, id)
}

func (r *PantryRepository) Overview(ctx context.Context, gid uuid.UUID) (PantryOverview, error) {
	products, err := r.List(ctx, gid)
	if err != nil {
		return PantryOverview{}, fmt.Errorf("list pantry products: %w", err)
	}

	soon := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
	overview := PantryOverview{
		Products:     products,
		LowStock:     make([]PantryProduct, 0),
		Expiring:     make([]PantryProduct, 0),
		ShoppingList: make([]PantryProduct, 0),
	}
	for _, product := range products {
		if product.MinQuantity > 0 && product.Quantity <= product.MinQuantity {
			overview.LowStock = append(overview.LowStock, product)
		}
		if product.ExpiryDate != "" && product.ExpiryDate <= soon {
			overview.Expiring = append(overview.Expiring, product)
		}
		if product.OnShoppingList {
			overview.ShoppingList = append(overview.ShoppingList, product)
		}
	}
	return overview, nil
}

package v1

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/hay-kot/httpkit/errchain"
	"github.com/sysadminsmedia/homebox/backend/internal/core/services"
	"github.com/sysadminsmedia/homebox/backend/internal/data/repo"
	"github.com/sysadminsmedia/homebox/backend/internal/sys/validate"
	"github.com/sysadminsmedia/homebox/backend/internal/web/adapters"
)

func pantryError(err error) error {
	if errors.Is(err, repo.ErrPantryProductNotFound) {
		return validate.NewRequestError(err, http.StatusNotFound)
	}
	return err
}

// HandlePantryOverview godoc
//
//	@Summary	Get the internal pantry overview
//	@Tags		Pantry
//	@Produce	json
//	@Success	200	{object}	repo.PantryOverview
//	@Router		/v1/pantry [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandlePantryOverview() errchain.HandlerFunc {
	fn := func(r *http.Request) (repo.PantryOverview, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Pantry.Overview(r.Context(), auth.GID)
	}

	return adapters.Command(fn, http.StatusOK)
}

// HandlePantryProductsGetAll godoc
//
//	@Summary	List pantry products
//	@Tags		Pantry
//	@Produce	json
//	@Success	200	{array}	repo.PantryProduct
//	@Router		/v1/pantry/products [GET]
//	@Security	Bearer
func (ctrl *V1Controller) HandlePantryProductsGetAll() errchain.HandlerFunc {
	fn := func(r *http.Request) ([]repo.PantryProduct, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Pantry.List(r.Context(), auth.GID)
	}

	return adapters.Command(fn, http.StatusOK)
}

// HandlePantryProductCreate godoc
//
//	@Summary	Create a pantry product
//	@Tags		Pantry
//	@Produce	json
//	@Param		payload	body	repo.PantryProductCreate	true	"Pantry product"
//	@Success	201	{object}	repo.PantryProduct
//	@Router		/v1/pantry/products [POST]
//	@Security	Bearer
func (ctrl *V1Controller) HandlePantryProductCreate() errchain.HandlerFunc {
	fn := func(r *http.Request, data repo.PantryProductCreate) (repo.PantryProduct, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.Pantry.Create(r.Context(), auth.GID, data)
	}

	return adapters.Action(fn, http.StatusCreated)
}

// HandlePantryProductUpdate godoc
//
//	@Summary	Update a pantry product
//	@Tags		Pantry
//	@Produce	json
//	@Param		id		path	string	true	"Pantry product ID"
//	@Param		payload	body	repo.PantryProductUpdate	true	"Pantry product"
//	@Success	200	{object}	repo.PantryProduct
//	@Router		/v1/pantry/products/{id} [PUT]
//	@Security	Bearer
func (ctrl *V1Controller) HandlePantryProductUpdate() errchain.HandlerFunc {
	fn := func(r *http.Request, id uuid.UUID, data repo.PantryProductUpdate) (repo.PantryProduct, error) {
		auth := services.NewContext(r.Context())
		product, err := ctrl.repo.Pantry.Update(r.Context(), auth.GID, id, data)
		return product, pantryError(err)
	}

	return adapters.ActionID("id", fn, http.StatusOK)
}

// HandlePantryProductConsume godoc
//
//	@Summary	Consume pantry stock
//	@Tags		Pantry
//	@Produce	json
//	@Param		id		path	string	true	"Pantry product ID"
//	@Param		payload	body	repo.PantryConsume	true	"Amount to consume"
//	@Success	200	{object}	repo.PantryProduct
//	@Router		/v1/pantry/products/{id}/consume [POST]
//	@Security	Bearer
func (ctrl *V1Controller) HandlePantryProductConsume() errchain.HandlerFunc {
	fn := func(r *http.Request, id uuid.UUID, data repo.PantryConsume) (repo.PantryProduct, error) {
		auth := services.NewContext(r.Context())
		product, err := ctrl.repo.Pantry.Consume(r.Context(), auth.GID, id, data.Amount)
		return product, pantryError(err)
	}

	return adapters.ActionID("id", fn, http.StatusOK)
}

// HandlePantryProductShoppingToggle godoc
//
//	@Summary	Toggle a product on the shopping list
//	@Tags		Pantry
//	@Produce	json
//	@Param		id	path	string	true	"Pantry product ID"
//	@Success	200	{object}	repo.PantryProduct
//	@Router		/v1/pantry/products/{id}/shopping [POST]
//	@Security	Bearer
func (ctrl *V1Controller) HandlePantryProductShoppingToggle() errchain.HandlerFunc {
	fn := func(r *http.Request, id uuid.UUID) (repo.PantryProduct, error) {
		auth := services.NewContext(r.Context())
		product, err := ctrl.repo.Pantry.ToggleShopping(r.Context(), auth.GID, id)
		return product, pantryError(err)
	}

	return adapters.CommandID("id", fn, http.StatusOK)
}

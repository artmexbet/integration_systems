package v1

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for the Nobel Prize API
type Handler struct {
	service   Service
	validator *validator.Validate
}

// NewHandler creates a new Handler instance
func NewHandler(service Service) *Handler {
	return &Handler{
		service:   service,
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}

// GetStats godoc
//
//	@Summary		Get dataset statistics
//	@Description	Returns count of laureates, prizes, and categories
//	@Tags			Stats
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		ApiKeyAuth
//	@Success		200	{object}	StatsResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/api/v1/stats [get]
//	@security		ApiKeyAuth
func (h *Handler) GetStats(c *fiber.Ctx) error {
	stats, err := h.service.GetStats(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}
	return c.JSON(stats)
}

// GetLastUpdate godoc
//
//	@Summary		Get last update timestamp
//	@Description	Returns the timestamp of the last dataset update
//	@Tags			Stats
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		ApiKeyAuth
//	@Success		200	{object}	LastUpdateResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/api/v1/stats/last-update [get]
//	@security		ApiKeyAuth
func (h *Handler) GetLastUpdate(c *fiber.Ctx) error {
	lastUpdate, err := h.service.GetLastUpdate(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}
	return c.JSON(lastUpdate)
}

// ListLaureates godoc
//
//	@Summary		List laureates
//	@Description	Returns a paginated list of Nobel laureates
//	@Tags			Laureates
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		ApiKeyAuth
//	@Param			page		query		int	false	"Page number"		default(1)
//	@Param			per_page	query		int	false	"Items per page"	default(10)	maximum(100)
//	@Success		200			{object}	LaureateListResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/v1/laureates [get]
//	@security		ApiKeyAuth
func (h *Handler) ListLaureates(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	result, err := h.service.ListLaureates(c.Context(), page, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}
	return c.JSON(result)
}

// GetLaureate godoc
//
//	@Summary		Get laureate by ID
//	@Description	Returns a single laureate by their ID
//	@Tags			Laureates
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Laureate ID"
//	@Success		200	{object}	LaureateResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/api/v1/laureates/{id} [get]
//	@security		ApiKeyAuth
func (h *Handler) GetLaureate(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid laureate ID",
		})
	}

	laureate, err := h.service.GetLaureate(c.Context(), int32(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:   "Not Found",
			Message: err.Error(),
		})
	}
	return c.JSON(laureate)
}

// CreateLaureate godoc
//
//	@Summary		Create a new laureate
//	@Description	Creates a new Nobel laureate
//	@Tags			Laureates
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		ApiKeyAuth
//	@Param			laureate	body		CreateLaureateRequest	true	"Laureate data"
//	@Success		201			{object}	LaureateResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/v1/laureates [post]
//	@security		ApiKeyAuth
func (h *Handler) CreateLaureate(c *fiber.Ctx) error {
	var req CreateLaureateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Validation Error",
			Message: err.Error(),
		})
	}

	laureate, err := h.service.CreateLaureate(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(laureate)
}

// UpdateLaureate godoc
//
//	@Summary		Update a laureate
//	@Description	Updates an existing Nobel laureate
//	@Tags			Laureates
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		ApiKeyAuth
//	@Param			id			path		int						true	"Laureate ID"
//	@Param			laureate	body		UpdateLaureateRequest	true	"Laureate data"
//	@Success		200			{object}	LaureateResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/v1/laureates/{id} [put]
//	@security		ApiKeyAuth
func (h *Handler) UpdateLaureate(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid laureate ID",
		})
	}

	var req UpdateLaureateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Validation Error",
			Message: err.Error(),
		})
	}

	laureate, err := h.service.UpdateLaureate(c.Context(), int32(id), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}
	return c.JSON(laureate)
}

// DeleteLaureate godoc
//
//	@Summary		Delete a laureate
//	@Description	Deletes an existing Nobel laureate
//	@Tags			Laureates
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Laureate ID"
//	@Success		200	{object}	SuccessResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/api/v1/laureates/{id} [delete]
//	@security		ApiKeyAuth
func (h *Handler) DeleteLaureate(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid laureate ID",
		})
	}

	err = h.service.DeleteLaureate(c.Context(), int32(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}
	return c.JSON(SuccessResponse{Message: "Laureate deleted successfully"})
}

// ListPrizes godoc
//
//	@Summary		List prizes
//	@Description	Returns a paginated list of Nobel prizes
//	@Tags			Prizes
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		ApiKeyAuth
//	@Param			page		query		int	false	"Page number"		default(1)
//	@Param			per_page	query		int	false	"Items per page"	default(10)	maximum(100)
//	@Success		200			{object}	PrizeListResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/v1/prizes [get]
//	@security		ApiKeyAuth
func (h *Handler) ListPrizes(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	result, err := h.service.ListPrizes(c.Context(), page, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}
	return c.JSON(result)
}

// GetPrize godoc
//
//	@Summary		Get prize by ID
//	@Description	Returns a single prize by its ID with associated laureates
//	@Tags			Prizes
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Prize ID"
//	@Success		200	{object}	PrizeResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/api/v1/prizes/{id} [get]
//	@security		ApiKeyAuth
func (h *Handler) GetPrize(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid prize ID",
		})
	}

	prize, err := h.service.GetPrize(c.Context(), int32(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:   "Not Found",
			Message: err.Error(),
		})
	}
	return c.JSON(prize)
}

// GetPrizesByCategory godoc
//
//	@Summary		Get prizes by category
//	@Description	Returns all prizes for a specific category with their laureates
//	@Tags			Prizes
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		ApiKeyAuth
//	@Param			category	path		string	true	"Prize category (e.g., physics, chemistry, medicine, literature, peace, economics)"
//	@Success		200			{array}		PrizeResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/v1/prizes/category/{category} [get]
//	@security		ApiKeyAuth
func (h *Handler) GetPrizesByCategory(c *fiber.Ctx) error {
	category := c.Params("category")
	if category == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Bad Request",
			Message: "Category is required",
		})
	}

	prizes, err := h.service.GetPrizesByCategory(c.Context(), category)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}
	return c.JSON(prizes)
}

// GetPrizesByYear godoc
//
//	@Summary		Get prizes by year
//	@Description	Returns all prizes for a specific year
//	@Tags			Prizes
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		ApiKeyAuth
//	@Param			year	path		int	true	"Prize year"
//	@Success		200		{array}		PrizeResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/api/v1/prizes/year/{year} [get]
//	@security		ApiKeyAuth
func (h *Handler) GetPrizesByYear(c *fiber.Ctx) error {
	year, err := strconv.ParseInt(c.Params("year"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid year",
		})
	}

	prizes, err := h.service.GetPrizesByYear(c.Context(), int32(year))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}
	return c.JSON(prizes)
}

// CreatePrize godoc
//
//	@Summary		Create a new prize
//	@Description	Creates a new Nobel prize
//	@Tags			Prizes
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		ApiKeyAuth
//	@Param			prize	body		CreatePrizeRequest	true	"Prize data"
//	@Success		201		{object}	PrizeResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/api/v1/prizes [post]
//	@security		ApiKeyAuth
func (h *Handler) CreatePrize(c *fiber.Ctx) error {
	var req CreatePrizeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Validation Error",
			Message: err.Error(),
		})
	}

	prize, err := h.service.CreatePrize(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(prize)
}

// UpdatePrize godoc
//
//	@Summary		Update a prize
//	@Description	Updates an existing Nobel prize
//	@Tags			Prizes
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		ApiKeyAuth
//	@Param			id		path		int					true	"Prize ID"
//	@Param			prize	body		UpdatePrizeRequest	true	"Prize data"
//	@Success		200		{object}	PrizeResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/api/v1/prizes/{id} [put]
//	@security		ApiKeyAuth
func (h *Handler) UpdatePrize(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid prize ID",
		})
	}

	var req UpdatePrizeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Validation Error",
			Message: err.Error(),
		})
	}

	prize, err := h.service.UpdatePrize(c.Context(), int32(id), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}
	return c.JSON(prize)
}

// DeletePrize godoc
//
//	@Summary		Delete a prize
//	@Description	Deletes an existing Nobel prize
//	@Tags			Prizes
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Prize ID"
//	@Success		200	{object}	SuccessResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/api/v1/prizes/{id} [delete]
//	@security		ApiKeyAuth
func (h *Handler) DeletePrize(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid prize ID",
		})
	}

	err = h.service.DeletePrize(c.Context(), int32(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}
	return c.JSON(SuccessResponse{Message: "Prize deleted successfully"})
}

// GetCategories godoc
//
//	@Summary		Get all categories
//	@Description	Returns a list of all unique prize categories
//	@Tags			Prizes
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		ApiKeyAuth
//	@Success		200	{object}	CategoriesResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/api/v1/categories [get]
//	@security		ApiKeyAuth
func (h *Handler) GetCategories(c *fiber.Ctx) error {
	categories, err := h.service.GetCategories(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}
	return c.JSON(categories)
}

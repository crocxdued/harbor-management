package handlers

import (
	"harbor/middleware"
	"harbor/models"
	"harbor/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct{ svc *service.Service }

func New(svc *service.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(r *gin.Engine) {
	api := r.Group("/api")

	api.POST("/auth/login", h.login)

	g := api.Group("")
	g.Use(middleware.Auth(h.svc))

	// ── Users (обязательные эндпоинты по заданию) ──
	g.GET("/users",    h.getUsers)
	g.GET("/users/:id", h.getUserByID)
	g.POST("/users",   middleware.RequireRoles("admin"), h.createUser)
	g.PUT("/users/:id", middleware.RequireRoles("admin"), h.updateUser)
	g.DELETE("/users/:id", middleware.RequireRoles("admin"), h.deleteUser)

	// ── Ships ──
	g.GET("/ships",    h.getShips)
	g.GET("/ships/:id", h.getShipByID)
	g.POST("/ships",   middleware.RequireRoles("admin", "dispatcher"), h.createShip)
	g.PUT("/ships/:id", middleware.RequireRoles("admin", "dispatcher"), h.updateShip)
	g.DELETE("/ships/:id", middleware.RequireRoles("admin"), h.deleteShip)

	// ── Visits ──
	g.GET("/visits",    h.getVisits)
	g.GET("/visits/:id", h.getVisitByID)
	g.POST("/visits",   middleware.RequireRoles("admin", "dispatcher"), h.createVisit)
	g.PUT("/visits/:id", middleware.RequireRoles("admin", "dispatcher"), h.updateVisit)
	g.DELETE("/visits/:id", middleware.RequireRoles("admin"), h.deleteVisit)

	g.GET("/me", h.me)
}

// ── helpers ──────────────────────────────────────────────────

func id(c *gin.Context) (int, bool) {
	n, err := strconv.Atoi(c.Param("id"))
	if err != nil { c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "неверный id"}); return 0, false }
	return n, true
}

func fail(c *gin.Context, status int, err error) {
	c.JSON(status, models.ErrorResponse{Error: err.Error()})
}

// ── auth ──────────────────────────────────────────────────────

func (h *Handler) login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil { fail(c, 400, err); return }
	resp, err := h.svc.Login(req)
	if err != nil { fail(c, 401, err); return }
	c.JSON(200, resp)
}

func (h *Handler) me(c *gin.Context) {
	uid, _ := c.Get("user_id")
	u, err := h.svc.GetUserByID(uid.(int))
	if err != nil { fail(c, 404, err); return }
	c.JSON(200, u)
}

// ── users ─────────────────────────────────────────────────────

func (h *Handler) getUsers(c *gin.Context) {
	list, err := h.svc.GetAllUsers()
	if err != nil { fail(c, 500, err); return }
	c.JSON(200, list)
}

func (h *Handler) getUserByID(c *gin.Context) {
	n, ok := id(c); if !ok { return }
	u, err := h.svc.GetUserByID(n)
	if err != nil { fail(c, 404, err); return }
	c.JSON(200, u)
}

func (h *Handler) createUser(c *gin.Context) {
	var req models.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil { fail(c, 400, err); return }
	u, err := h.svc.CreateUser(req)
	if err != nil { fail(c, 409, err); return }
	c.JSON(201, u)
}

func (h *Handler) updateUser(c *gin.Context) {
	n, ok := id(c); if !ok { return }
	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil { fail(c, 400, err); return }
	u, err := h.svc.UpdateUser(n, req)
	if err != nil { fail(c, 404, err); return }
	c.JSON(200, u)
}

func (h *Handler) deleteUser(c *gin.Context) {
	n, ok := id(c); if !ok { return }
	if err := h.svc.DeleteUser(n); err != nil { fail(c, 404, err); return }
	c.JSON(200, models.MsgResponse{Message: "пользователь удалён"})
}

// ── ships ─────────────────────────────────────────────────────

func (h *Handler) getShips(c *gin.Context) {
	list, err := h.svc.GetAllShips()
	if err != nil { fail(c, 500, err); return }
	c.JSON(200, list)
}

func (h *Handler) getShipByID(c *gin.Context) {
	n, ok := id(c); if !ok { return }
	s, err := h.svc.GetShipByID(n)
	if err != nil { fail(c, 404, err); return }
	c.JSON(200, s)
}

func (h *Handler) createShip(c *gin.Context) {
	var req models.ShipCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil { fail(c, 400, err); return }
	s, err := h.svc.CreateShip(req)
	if err != nil { fail(c, 409, err); return }
	c.JSON(201, s)
}

func (h *Handler) updateShip(c *gin.Context) {
	n, ok := id(c); if !ok { return }
	var req models.ShipUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil { fail(c, 400, err); return }
	s, err := h.svc.UpdateShip(n, req)
	if err != nil { fail(c, 404, err); return }
	c.JSON(200, s)
}

func (h *Handler) deleteShip(c *gin.Context) {
	n, ok := id(c); if !ok { return }
	if err := h.svc.DeleteShip(n); err != nil { fail(c, 404, err); return }
	c.JSON(200, models.MsgResponse{Message: "судно удалено"})
}

// ── visits ────────────────────────────────────────────────────

func (h *Handler) getVisits(c *gin.Context) {
	list, err := h.svc.GetAllVisits()
	if err != nil { fail(c, 500, err); return }
	c.JSON(200, list)
}

func (h *Handler) getVisitByID(c *gin.Context) {
	n, ok := id(c); if !ok { return }
	v, err := h.svc.GetVisitByID(n)
	if err != nil { fail(c, 404, err); return }
	c.JSON(200, v)
}

func (h *Handler) createVisit(c *gin.Context) {
	var req models.VisitCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil { fail(c, 400, err); return }
	v, err := h.svc.CreateVisit(req)
	if err != nil { fail(c, 409, err); return }
	c.JSON(201, v)
}

func (h *Handler) updateVisit(c *gin.Context) {
	n, ok := id(c); if !ok { return }
	var req models.VisitUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil { fail(c, 400, err); return }
	v, err := h.svc.UpdateVisit(n, req)
	if err != nil { fail(c, 404, err); return }
	c.JSON(200, v)
}

func (h *Handler) deleteVisit(c *gin.Context) {
	n, ok := id(c); if !ok { return }
	if err := h.svc.DeleteVisit(n); err != nil { fail(c, 404, err); return }
	c.JSON(200, models.MsgResponse{Message: "визит удалён"})
}

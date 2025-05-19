package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
	"homelab.com/homelab-server/homeLab-server/internal/usecase"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body entities.UserCredentials true "Login credentials"
// @Success 200 {object} entities.AuthResponse "JWT token"
// @Failure 400 {object} object{error=string} "Invalid request format"
// @Failure 401 {object} object{error=string} "Invalid credentials"
// @Router /api/v1/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req entities.UserCredentials

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	response, err := h.userUseCase.Login(req.Username, req.Password, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Register godoc
// @Summary Register new user
// @Description Create a new user account
// @Tags authentication
// @Accept json
// @Produce json
// @Param user body entities.UserCredentials true "User registration details"
// @Success 201 {object} entities.AuthResponse "JWT token"
// @Failure 400 {object} object{error=string} "Invalid request format or validation error"
// @Router /api/v1/auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req entities.UserCredentials

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	response, err := h.userUseCase.Register(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Logout godoc
// @Summary User logout
// @Description Invalidate user's current session
// @Tags authentication
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Security BearerAuth
// @Success 200 {object} object{message=string} "Successfully logged out"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Router /api/v1/auth/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.userUseCase.Logout(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// GetUser godoc
// @Summary Get user details
// @Description Get details of a specific user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Security BearerAuth
// @Success 200 {object} entities.UserResponse
// @Failure 404 {object} object{error=string} "User not found"
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.userUseCase.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ListUsers godoc
// @Summary List all users
// @Description Get paginated list of all users
// @Tags users
// @Accept json
// @Produce json
// @Param page query integer false "Page number" default(1)
// @Param page_size query integer false "Items per page" default(10)
// @Security BearerAuth
// @Success 200 {object} object{users=[]entities.UserResponse,total=integer,page=integer,page_size=integer}
// @Failure 500 {object} object{error=string} "Server error"
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, err := h.userUseCase.ListUsers(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":     users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetUserRoles godoc
// @Summary Get user roles
// @Description Get all roles assigned to a specific user
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Security BearerAuth
// @Success 200 {array} entities.Role
// @Failure 404 {object} object{error=string} "User not found"
// @Router /api/v1/users/{id}/roles [get]
func (h *UserHandler) GetUserRoles(c *gin.Context) {
	userID := c.Param("id")

	roles, err := h.userUseCase.GetUserRoles(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// AssignRole godoc
// @Summary Assign role to user
// @Description Assign a specific role to a user
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param roleId path string true "Role ID"
// @Security BearerAuth
// @Success 200 {object} object{message=string} "Role assigned successfully"
// @Failure 400 {object} object{error=string} "Invalid request"
// @Router /api/v1/users/{id}/roles/{roleId} [post]
func (h *UserHandler) AssignRole(c *gin.Context) {
	userID := c.Param("id")
	roleID := c.Param("roleId")

	if err := h.userUseCase.AssignRole(userID, roleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assigned successfully"})
}

// RemoveRole godoc
// @Summary Remove role from user
// @Description Remove a specific role from a user
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param roleId path string true "Role ID"
// @Security BearerAuth
// @Success 200 {object} object{message=string} "Role removed successfully"
// @Failure 400 {object} object{error=string} "Invalid request"
// @Router /api/v1/users/{id}/roles/{roleId} [delete]
func (h *UserHandler) RemoveRole(c *gin.Context) {
	userID := c.Param("id")
	roleID := c.Param("roleId")

	if err := h.userUseCase.RemoveRole(userID, roleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role removed successfully"})
}

// CreateAPIKey godoc
// @Summary Create API key
// @Description Create a new API key for the authenticated user
// @Tags api-keys
// @Accept json
// @Produce json
// @Param request body object{name=string} true "API key details"
// @Security BearerAuth
// @Success 201 {object} entities.APIKey
// @Failure 400 {object} object{error=string} "Invalid request format"
// @Failure 500 {object} object{error=string} "Server error"
// @Router /api/v1/users/api-keys [post]
func (h *UserHandler) CreateAPIKey(c *gin.Context) {
	userID := c.GetString("user_id")
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	apiKey, err := h.userUseCase.CreateAPIKey(userID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create API key"})
		return
	}

	c.JSON(http.StatusCreated, apiKey)
}

// ListAPIKeys godoc
// @Summary List API keys
// @Description Get all API keys for the authenticated user
// @Tags api-keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entities.APIKey
// @Failure 500 {object} object{error=string} "Server error"
// @Router /api/v1/users/api-keys [get]
func (h *UserHandler) ListAPIKeys(c *gin.Context) {
	userID := c.GetString("user_id")

	keys, err := h.userUseCase.ListAPIKeys(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch API keys"})
		return
	}

	c.JSON(http.StatusOK, keys)
}

// RevokeAPIKey godoc
// @Summary Revoke API key
// @Description Revoke a specific API key
// @Tags api-keys
// @Accept json
// @Produce json
// @Param keyId path string true "API Key ID"
// @Security BearerAuth
// @Success 200 {object} object{message=string} "API key revoked successfully"
// @Failure 400 {object} object{error=string} "Invalid request"
// @Router /api/v1/users/api-keys/{keyId} [delete]
func (h *UserHandler) RevokeAPIKey(c *gin.Context) {
	userID := c.GetString("user_id")
	keyID := c.Param("keyId")

	if err := h.userUseCase.RevokeAPIKey(userID, keyID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key revoked successfully"})
}

// GetAuditLogs godoc
// @Summary Get user audit logs
// @Description Get paginated audit logs for the authenticated user
// @Tags audit
// @Accept json
// @Produce json
// @Param page query integer false "Page number" default(1)
// @Param page_size query integer false "Items per page" default(10)
// @Security BearerAuth
// @Success 200 {object} object{logs=[]entities.AuditLog,total=integer,page=integer,page_size=integer}
// @Failure 500 {object} object{error=string} "Server error"
// @Router /api/v1/users/audit-logs [get]
func (h *UserHandler) GetAuditLogs(c *gin.Context) {
	userID := c.GetString("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	logs, total, err := h.userUseCase.GetAuditLogs(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch audit logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":      logs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetMe godoc
// @Summary Get current user details
// @Description Get details of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} entities.UserResponse
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Router /api/v1/users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	response, err := h.userUseCase.GetMe(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListRoles godoc
// @Summary List all roles
// @Description Get all available roles in the system
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entities.Role
// @Failure 500 {object} object{error=string} "Server error"
// @Router /api/v1/users/roles [get]
func (h *UserHandler) ListRoles(c *gin.Context) {
	roles, err := h.userUseCase.ListRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch roles"})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// ListPermissions godoc
// @Summary List all permissions
// @Description Get all available permissions in the system
// @Tags permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entities.Permission
// @Failure 500 {object} object{error=string} "Server error"
// @Router /api/v1/users/permissions [get]
func (h *UserHandler) ListPermissions(c *gin.Context) {
	permissions, err := h.userUseCase.ListPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch permissions"})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// GetUserPermissions godoc
// @Summary Get user permissions
// @Description Get all permissions assigned to a specific user through their roles
// @Tags permissions
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Security BearerAuth
// @Success 200 {array} entities.Permission
// @Failure 404 {object} object{error=string} "User not found"
// @Router /api/v1/users/{id}/permissions [get]
func (h *UserHandler) GetUserPermissions(c *gin.Context) {
	userID := c.Param("id")

	permissions, err := h.userUseCase.GetUserPermissions(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// AssignPermission godoc
// @Summary Assign permission to role
// @Description Assign a specific permission to a role
// @Tags permissions
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Param permissionId path string true "Permission ID"
// @Security BearerAuth
// @Success 200 {object} object{message=string} "Permission assigned successfully"
// @Failure 400 {object} object{error=string} "Invalid request"
// @Router /api/v1/users/roles/{id}/permissions/{permissionId} [post]
func (h *UserHandler) AssignPermission(c *gin.Context) {
	roleID := c.Param("id")
	permissionID := c.Param("permissionId")

	if err := h.userUseCase.AssignPermission(roleID, permissionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission assigned successfully"})
}

// RemovePermission godoc
// @Summary Remove permission from role
// @Description Remove a specific permission from a role
// @Tags permissions
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Param permissionId path string true "Permission ID"
// @Security BearerAuth
// @Success 200 {object} object{message=string} "Permission removed successfully"
// @Failure 400 {object} object{error=string} "Invalid request"
// @Router /api/v1/users/roles/{id}/permissions/{permissionId} [delete]
func (h *UserHandler) RemovePermission(c *gin.Context) {
	roleID := c.Param("id")
	permissionID := c.Param("permissionId")

	if err := h.userUseCase.RemovePermission(roleID, permissionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission removed successfully"})
}

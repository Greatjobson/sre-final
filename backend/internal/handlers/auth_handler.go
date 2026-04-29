package handlers

import (
	"net/http"

	"github.com/Tedra-ez/AdvancedProgramming_Final/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth *services.AuthService
}

func NewAuthHandler(auth *services.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var fullName, email, password string
	if c.GetHeader("Content-Type") == "application/json" {
		var req struct {
			FullName string `json:"fullName"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		fullName, email, password = req.FullName, req.Email, req.Password
	} else {
		fullName = c.PostForm("fullName")
		email = c.PostForm("email")
		password = c.PostForm("password")
		if fullName == "" || email == "" || password == "" {
			c.Redirect(http.StatusFound, "/register?error=invalid+input")
			return
		}
	}
	err := h.auth.Register(c.Request.Context(), fullName, email, password)
	if err != nil {
		if c.GetHeader("Content-Type") == "application/json" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.Redirect(http.StatusFound, "/register?error=email+exists")
		return
	}
	if c.GetHeader("Content-Type") == "application/json" {
		c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
		return
	}
	c.Redirect(http.StatusFound, "/login")
}

func (h *AuthHandler) Login(c *gin.Context) {
	var email, password string
	if c.GetHeader("Content-Type") == "application/json" {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		email, password = req.Email, req.Password
	} else {
		email = c.PostForm("email")
		password = c.PostForm("password")
		if email == "" || password == "" {
			c.Redirect(http.StatusFound, "/login?error=invalid+input")
			return
		}
	}
	token, err := h.auth.Login(c.Request.Context(), email, password)
	if err != nil {
		if c.GetHeader("Content-Type") == "application/json" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.Redirect(http.StatusFound, "/login?error=invalid+credentials")
		return
	}
	c.SetCookie("auth_token", token, 24*3600, "/", "", false, true)
	if c.GetHeader("Content-Type") == "application/json" {
		c.JSON(http.StatusOK, gin.H{"token": token})
		return
	}
	c.Redirect(http.StatusFound, "/account")
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	token, err := c.Cookie("auth_token")
	if err != nil || token == "" {
		c.JSON(http.StatusOK, gin.H{"message": "token refresh check", "status": "ok"})
		return
	}

	newToken, err := h.auth.RefreshToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "token refresh check", "status": "ok"})
		return
	}

	c.SetCookie("auth_token", newToken, 24*3600, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"token": newToken})
}

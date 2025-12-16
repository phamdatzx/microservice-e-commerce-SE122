package middleware

import (
	"net/http"
	"user-service/error"

	"github.com/gin-gonic/gin"
)

const (
	RoleAdmin    = "admin"
	RoleSeller   = "seller"
	RoleCustomer = "customer"
)

// RequireAdmin middleware ensures only admin users can access the route
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetHeader("X-User-Role")
		
		if role == "" {
			c.Error(error.NewAppError(http.StatusUnauthorized, "Authentication required"))
			c.Abort()
			return
		}

		if role != RoleAdmin {
			c.Error(error.NewAppError(http.StatusForbidden, "Admin access required"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireSeller middleware ensures only seller users can access the route
func RequireSeller() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetHeader("X-User-Role")
		
		if role == "" {
			c.Error(error.NewAppError(http.StatusUnauthorized, "Authentication required"))
			c.Abort()
			return
		}

		if role != RoleSeller {
			c.Error(error.NewAppError(http.StatusForbidden, "Seller access required"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireSellerOrAdmin middleware ensures only seller or admin users can access the route
func RequireSellerOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetHeader("X-User-Role")
		
		if role == "" {
			c.Error(error.NewAppError(http.StatusUnauthorized, "Authentication required"))
			c.Abort()
			return
		}

		if role != RoleSeller && role != RoleAdmin {
			c.Error(error.NewAppError(http.StatusForbidden, "Seller or Admin access required"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireCustomer middleware ensures only customer users can access the route
func RequireCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetHeader("X-User-Role")
		
		if role == "" {
			c.Error(error.NewAppError(http.StatusUnauthorized, "Authentication required"))
			c.Abort()
			return
		}

		if role != RoleCustomer {
			c.Error(error.NewAppError(http.StatusForbidden, "Customer access required"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAuthenticated middleware ensures the user has any valid role
func RequireAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetHeader("X-User-Role")
		
		if role == "" {
			c.Error(error.NewAppError(http.StatusUnauthorized, "Authentication required"))
			c.Abort()
			return
		}

		// Validate that the role is one of the known roles
		if role != RoleAdmin && role != RoleSeller && role != RoleCustomer {
			c.Error(error.NewAppError(http.StatusUnauthorized, "Invalid role"))
			c.Abort()
			return
		}

		c.Next()
	}
}

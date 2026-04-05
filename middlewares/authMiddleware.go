package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set Secret Key dari env
		clerk.SetKey(os.Getenv("CLERK_SECRET_KEY"))

		// 1. Cek keberadaan header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Akses Ditolak: Token tidak ditemukan"})
			return
		}

		// 2. Format token harus "Bearer <token>"
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// 3. Verifikasi token ke server Clerk
		claims, err := jwt.Verify(c.Request.Context(), &jwt.VerifyParams{
			Token: token,
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Akses Ditolak: Token kedaluwarsa atau tidak valid"})
			return
		}

		// 4. (Opsional) Simpan ID User dari token ke dalam context jika ingin dipakai di Controller
		c.Set("userID", claims.Subject)

		// Lolos verifikasi, lanjutkan ke controller tujuan
		c.Next()
	}
}
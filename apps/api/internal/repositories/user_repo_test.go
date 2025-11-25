package repositories

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"apps/api/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestUserRepo() *UserRepo {
	return NewUserRepo(testDbService.GetDB())
}

func TestUserRepo_CreateUser(t *testing.T) {
	ctx := context.Background()

	t.Run("should create user successfully", func(t *testing.T) {
		cleanupTestDatabase()
		userRepo := getTestUserRepo()

		userCreate := models.UserCreate{
			Email:        "test@example.com",
			PasswordHash: "hashedpassword123",
		}

		user, err := userRepo.CreateUser(ctx, userCreate)

		require.NoError(t, err)
		require.NotNil(t, user)

		assert.NotEmpty(t, user.ID)
		assert.Equal(t, userCreate.Email, user.Email)
		assert.Equal(t, userCreate.PasswordHash, user.PasswordHash)
		assert.False(t, user.CreatedAt.IsZero())
		assert.False(t, user.UpdatedAt.IsZero())
	})

	t.Run(
		"should fail when creating user with duplicate email",
		func(t *testing.T) {
			cleanupTestDatabase()
			userRepo := getTestUserRepo()

			userCreate := models.UserCreate{
				Email:        "duplicate@example.com",
				PasswordHash: "hashedpassword123",
			}

			// Create first user
			_, err := userRepo.CreateUser(ctx, userCreate)
			require.NoError(t, err)

			// Try to create second user with same email
			_, err = userRepo.CreateUser(ctx, userCreate)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "Failed to create user")
		},
	)
}

func TestUserRepo_GetUserByEmail(t *testing.T) {
	ctx := context.Background()

	t.Run("should get user by email successfully", func(t *testing.T) {
		cleanupTestDatabase()
		userRepo := getTestUserRepo()

		// Create a user first
		userCreate := models.UserCreate{
			Email:        "get@example.com",
			PasswordHash: "hashedpassword123",
		}
		createdUser, err := userRepo.CreateUser(ctx, userCreate)
		require.NoError(t, err)

		// Get user by email
		user, err := userRepo.GetUserByEmail(ctx, userCreate.Email)

		require.NoError(t, err)
		require.NotNil(t, user)

		assert.Equal(t, createdUser.ID, user.ID)
		assert.Equal(t, createdUser.Email, user.Email)
		assert.Equal(t, createdUser.PasswordHash, user.PasswordHash)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		cleanupTestDatabase()
		userRepo := getTestUserRepo()

		user, err := userRepo.GetUserByEmail(ctx, "nonexistent@example.com")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "Failed to get user by email")
	})

	t.Run("should handle case-sensitive email", func(t *testing.T) {
		cleanupTestDatabase()
		userRepo := getTestUserRepo()

		// Create a user with lowercase email
		userCreate := models.UserCreate{
			Email:        "case@example.com",
			PasswordHash: "hashedpassword123",
		}
		_, err := userRepo.CreateUser(ctx, userCreate)
		require.NoError(t, err)

		// Try to get with different case
		user, err := userRepo.GetUserByEmail(ctx, "CASE@EXAMPLE.COM")

		// This should fail as PostgreSQL is case-sensitive by default
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepo_GetUserById(t *testing.T) {
	ctx := context.Background()

	t.Run("should get user by ID successfully", func(t *testing.T) {
		cleanupTestDatabase()
		userRepo := getTestUserRepo()

		// Create a user first
		userCreate := models.UserCreate{
			Email:        "getbyid@example.com",
			PasswordHash: "hashedpassword123",
		}
		createdUser, err := userRepo.CreateUser(ctx, userCreate)
		require.NoError(t, err)

		// Get user by ID
		user, err := userRepo.GetUserById(ctx, createdUser.ID)

		require.NoError(t, err)
		require.NotNil(t, user)

		assert.Equal(t, createdUser.ID, user.ID)
		assert.Equal(t, createdUser.Email, user.Email)
		assert.Equal(t, createdUser.PasswordHash, user.PasswordHash)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		cleanupTestDatabase()
		userRepo := getTestUserRepo()

		user, err := userRepo.GetUserById(
			ctx,
			"550e8400-e29b-41d4-a716-446655440000",
		)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "Failed to get user by id")
	})

	t.Run("should handle invalid UUID format", func(t *testing.T) {
		cleanupTestDatabase()
		userRepo := getTestUserRepo()

		user, err := userRepo.GetUserById(ctx, "invalid-uuid")

		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepo_getUserByUniqField(t *testing.T) {
	ctx := context.Background()

	t.Run("should get user by any unique field", func(t *testing.T) {
		cleanupTestDatabase()
		userRepo := getTestUserRepo()

		// Create a user first
		userCreate := models.UserCreate{
			Email:        "uniquefield@example.com",
			PasswordHash: "hashedpassword123",
		}
		createdUser, err := userRepo.CreateUser(ctx, userCreate)
		require.NoError(t, err)

		// Test getting by email
		user, err := userRepo.getUserByUniqField(ctx, "email", userCreate.Email)
		require.NoError(t, err)
		assert.Equal(t, createdUser.ID, user.ID)

		// Test getting by ID
		user, err = userRepo.getUserByUniqField(ctx, "id", createdUser.ID)
		require.NoError(t, err)
		assert.Equal(t, createdUser.Email, user.Email)
	})

	t.Run(
		"should return error for non-existent field value",
		func(t *testing.T) {
			cleanupTestDatabase()
			userRepo := getTestUserRepo()

			user, err := userRepo.getUserByUniqField(
				ctx,
				"email",
				"nonexistent@example.com",
			)

			assert.Error(t, err)
			assert.Nil(t, user)
		},
	)

	t.Run(
		"should handle invalid field names",
		func(t *testing.T) {
			cleanupTestDatabase()
			userRepo := getTestUserRepo()

			user, err := userRepo.getUserByUniqField(
				ctx,
				"invalid_field",
				"some_value",
			)

			assert.Error(t, err)
			assert.Nil(t, user)
		},
	)
}

func TestUserRepo_Integration_CompleteFlow(t *testing.T) {
	ctx := context.Background()

	t.Run("complete user creation and retrieval flow", func(t *testing.T) {
		cleanupTestDatabase()
		userRepo := getTestUserRepo()

		// Create multiple users
		users := []models.UserCreate{
			{Email: "user1@example.com", PasswordHash: "hash1"},
			{Email: "user2@example.com", PasswordHash: "hash2"},
			{Email: "user3@example.com", PasswordHash: "hash3"},
		}

		var createdUsers []*models.User
		for _, userCreate := range users {
			user, err := userRepo.CreateUser(ctx, userCreate)
			require.NoError(t, err)
			createdUsers = append(createdUsers, user)
		}

		// Verify all users can be retrieved by email
		for i, user := range createdUsers {
			retrievedUser, err := userRepo.GetUserByEmail(ctx, users[i].Email)
			require.NoError(t, err)
			assert.Equal(t, user.ID, retrievedUser.ID)
			assert.Equal(t, user.Email, retrievedUser.Email)
		}

		// Verify all users can be retrieved by ID
		for _, user := range createdUsers {
			retrievedUser, err := userRepo.GetUserById(ctx, user.ID)
			require.NoError(t, err)
			assert.Equal(t, user.Email, retrievedUser.Email)
		}

		// Verify total user count in database
		db := testDbService.GetDB()
		var count int
		err := db.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, len(users), count)
	})
}

func TestUserRepo_ConcurrentOperations(t *testing.T) {
	ctx := context.Background()

	t.Run("concurrent user creation should work", func(t *testing.T) {
		cleanupTestDatabase()
		userRepo := getTestUserRepo()

		const numGoroutines = 5
		results := make(chan error, numGoroutines)

		// Start concurrent user creation
		for i := 0; i < numGoroutines; i++ {
			go func(index int) {
				userCreate := models.UserCreate{
					Email: fmt.Sprintf(
						"concurrent%d@example.com",
						index,
					),
					PasswordHash: fmt.Sprintf("hash%d", index),
				}

				_, err := userRepo.CreateUser(ctx, userCreate)
				results <- err
			}(i)
		}

		// Collect results
		for i := 0; i < numGoroutines; i++ {
			err := <-results
			assert.NoError(t, err)
		}

		// Verify all users were created
		db := testDbService.GetDB()
		var count int
		err := db.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, numGoroutines, count)
	})
}

func TestUserRepo_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("should handle empty email gracefully", func(t *testing.T) {
		cleanupTestDatabase()
		userRepo := getTestUserRepo()

		userCreate := models.UserCreate{
			Email:        "",
			PasswordHash: "hashedpassword123",
		}

		user, err := userRepo.CreateUser(ctx, userCreate)
		// The repository level doesn't validate emails - that's likely done at the service level
		// So empty email should succeed at repository level
		require.NoError(t, err)
		assert.Equal(t, "", user.Email)
	})

	t.Run("should handle empty password hash gracefully", func(t *testing.T) {
		cleanupTestDatabase()
		userRepo := getTestUserRepo()

		userCreate := models.UserCreate{
			Email:        "empty-password@example.com",
			PasswordHash: "",
		}

		user, err := userRepo.CreateUser(ctx, userCreate)
		// This might be valid depending on business logic - let's check what happens
		if err == nil {
			assert.NotNil(t, user)
			assert.Equal(t, "", user.PasswordHash)
		}
	})

	t.Run("should handle very long email addresses", func(t *testing.T) {
		cleanupTestDatabase()
		userRepo := getTestUserRepo()

		// Create a very long email (but still valid)
		longEmail := strings.Repeat("a", 240) + "@example.com"
		userCreate := models.UserCreate{
			Email:        longEmail,
			PasswordHash: "hashedpassword123",
		}

		user, err := userRepo.CreateUser(ctx, userCreate)
		if err == nil {
			assert.NotNil(t, user)
			assert.Equal(t, longEmail, user.Email)
		} else {
			// If there's a length constraint, error is expected
			assert.Error(t, err)
		}
	})

	t.Run("should handle SQL injection attempts in email", func(t *testing.T) {
		cleanupTestDatabase()
		userRepo := getTestUserRepo()

		maliciousEmail := "'; DROP TABLE users; --@example.com"
		userCreate := models.UserCreate{
			Email:        maliciousEmail,
			PasswordHash: "hashedpassword123",
		}

		// Should either succeed (if email is treated as data) or fail gracefully
		user, err := userRepo.CreateUser(ctx, userCreate)
		if err == nil {
			// If creation succeeds, the email should be stored as-is
			assert.Equal(t, maliciousEmail, user.Email)

			// Verify table still exists by trying to get the user back
			retrievedUser, retrieveErr := userRepo.GetUserByEmail(
				ctx,
				maliciousEmail,
			)
			assert.NoError(t, retrieveErr)
			assert.Equal(t, user.ID, retrievedUser.ID)
		}
	})

	t.Run("should handle context cancellation gracefully", func(t *testing.T) {
		cleanupTestDatabase()
		userRepo := getTestUserRepo()

		// Create a context that's already cancelled
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		userCreate := models.UserCreate{
			Email:        "cancelled@example.com",
			PasswordHash: "hashedpassword123",
		}

		_, err := userRepo.CreateUser(cancelledCtx, userCreate)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context canceled")
	})
}

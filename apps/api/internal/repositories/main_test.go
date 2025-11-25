package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"apps/api/internal/config"
	"apps/api/internal/database"
	"apps/api/internal/models"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDbConfig *config.DbConfig
var testDbService database.Service

func mustStartPostgresContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	var (
		dbName = "testdb"
		dbPwd  = "password"
		dbUser = "user"
	)

	dbContainer, err := tcpostgres.Run(
		context.Background(),
		"postgres:latest",
		tcpostgres.WithDatabase(dbName),
		tcpostgres.WithUsername(dbUser),
		tcpostgres.WithPassword(dbPwd),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}

	testDbConfig = &config.DbConfig{
		DbHost:     dbHost,
		DbPort:     dbPort.Int(),
		DbName:     dbName,
		DbUsername: dbUser,
		DbPassword: dbPwd,
		DbSchema:   "public",
	}

	return dbContainer.Terminate, nil
}

func setupTestDatabase() error {
	// Use standard database/sql DB for migrations
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		testDbConfig.DbUsername,
		testDbConfig.DbPassword,
		testDbConfig.DbHost,
		testDbConfig.DbPort,
		testDbConfig.DbName,
	)

	// Open *sql.DB for migration
	sqlDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	driver, err := migratepg.WithInstance(sqlDB, &migratepg.Config{})
	if err != nil {
		return err
	}

	migrationsPath := "file://../database/migrations"

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		return err
	}

	// Now initialize your app DB service for use in tests
	testDbService = database.New(testDbConfig)
	return nil
}

func cleanupTestDatabase() {
	db := testDbService.GetDB()
	_, _ = db.Exec(
		context.Background(),
		"TRUNCATE TABLE users, folders CASCADE",
	)
}

func TestMain(m *testing.M) {
	teardown, err := mustStartPostgresContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	if err := setupTestDatabase(); err != nil {
		log.Fatalf("could not setup test database: %v", err)
	}

	code := m.Run()

	if teardown != nil {
		if err := teardown(context.Background()); err != nil {
			log.Fatalf("could not teardown postgres container: %v", err)
		}
	}

	log.Printf("Tests completed with exit code %d", code)
}

func createTestUSer(t *testing.T, ctx context.Context) *models.User {
	userRepo := getTestUserRepo()

	userCreate := models.UserCreate{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword123",
	}

	user, err := userRepo.CreateUser(ctx, userCreate)

	require.NoError(t, err)
	require.NotNil(t, user)

	assert.NotEmpty(t, user.ID)

	return user
}

func createTestFolder(
	t *testing.T,
	ctx context.Context,
	userId string,
) *models.Folder {
	folderRepo := getTestFolderRepo()

	folderCreate := models.FolderCreate{
		Name:     "Test Folder",
		ParentId: nil,
		UserId:   userId,
	}

	folder, err := folderRepo.CreateFolder(ctx, folderCreate)

	require.NoError(t, err)
	require.NotNil(t, folder)

	assert.NotEmpty(t, folder.ID)
	assert.Equal(t, folderCreate.UserId, userId)

	return folder
}

func createTestLink(
	t *testing.T,
	ctx context.Context,
	userId, folderId string,
) *models.Link {
	linkRepo := getTestLinkRepo()

	linkCreate := models.LinkCreate{
		Url:         "https://example.com",
		Name:        "Example Link",
		Description: "A test link for example.com",
		UserId:      userId,
		FolderId:    folderId,
	}

	link, err := linkRepo.CreateLink(ctx, linkCreate)
	require.NoError(t, err)
	require.NotNil(t, link)

	assert.NotEmpty(t, link.ID)
	assert.Equal(t, linkCreate.UserId, userId)
	assert.Equal(t, linkCreate.FolderId, folderId)

	return link
}

func assertLinkFolder(
	t *testing.T,
	ctx context.Context,
	linkId, expectedFolderId string,
) {
	linkRepo := getTestLinkRepo()

	link, err := linkRepo.GetLinkById(ctx, linkId)
	require.NoError(t, err)
	require.NotNil(t, link)

	assert.Equal(t, expectedFolderId, link.FolderId)
}

func assertFolderParent(
	t *testing.T,
	ctx context.Context,
	folderId, expectedParentId string,
) {
	folderRepo := getTestFolderRepo()

	folder, err := folderRepo.GetFolderById(ctx, folderId)
	require.NoError(t, err)
	require.NotNil(t, folder)

	if folder.ParentId != nil {
		assert.Equal(t, expectedParentId, *folder.ParentId)
	} else {
		assert.Fail(t, "ParentId is nil")
	}
}

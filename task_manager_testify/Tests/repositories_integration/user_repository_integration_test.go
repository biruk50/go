package repositories_integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"task_manager_testify/Domain"
	"task_manager_testify/Infrastructure"
	"task_manager_testify/Repositories"
)

type UserRepoIntegrationSuite struct {
	suite.Suite
	repo Repositories.UserRepository
	db   interface{} // any type so go test doesn’t complain
}

func (s *UserRepoIntegrationSuite) SetupSuite() {
	Infrastructure.LoadEnv()

	uri := Infrastructure.GetEnv("MONGODB_URL", "")
	if uri == "" {
		s.T().Skip("MONGODB_URL not set – skipping integration test")
	}

	err := Infrastructure.InitMongo()
	s.Require().NoError(err)

	s.db = Infrastructure.GetDB()
	s.repo = Repositories.NewMongoUserRepository(Infrastructure.GetDB())
}

func (s *UserRepoIntegrationSuite) TearDownSuite() {
	Infrastructure.CloseMongo()
}

func (s *UserRepoIntegrationSuite) SetupTest() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	Infrastructure.GetDB().Collection("users").Drop(ctx)
}

func (s *UserRepoIntegrationSuite) TestCreateAndFind() {
	user := &Domain.User{
		Username:     "alice",
		PasswordHash: "hash",
		Role:         "user",
	}

	err := s.repo.Create(user)
	s.NoError(err)
	s.NotEmpty(user.ID)

	found, err := s.repo.FindByUsername("alice")
	s.NoError(err)
	s.Equal("alice", found.Username)
}

func (s *UserRepoIntegrationSuite) TestPromote() {
	user := &Domain.User{
		Username:     "bob",
		PasswordHash: "hash",
		Role:         "user",
	}
	s.repo.Create(user)

	err := s.repo.Promote("bob")
	s.NoError(err)

	found, err := s.repo.FindByUsername("bob")
	s.NoError(err)
	s.Equal("admin", found.Role)
}

func TestUserRepoIntegrationSuite(t *testing.T) {
	suite.Run(t, new(UserRepoIntegrationSuite))
}

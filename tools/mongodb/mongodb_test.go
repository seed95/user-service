package mongodb_test

import (
	"testing"
	"user-service/internal/user"
	"user-service/pkg"
	"user-service/test/mock"
	"user-service/tools/dockertest"
	"user-service/tools/mongodb"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var mongoConfig = mock.MongoConfig

func TestMongodb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mongodb Suite")
}

var _ = Describe("MongoDB", Ordered, func() {

	var mongo mongodb.MongoDBInterface
	var dockerContainer *dockertest.Dockertest

	BeforeAll(func() {
		dockerContainer = dockertest.NewDockertest("")
		err := dockerContainer.RunMongoDB(mongoConfig)
		Expect(err).Should(BeNil())
	})

	AfterAll(func() {
		dockerContainer.Purge()
	})

	Context("Connect", func() {
		It("Should return database", func() {
			mongo = mongodb.NewMongoDB(mongoConfig)
			db, err := mongo.Connect()
			Expect(err).Should(BeNil())
			Expect(db).ShouldNot(BeNil())
		})
	})

	Context("GetMongoDBURI", func() {
		It("should return mongodb uri", func() {
			mongodbURI := mongo.GetMongoDBURI()
			Expect(mongodbURI).Should(Equal("mongodb://" + mongoConfig.Username + ":" + mongoConfig.Password + "@" + mongoConfig.Host + ":" + mongoConfig.Port))
		})
	})

	Context("GetUserCollection", func() {
		It("should return mongodb user collection", func() {
			userCollection := mongo.GetUserCollection()
			Expect(userCollection).ShouldNot(BeNil())
		})
	})

	Context("GetDatabase", func() {
		It("should return mongodb database", func() {
			database := mongo.GetDatabase()
			Expect(database).ShouldNot(BeNil())
		})
	})

	When("User collection is not empty", func() {
		Context("Upsert", func() {
			It("should upsert user", func() {
				user := mock.MockUser
				err := mongo.Upsert(user)
				Expect(err).Should(BeNil())
			})
		})

		Context("IsEmailExists", func() {
			It("should return true if email exists", func() {
				email := mock.MockUser.Email
				exists, err := mongo.IsEmailExists(email)
				Expect(err).Should(BeNil())
				Expect(exists).Should(BeTrue())
			})
		})

		Context("IsNicknameExists", func() {
			It("should return true if nickname exists", func() {
				nickname := mock.MockUser.Nickname
				exists, err := mongo.IsNicknameExists(nickname)
				Expect(err).Should(BeNil())
				Expect(exists).Should(BeTrue())
			})
		})

		Context("GetUserByID", func() {
			It("should return user", func() {
				user, err := mongo.GetUserByID(mock.MockUser.ID)
				Expect(err).Should(BeNil())
				Expect(user.ID).Should(Equal(mock.MockUser.ID))
				Expect(user.Email).Should(Equal(mock.MockUser.Email))
				Expect(user.Nickname).Should(Equal(mock.MockUser.Nickname))
				Expect(user.FirstName).Should(Equal(mock.MockUser.FirstName))
				Expect(user.LastName).Should(Equal(mock.MockUser.LastName))
			})
		})

		Context("GetUsers", func() {
			It("should return users", func() {
				users, err := mongo.GetUsers()
				Expect(err).Should(BeNil())
				Expect(len(users)).Should(Equal(1))
				Expect(users[0].ID).Should(Equal(mock.MockUser.ID))
				Expect(users[0].Email).Should(Equal(mock.MockUser.Email))
				Expect(users[0].Nickname).Should(Equal(mock.MockUser.Nickname))
				Expect(users[0].FirstName).Should(Equal(mock.MockUser.FirstName))
				Expect(users[0].LastName).Should(Equal(mock.MockUser.LastName))
			})
		})

		Context("DeleteUserByID", func() {
			It("should delete user", func() {
				err := mongo.DeleteUserByID(mock.MockUser.ID)
				Expect(err).Should(BeNil())
			})
		})
	})

	When("User collection is empty", func() {
		Context("IsEmailExists", func() {
			It("should return true if email exists", func() {
				email := mock.MockUser.Email
				exists, err := mongo.IsEmailExists(email)
				Expect(err).Should(BeNil())
				Expect(exists).Should(BeFalse())
			})
		})

		Context("IsNicknameExists", func() {
			It("should return true if nickname exists", func() {
				nickname := mock.MockUser.Nickname
				exists, err := mongo.IsNicknameExists(nickname)
				Expect(err).Should(BeNil())
				Expect(exists).Should(BeFalse())
			})
		})

		Context("GetUserByID", func() {
			It("should return user", func() {
				userObj, err := mongo.GetUserByID(mock.MockUser.ID)
				Expect(err).ShouldNot(BeNil())
				Expect(err).Should(Equal(pkg.ErrUserNotFound))
				Expect(userObj).Should(Equal(new(user.User)))
			})
		})

		Context("GetUsers", func() {
			It("should return users", func() {
				users, err := mongo.GetUsers()
				Expect(err).ShouldNot(BeNil())
				Expect(err).Should(Equal(pkg.ErrUserNotFound))
				Expect(users).Should(BeNil())
			})
		})
	})

	Context("CreateUsers", func() {
		It("should create users", func() {
			users := mock.MockUsers
			err := mongo.CreateUsers(users)

			Expect(err).Should(BeNil())
		})

		It("should return error if users is empty", func() {
			err := mongo.CreateUsers([]user.User{})

			Expect(err).ShouldNot(BeNil())
			Expect(err).Should(Equal(pkg.ErrUsersEmpty))
		})

		It("should return error if users is nil", func() {
			err := mongo.CreateUsers(nil)

			Expect(err).ShouldNot(BeNil())
			Expect(err).Should(Equal(pkg.ErrUsersEmpty))
		})
	})
})
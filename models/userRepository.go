package model

import (
	"time"

	"gorm.io/gorm"

	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
)

type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Name         string    `gorm:"type:varchar(100);not null"`
	Email        string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password     string    `gorm:"type:varchar(255);not null"`
	UserName     string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	PhoneNumber  string    `gorm:"type:varchar(15);uniqueIndex;not null"`
	Cash         uint32    `gorm:"default:0"`
	Bronze       uint32    `gorm:"default:0"`
	Silver       uint32    `gorm:"default:0"`
	Gold         uint32    `gorm:"default:0"`
	IsAuthorized bool      `gorm:"default:false"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

type mysqlUserRepo struct {
	db *gorm.DB
}

type UserRepository interface {
	FindByUserName(username string) (User, *helpers.CustomEror)
	RegisterUser(tx *gorm.DB, user User) (err *helpers.CustomEror)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &mysqlUserRepo{
		db: db,
	}
}

func (repo *mysqlUserRepo) FindByUserName(userName string) (user User, err *helpers.CustomEror) {
	return
}

func (repo *mysqlUserRepo) RegisterUser(tx *gorm.DB, user User) (err *helpers.CustomEror) {
	derr := repo.db.Create(user).Error
	if derr != nil {
		return err.System(derr.Error())
	}
	return nil
}

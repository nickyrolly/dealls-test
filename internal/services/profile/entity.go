package profile

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserProfile struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;column:id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;column:user_id"`
	Height     int       `gorm:"type:int;column:height"`
	Weight     int       `gorm:"type:int;column:weight"`
	Occupation string    `gorm:"type:varchar(100);column:occupation"`
	Education  string    `gorm:"type:varchar(100);column:education"`
	Religion   string    `gorm:"type:varchar(50);column:religion"`
	Ethnicity  string    `gorm:"type:varchar(50);column:ethnicity"`
	Interests  string    `gorm:"type:text;column:interests"`
	AboutMe    string    `gorm:"type:text;column:about_me"`
	CreatedAt  time.Time `gorm:"type:timestamp;not null;default:current_timestamp;column:created_at"`
	UpdatedAt  time.Time `gorm:"type:timestamp;not null;default:current_timestamp;column:updated_at"`
}

type UserPhoto struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;column:id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;column:user_id"`
	URL       string    `gorm:"type:varchar(255);not null;column:url"`
	IsPrimary bool      `gorm:"type:boolean;not null;default:false;column:is_primary"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:current_timestamp;column:created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:current_timestamp;column:updated_at"`
}

type UserPreference struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;column:id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;column:user_id"`
	Gender    string    `gorm:"type:varchar(20);not null;column:gender"`
	MinAge    int       `gorm:"type:int;column:min_age"`
	MaxAge    int       `gorm:"type:int;column:max_age"`
	MinHeight int       `gorm:"type:int;column:min_height"`
	MaxHeight int       `gorm:"type:int;column:max_height"`
	Religion  string    `gorm:"type:varchar(50);column:religion"`
	Ethnicity string    `gorm:"type:varchar(50);column:ethnicity"`
	Distance  int       `gorm:"type:int;column:distance"` // in kilometers
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:current_timestamp;column:created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:current_timestamp;column:updated_at"`
}

type UserMatch struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;column:id"`
	User1ID   uuid.UUID `gorm:"type:uuid;not null;column:user1_id"`
	User2ID   uuid.UUID `gorm:"type:uuid;not null;column:user2_id"`
	Status    string    `gorm:"type:varchar(20);not null;default:'pending';column:status"` // pending, matched, unmatched
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:current_timestamp;column:created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:current_timestamp;column:updated_at"`
}

type UserLike struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;column:id"`
	LikerID   uuid.UUID `gorm:"type:uuid;not null;column:liker_id"`
	LikedID   uuid.UUID `gorm:"type:uuid;not null;column:liked_id"`
	Status    string    `gorm:"type:varchar(20);not null;default:'active';column:status"` // active, withdrawn
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:current_timestamp;column:created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:current_timestamp;column:updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (e *UserProfile) BeforeCreate(tx *gorm.DB) error {
	e.ID = uuid.New()
	return nil
}

func (e *UserPhoto) BeforeCreate(tx *gorm.DB) error {
	e.ID = uuid.New()
	return nil
}

func (e *UserPreference) BeforeCreate(tx *gorm.DB) error {
	e.ID = uuid.New()
	return nil
}

func (e *UserMatch) BeforeCreate(tx *gorm.DB) error {
	e.ID = uuid.New()
	return nil
}

func (e *UserLike) BeforeCreate(tx *gorm.DB) error {
	e.ID = uuid.New()
	return nil
}

func (e *UserProfile) TableName() string {
	return "user_profiles"
}

func (e *UserPhoto) TableName() string {
	return "user_photos"
}

func (e *UserPreference) TableName() string {
	return "user_preferences"
}

func (e *UserMatch) TableName() string {
	return "user_matches"
}

func (e *UserLike) TableName() string {
	return "user_likes"
}

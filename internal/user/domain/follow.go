package domain

type Follow struc {
	ID int `gorm:"primaryKey"`
	FollowerID int
	Seller ID int
}
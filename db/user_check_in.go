package db

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// FortuneType represents the fortune assigned to a user's daily check-in.
type FortuneType int

const (
	FortuneTypeNone FortuneType = iota
	FortuneTypeGreatBlessing
	FortuneTypeMiddleBlessing
	FortuneTypeSmallBlessing
	FortuneTypeBlessing
	FortuneTypeFutureBlessing
	FortuneTypeBadLuck
	FortuneTypeGreatBadLuck
)

func (f FortuneType) Valid() bool {
	return f >= FortuneTypeGreatBlessing && f <= FortuneTypeGreatBadLuck
}

// UserCheckIn stores only the latest check-in state for a user.
// It intentionally does not retain daily check-in history.
type UserCheckIn struct {
	UserID          int         `gorm:"primaryKey;autoIncrement:false" json:"userId"`
	LastCheckInDate time.Time   `gorm:"type:date;not null" json:"lastCheckInDate"`
	CurrentStreak   int         `gorm:"not null;default:1" json:"currentStreak"`
	LastFortune     FortuneType `gorm:"not null" json:"lastFortune"`
	CreatedAt       time.Time   `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time   `gorm:"autoUpdateTime" json:"updatedAt"`

	User *User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

type UserCheckInResult struct {
	State            UserCheckIn
	AlreadyCheckedIn bool
}

var taipeiTimeZone = time.FixedZone("Asia/Taipei", 8*60*60)

// CheckInUser records a user's daily check-in using the Asia/Taipei calendar day.
// Repeated calls on the same day return the persisted fortune without incrementing the streak.
func CheckInUser(database *gorm.DB, userID int, fortune FortuneType, now time.Time) (UserCheckInResult, error) {
	var result UserCheckInResult
	if database == nil || userID <= 0 {
		return result, ErrParameterNotFound
	}
	if !fortune.Valid() {
		return result, fmt.Errorf("invalid fortune type: %d", fortune)
	}

	today := checkInDate(now)
	err := database.Transaction(func(tx *gorm.DB) error {
		// Lock the user row so concurrent first-time check-ins are serialized too.
		var user User
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("id").
			First(&user, userID).Error; err != nil {
			return err
		}

		var state UserCheckIn
		err := tx.First(&state, "user_id = ?", userID).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			state = UserCheckIn{
				UserID:          userID,
				LastCheckInDate: today,
				CurrentStreak:   1,
				LastFortune:     fortune,
			}
			if err := tx.Create(&state).Error; err != nil {
				return err
			}
			result.State = state
			return nil
		}

		if sameCalendarDate(state.LastCheckInDate, today) {
			result.State = state
			result.AlreadyCheckedIn = true
			return nil
		}

		state.CurrentStreak = nextCheckInStreak(state, today)
		state.LastCheckInDate = today
		state.LastFortune = fortune
		state.UpdatedAt = time.Now()
		if err := tx.Model(&state).
			Select("last_check_in_date", "current_streak", "last_fortune", "updated_at").
			Updates(&state).Error; err != nil {
			return err
		}

		result.State = state
		return nil
	})

	return result, err
}

func checkInDate(now time.Time) time.Time {
	localNow := now.In(taipeiTimeZone)
	return time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, time.UTC)
}

func nextCheckInStreak(state UserCheckIn, today time.Time) int {
	yesterday := today.AddDate(0, 0, -1)
	if sameCalendarDate(state.LastCheckInDate, yesterday) {
		return state.CurrentStreak + 1
	}
	return 1
}

func sameCalendarDate(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

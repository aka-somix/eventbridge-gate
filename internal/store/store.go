package store

import (
	"errors"
)

// ProfileStore manages the currently active AWS CLI profile.
type ProfileStore struct {
	activeProfile string
}

// SetActiveProfile sets the active AWS CLI profile.
func (ps *ProfileStore) SetActiveProfile(profile string) {
	ps.activeProfile = profile
}

// GetActiveProfile retrieves the currently active AWS CLI profile.
// Returns an error if no profile is set.
func (ps *ProfileStore) GetActiveProfile() (string, error) {
	if ps.IsThereAProfile() {
		return ps.activeProfile, nil
	}
	return "", errors.New("no active profile set")
}

// ClearActiveProfile clears the active AWS CLI profile.
func (ps *ProfileStore) ClearActiveProfile() {
	ps.activeProfile = ""
}

func (ps *ProfileStore) IsThereAProfile() bool {
	return ps.activeProfile != ""
}

/*
 * ---- EXPORTED SINGLETON ----
 */
 var storeInstance *ProfileStore
 
 func GetProfileStore() *ProfileStore {
	if storeInstance == nil {
		storeInstance = &ProfileStore{}
	}
	return storeInstance
 }

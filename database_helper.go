package main

import (
	"database/sql"
	"errors"
)

// WechatURLData represents the data structure for Wechat URLs
type WechatURLData struct {
	ID  int
	URL string
}

// NewDatabaseHelper creates a new instance of DatabaseHelper
func NewDatabaseHelper(db *sql.DB) *DatabaseHelper {
	return &DatabaseHelper{
		db: db,
	}
}

// GetUserByUsername checks if a user exists by username
func (dh *DatabaseHelper) GetUserByUsername(username string) (bool, error) {
	var exists bool
	err := dh.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
	return exists, err
}

// CheckUserCredentials validates user credentials
func (dh *DatabaseHelper) CheckUserCredentials(username, password string) (bool, error) {
	var storedPassword string
	err := dh.db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	// In a real implementation, you would use proper password hashing
	return storedPassword == password, nil
}

// GetLoginPassByVendor retrieves login pass information for a vendor
func (dh *DatabaseHelper) GetLoginPassByVendor(vendorID string) (map[string]string, error) {
	var co string
	err := dh.db.QueryRow("SELECT co FROM login_pass WHERE vendor_id = ?", vendorID).Scan(&co)
	if err != nil {
		if err == sql.ErrNoRows {
			return map[string]string{}, nil
		}
		return nil, err
	}

	return map[string]string{"co": co}, nil
}

// GetAPIURLByCode retrieves API URL by currency code
func (dh *DatabaseHelper) GetAPIURLByCode(code string) (string, error) {
	var url string
	err := dh.db.QueryRow("SELECT url FROM api_urls WHERE currency_code = ?", code).Scan(&url)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("API URL not found")
		}
		return "", err
	}
	return url, nil
}

// GetRandomWechatURL retrieves a random Wechat URL
func (dh *DatabaseHelper) GetRandomWechatURL() (*WechatURLData, error) {
	var urlData WechatURLData
	err := dh.db.QueryRow("SELECT id, url FROM wechat_urls ORDER BY RAND() LIMIT 1").Scan(&urlData.ID, &urlData.URL)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no Wechat URLs available")
		}
		return nil, err
	}
	return &urlData, nil
}

// UpdateWechatURLUseCount increments the use count for a Wechat URL
func (dh *DatabaseHelper) UpdateWechatURLUseCount(id int) error {
	_, err := dh.db.Exec("UPDATE wechat_urls SET use_count = use_count + 1 WHERE id = ?", id)
	return err
}

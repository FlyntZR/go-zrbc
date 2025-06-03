package main

import (
	"go-zrbc/pkg/utils"
	"strings"
	"testing"
)

func TestDismantSID(t *testing.T) {
	// Test case 1: Empty SID
	info, err := utils.DismantSID("", 4)
	if err != nil {
		t.Errorf("Expected no error for empty SID, got %v", err)
	}
	if info != nil {
		t.Errorf("Expected nil info for empty SID, got %+v", info)
	}

	// Test case 2: Valid SID
	sid := "8813C3F6G80001A1688M794132" // Last digit 7 indicates first part length is 7
	info, err = utils.DismantSID(sid, 4)
	if err != nil {
		t.Errorf("Failed to dismant SID: %v", err)
	}
	if info == nil {
		t.Fatal("Expected non-nil info for valid SID")
	}
	t.Logf("info:%+v", info)

	// Check the components
	if len(info.Website) != 4 {
		t.Errorf("Expected website length to be 4, got %d", len(info.Website))
	}
	if info.Utp == "" {
		t.Error("Expected non-empty Utp")
	}
	if info.Ulv == "" {
		t.Error("Expected non-empty Ulv")
	}
	if info.Uid == "" {
		t.Error("Expected non-empty Uid")
	}
	if info.VftR == "" {
		t.Error("Expected non-empty VftR")
	}
}

func TestProSIDCreate(t *testing.T) {
	// Test case 1: Basic SID creation
	wcode := "a168"
	uid := "7941388"
	ulv := "M"
	utp := "8"
	sidlen := 13

	sid := utils.ProSIDCreate(wcode, uid, ulv, utp, sidlen)

	// Test if SID is not empty
	if sid == "" {
		t.Error("Generated SID should not be empty")
	}

	t.Logf("sid:%s", sid)

	// Test if we can dismant the generated SID
	info, err := utils.DismantSID(sid, len(wcode))
	if err != nil {
		t.Errorf("Failed to dismant generated SID: %v", err)
	}
	if info == nil {
		t.Fatal("Expected non-nil info for generated SID")
	}
	t.Logf("info:%+v", info)

	// Verify the components
	if info.Website != strings.ToLower(wcode) {
		t.Errorf("Expected website code %s, got %s", strings.ToLower(wcode), info.Website)
	}
	if info.Uid != uid {
		t.Errorf("Expected uid %s, got %s", uid, info.Uid)
	}
	if info.Ulv != ulv {
		t.Errorf("Expected ulv %s, got %s", ulv, info.Ulv)
	}
	if info.Utp != utp {
		t.Errorf("Expected utp %s, got %s", utp, info.Utp)
	}

	// Test case 2: SID with default length (when sidlen = 0)
	sid2 := utils.ProSIDCreate(wcode, uid, ulv, utp, 0)
	if sid2 == "" {
		t.Error("Generated SID with default length should not be empty")
	}

	info2, err := utils.DismantSID(sid2, len(wcode))
	if err != nil {
		t.Errorf("Failed to dismant generated SID with default length: %v", err)
	}
	if info2 == nil {
		t.Fatal("Expected non-nil info for generated SID with default length")
	}
}

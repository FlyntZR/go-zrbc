package utils

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

// SIDInfo represents the dismantled components of a SID
type SIDInfo struct {
	Website string
	Utp     string
	Ulv     string
	Uid     string
	VftR    string
	VftN    []string
}

// ProSIDCreate generates a SID with specified length
// Parameters:
// wcode: website code
// uid: user ID
// ulv: user level
// utp: user type
// sidlen: SID string length (default 13)
func ProSIDCreate(wcode, ulv, utp string, uid int64, sidlen int) string {
	if sidlen == 0 {
		sidlen = 13
	}

	uidStr := strconv.FormatInt(uid, 10)
	wlen := len(wcode)
	ilen := len(uidStr)
	llen := len(ulv)
	tlen := len(utp)

	// Calculate required SID length
	sidlen = sidlen - (wlen + ilen + llen - 5)

	var strTmp1, strTmp2 strings.Builder
	for i := 0; i < sidlen; i++ {
		rand1 := rand.Intn(9) + 1
		rand2 := rand.Intn(9) + 1
		strTmp1.WriteString(strconv.Itoa(rand1))
		strTmp2.WriteString(strconv.Itoa(rand2))
		if strTmp1.Len()+strTmp2.Len() >= sidlen {
			break
		}
	}

	intTmp1, _ := strconv.ParseFloat(strTmp1.String(), 64)
	intTmp2, _ := strconv.ParseFloat(strTmp2.String(), 64)
	strTmp2Eng := chgASCIIToEng(intTmp2)
	intTmp3 := math.Abs(intTmp1 - intTmp2)
	strTmp3 := fmt.Sprintf("%.0f", intTmp3)

	if len(strTmp3) > 4 {
		strTmp3 = strTmp3[len(strTmp3)-4:]
	}
	if len(strTmp3) < 4 {
		strTmp3 = fmt.Sprintf("%04s", strTmp3)
	}

	intTmp0 := wlen + ilen + llen + tlen

	sid := fmt.Sprintf("%d%s%s%s%s%s%s",
		intTmp0,
		BlendEngNum(fmt.Sprintf("%.0f", intTmp1), strTmp2Eng),
		strTmp3,
		wcode,
		utp,
		ulv,
		uidStr)

	return MixSID(strings.ToUpper(sid))
}

// DismantSID dismantles a SID string into its components
func DismantSID(sid string, wcodeLength int) (*SIDInfo, error) {
	if len(sid) == 0 {
		return nil, nil
	}

	// Get the last digit which indicates how many characters were mixed
	lastChar := sid[len(sid)-1:]
	intTmp0, err := strconv.Atoi(lastChar)
	if err != nil {
		return nil, fmt.Errorf("invalid SID format")
	}

	// Unmix the SID
	strTmp0 := sid[:intTmp0]
	var mixed []string
	for i := len(strTmp0) - 1; i >= 0; i-- {
		mixed = append(mixed, string(strTmp0[i]))
	}
	strTmp1 := sid[intTmp0 : len(sid)-1]
	sid = strTmp1 + strings.Join(mixed, "")

	// Get the length indicator
	intTmp1len := 1
	intTmp1, err := strconv.Atoi(sid[:intTmp1len])
	if err != nil || intTmp1 == 0 {
		return nil, fmt.Errorf("invalid SID format")
	}
	if intTmp1 < 5 {
		intTmp1len = 2
		intTmp1, _ = strconv.Atoi(sid[:intTmp1len])
	}

	// Extract components
	strTmp2 := sid[len(sid)-intTmp1:]
	info := &SIDInfo{
		Website: strings.ToLower(strTmp2[:wcodeLength]),
		Utp:     strTmp2[wcodeLength : wcodeLength+1],
		Ulv:     strTmp2[wcodeLength+1 : wcodeLength+2],
		Uid:     strTmp2[wcodeLength+2:],
	}

	// Extract verification tokens
	sid = sid[intTmp1len : len(sid)-intTmp1]
	info.VftR = sid[len(sid)-4:]
	vftParts := dismantEngNum(sid[:len(sid)-4])
	info.VftN = []string{
		chgEngToASCII(vftParts[0]),
		vftParts[1],
	}

	return info, nil
}

// MixSID mixes up the SID string
func MixSID(sid string) string {
	intrnd := rand.Intn(6) + 2 // 2 to 7
	strTmp := sid[len(sid)-intrnd:]
	sid = sid[:len(sid)-intrnd]

	var mixed []string
	for i := len(strTmp) - 1; i >= 0; i-- {
		mixed = append(mixed, string(strTmp[i]))
	}

	return fmt.Sprintf("%s%s%d", strings.Join(mixed, ""), sid, intrnd)
}

// Helper functions

// dismantEngNum splits mixed English and numbers
func dismantEngNum(engnum string) []string {
	var ary1, ary2 []string
	length := len(engnum) / 2

	for i := 0; i < length; i++ {
		strTmp := engnum[i*2 : i*2+2]
		ary1 = append(ary1, string(strTmp[0]))
		ary2 = append(ary2, string(strTmp[1]))
	}

	return []string{strings.Join(ary1, ""), strings.Join(ary2, "")}
}

// BlendEngNum mixes English and numbers
func BlendEngNum(num, eng string) string {
	var result []string
	for i := 0; i < len(num); i++ {
		str1 := string(num[i])
		str2 := string(eng[i])
		result = append(result, str2+str1)
	}
	return strings.Join(result, "")
}

// chgEngToASCII converts English letters to ASCII numbers
func chgEngToASCII(eng string) string {
	var result []string
	for _, c := range eng {
		result = append(result, strconv.Itoa(int(c)-64))
	}
	return strings.Join(result, "")
}

// chgASCIIToEng converts ASCII numbers to English letters
func chgASCIIToEng(ascii float64) string {
	asciiStr := fmt.Sprintf("%.0f", ascii)
	var result []string
	for _, c := range asciiStr {
		num, _ := strconv.Atoi(string(c))
		result = append(result, string(rune(num+64)))
	}
	return strings.Join(result, "")
}

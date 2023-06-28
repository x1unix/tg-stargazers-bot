package auth

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type UserID = int64

func ParseUserID(str string) (UserID, error) {
	str = strings.TrimSpace(str)
	if str == "" {
		return 0, errors.New("empty user ID")
	}

	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID: %q", err)
	}

	return v, nil
}

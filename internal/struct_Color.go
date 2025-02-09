package deputysheriff

import (
	"errors"
	"slices"
	"strconv"
	"strings"
)

type Color struct {
	hexValue string
	intValue int
}

func (c *Color) colorHexToDecimal(s string) error {
	allowedChars := []string{"A", "B", "C", "D", "E", "F", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}

	if s[0] != '#' {
		return errors.New("missing # in color")
	}
	s = strings.Replace(s, "#", "", -1)
	if len(s) != 3 && len(s) != 6 {
		return errors.New("missing # in color")
	}

	for _, char := range s {
		if !slices.Contains(allowedChars, strings.ToUpper(string(char))) {
			return errors.New("invalid color input")
		}
	}

	if len(s) == 3 {
		s += s // FFF -> FFFFFF
	}
	decimal_num, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return err
	}

	c.hexValue = s
	c.intValue = int(decimal_num)
	return nil
}

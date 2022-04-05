package util

import (
	"fmt"
	"strings"
)

func StockTable(id string) string {
	return fmt.Sprintf("stock_%s", strings.ToLower(id))
}

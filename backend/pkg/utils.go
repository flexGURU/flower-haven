package pkg

import (
	"log"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

func PgTypeNumericToFloat64(n pgtype.Numeric) float64 {
	f, err := n.Float64Value()
	if err != nil {
		log.Println("not float")
		return 0
	}

	return f.Float64
}

func Float64ToPgTypeNumeric(f float64) pgtype.Numeric {
	var amount pgtype.Numeric
	if err := amount.Scan(strconv.FormatFloat(f, 'f', -1, 64)); err != nil {
		log.Println("not float")
		return pgtype.Numeric{
			Valid: false,
		}
	}

	return amount
}

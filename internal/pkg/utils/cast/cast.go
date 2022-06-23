package cast

import (
	"dbms/internal/pkg/utils/log"
	
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"github.com/jackc/pgx/pgtype"
)

func ToString(src []byte) string {
	return string(src)
}

func IntToStr(src uint64) string {
	return fmt.Sprint(src)
}

func FlToStr(src float64) string {
	return fmt.Sprintf("%.1f", src)
}

func TimeToStr(src time.Time, withTime bool) string {
	if withTime {
		return src.Format("2006.01.02 15:04:05")
	}
	return src.Format("2006.01.02")
}

func ToUint64(src []byte) uint64 {
	return binary.BigEndian.Uint64(src)
}

func ToInt(src []byte) int {
	tmp := binary.BigEndian.Uint32(src)
	return int(tmp)
}

func ToFloat64(src []byte) float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(src))
}

func ToTime(src []byte) time.Time {
	tmp := pgtype.Timestamp{}
	err := tmp.DecodeBinary(nil, src)
	if err != nil {
		log.Error(err)
		return time.Time{}
	}
	return tmp.Time
}

func ToDate(src []byte) time.Time {
	tmp := pgtype.Timestamp{}
	err := tmp.DecodeBinary(nil, src)
	if err != nil {
		log.Error(err)
		return time.Time{}
	}
	return tmp.Time
}

func ToBool(src []byte) bool {
	tmp := pgtype.Bool{}
	err := tmp.DecodeBinary(nil, src)
	if err != nil {
		log.Error(err)
		return tmp.Bool
	}
	return tmp.Bool
}

func DateToStringUnderscore(src []byte) (string, error) {
	timeBuffer := pgtype.Date{}
	err := timeBuffer.DecodeBinary(nil, src)
	timeString := timeBuffer.Time.Format("2006-01-02")
	if timeString == "0001-01-01" {
		return "", err
	}
	return timeString, err
}

func ToInt8Arr(src []byte) pgtype.Int8Array {
	tmp := pgtype.Int8Array{}
	err := tmp.DecodeBinary(nil, src)
	if err != nil {
		return pgtype.Int8Array{}
	}
	return tmp
}

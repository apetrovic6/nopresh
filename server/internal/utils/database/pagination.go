package database

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

// type Cursor map[string]any

type SortOrder string

const (
	ASC  SortOrder = "ASC"
	DESC SortOrder = "DESC"
)

type Cursor struct {
	Date      time.Time
	LastId    uint
	SortOrder SortOrder
}

func EncodeCursor(cursor *Cursor) string {
	if cursor == nil {
		return ""
	}

	serializedCursor, err := json.Marshal(cursor)

	if err != nil {
		return ""
	}

	encoded := base64.StdEncoding.EncodeToString(serializedCursor)

	return encoded
}

func DecodeCursor(cursor string) (*Cursor, error) {
	decodedCursor, err := base64.StdEncoding.DecodeString(cursor)

	if err != nil {
		return nil, err
	}

	var cur Cursor

	if err := json.Unmarshal(decodedCursor, &cur); err != nil {
		return nil, err
	}

	return &cur, nil
}

func CreateCursor(id uint, date time.Time, pointsNext bool) *Cursor {

	return &Cursor{
		LastId: id,
		Date:   date,
	}

	// return Cursor{
	// 	"id":            id,
	// 	"date_time_utc": date,
	// 	"points_next":   pointsNext,
	// }
}

// func GetPaginatedQuery(query *gorm.DB, pointsNext bool, cursor, sortOrder string) (*gorm.DB, bool, error) {
// 	if cursor != "" {
// 		decodedCursor, err := DecodeCursor(cursor)

// 		if err != nil {
// 			return nil, pointsNext, err
// 		}

// 		pointsNext = decodedCursor["points_next"] == true
// 		operator, order := getPaginatedOperator(pointsNext, sortOrder)

// 		whereStr := fmt.Sprintf("(date_time_utc %s OR (date_time_utc = ? AND id %s ))", operator, operator)

// 		query = query.Where(whereStr, decodedCursor["date_time_utc"], decodedCursor["date_time_utc"], decodedCursor["id"])

// 		if order != "" {
// 			sortOrder = order
// 		}
// 	}

// 	query = query.Order("date_time_utc " + sortOrder)

// 	return query, pointsNext, nil
// }

// func getPaginatedOperator(pointsNext bool, sortOrder string) (string, string) {
// 	if pointsNext && sortOrder == "asc" {
// 		return ">", ""
// 	}

// 	if pointsNext && sortOrder == "desc" {
// 		return "<", ""
// 	}

// 	if !pointsNext && sortOrder == "asc" {
// 		return "<", "desc"
// 	}

// 	if !pointsNext && sortOrder == "desc" {
// 		return "<", "asc"
// 	}

// 	return "", ""
// }

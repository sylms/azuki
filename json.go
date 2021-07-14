package main

import (
	"database/sql"
	"encoding/json"
)

func (ns *nullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return []byte(`null`), nil
}

func newNullString(s string, valid bool) nullString {
	return nullString{sql.NullString{String: s, Valid: valid}}
}

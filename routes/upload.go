package routes

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"pancast-server/models"
	"pancast-server/server-utils"
	"pancast-server/types"
	"unicode"
)

type UploadParameters struct {
	Entries []types.Entry
	Type    int64
}

func UploadController(input UploadParameters, db *sql.DB) error {
	// input validation
	if !isUploadInputSafe(input.Entries) {
		log.Println("error: unsafe or illegal input")
		return errors.New("unsafe or illegal input")
	}
	if input.Type == server_utils.RISK {
		success := models.CreateRiskEntries(input.Entries, db)
		if !success {
			return errors.New("error: create unsuccessful")
		}
	} else if input.Type == server_utils.EPI {
		success := models.CreateEpiEntries(input.Entries, db)
		if !success {
			return errors.New("error: create unsuccessful")
		}
	} else {
		log.Println("error: database type not valid")
		return errors.New("error: database type not valid")
	}
	return nil
}

func isUploadInputSafe(input []types.Entry) bool {
	for _, entry := range input {
		ephCond := checkInputType(entry.EphemeralID, BASE64)
		// assume that the integers in the struct are safe
		if !ephCond {
			return false
		}
	}
	return true
}

const (
	ALPHABETIC = 0
	NUMERIC    = 1
	BASE64     = 2
	SPECIAL    = 3
)

func checkInputType(input string, inputType int) bool {
	switch inputType {
	case ALPHABETIC:
		for _, c := range input {
			if !unicode.IsLetter(c) {
				return false
			}
		}
	case NUMERIC:
		for _, c := range input {
			if !unicode.IsDigit(c) {
				return false
			}
		}
	case BASE64:
		for _, c := range input {
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) && string(c) != "=" && string(c) != "+" && string(c) != "/" {
				return false
			}
		}
	case SPECIAL:
		for _, c := range input {
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) && string(c) != "_" && string(c) != "-" {
				return false
			}
		}
	default:
		return false
	}
	return true
}

func ConvertStringToUploadParam(input []byte) (UploadParameters, error) {
	var entries UploadParameters
	err := json.Unmarshal(input, &entries)
	if err != nil {
		return UploadParameters{}, err
	}
	return entries, nil
}

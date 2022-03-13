package request

import (
	"encoding/json"
	"fmt"
	"r3/cache"
	"r3/config"
	"r3/db"

	"github.com/jackc/pgtype"
)

func LookupGet(reqJson json.RawMessage, loginId int64) (interface{}, error) {

	var req struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}

	switch req.Name {
	case "access":
		return cache.GetAccessById(loginId)

	case "caption":
		var languageCode string

		if err := db.Pool.QueryRow(db.Ctx, `
			SELECT language_code
			FROM instance.login_setting
			WHERE login_id = $1
		`, loginId).Scan(&languageCode); err != nil {
			return nil, err
		}
		return cache.GetCaptions(languageCode), nil

	case "customizing":
		var res struct {
			CompanyName    string `json:"companyName"`
			CompanyWelcome string `json:"companyWelcome"`
		}
		res.CompanyName = config.GetString("companyName")
		res.CompanyWelcome = config.GetString("companyWelcome")
		return res, nil

	case "feedback":
		return config.GetUint64("repoFeedback"), nil

	case "loginKeys":
		var res struct {
			PrivateEnc       pgtype.Varchar `json:"privateEnc"`
			PrivateEncBackup pgtype.Varchar `json:"privateEncBackup"`
			PublicKey        pgtype.Varchar `json:"publicKey"`
		}

		err := db.Pool.QueryRow(db.Ctx, `
			SELECT key_private_enc, key_private_enc_backup, key_public
			FROM instance.login
			WHERE id = $1
		`, loginId).Scan(&res.PrivateEnc, &res.PrivateEncBackup, &res.PublicKey)

		return res, err
	}
	return nil, fmt.Errorf("unknown lookup name")
}

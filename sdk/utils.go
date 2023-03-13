package sdk

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
)

func EncodeUserEncryptionContext(encryptionContext string) ([]byte, error) {
	ctxMap := map[string]string{}
	err := json.Unmarshal([]byte(encryptionContext), &ctxMap)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("The specified parameter EncryptContext is not valid. %s", err.Error()))
	}
	var keys []string
	for key := range ctxMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var buf bytes.Buffer
	keysLen := len(keys)
	for i, key := range keys {
		buf.WriteString(key)
		buf.WriteByte('=')
		buf.WriteString(ctxMap[key])
		if i < keysLen-1 {
			buf.WriteByte('&')
		}
	}
	encData := sha256.Sum256(buf.Bytes())
	return encData[:], nil
}

package awspresence

import (
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"reflect"
	"regexp"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

// Base64Encode encodes data if the input isn't already encoded using base64.StdEncoding.EncodeToString.
// If the input is already base64 encoded, return the original input unchanged.
func base64Encode(data []byte) string {
	// Check whether the data is already Base64 encoded; don't double-encode
	if isBase64Encoded(data) {
		return string(data)
	}
	// data has not been encoded encode and return
	return base64.StdEncoding.EncodeToString(data)
}

func isBase64Encoded(data []byte) bool {
	_, err := base64.StdEncoding.DecodeString(string(data))
	return err == nil
}

func looksLikeJsonString(s interface{}) bool {
	return regexp.MustCompile(`^\s*{`).MatchString(s.(string))
}

func jsonBytesEqual(b1, b2 []byte) bool {
	var o1 interface{}
	if err := json.Unmarshal(b1, &o1); err != nil {
		return false
	}

	var o2 interface{}
	if err := json.Unmarshal(b2, &o2); err != nil {
		return false
	}

	return reflect.DeepEqual(o1, o2)
}

func isResourceNotFoundError(err error) bool {
	_, ok := err.(*resource.NotFoundError)
	return ok
}

func isResourceTimeoutError(err error) bool {
	timeoutErr, ok := err.(*resource.TimeoutError)
	return ok && timeoutErr.LastError == nil
}

// interfaceStringSlice converts a Terraform interface slice of strings into a slice of string pointers
func interfaceStringSlice(src []interface{}) []*string {
	dst := make([]*string, len(src))
	for i, v := range src {
		vs := v.(string)
		dst[i] = &vs
	}
	return dst
}

func randomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

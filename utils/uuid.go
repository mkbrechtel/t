package utils

import (
	"log"
	uuidv7 "github.com/gofrs/uuid/v5"
	"encoding/base64"
)

func NewUUID () (uid uuidv7.UUID) {
	return uuidv7.Must(uuidv7.NewV7())
}

func ShortEncodeUUID (uid uuidv7.UUID) (sid string) {
	return base64.RawURLEncoding.EncodeToString(uid.Bytes())
}

func LongEncodeUUID (uid uuidv7.UUID) (lid string) {
	return uid.String()
}

func EncodeUUID (uid uuidv7.UUID) (lid string) {
	return ShortEncodeUUID(uid)
}

func DecodeUUID  (encid string) (uid uuidv7.UUID) {
	if (len(encid) == 36) {
		u, err := uuidv7.FromString(encid)
		if err != nil {
			log.Fatalf("DecodeUUID: Failed to parse UUID %q: %v", encid, err)
		}
		return u
	} else {
		b, err := base64.RawURLEncoding.DecodeString(encid)
		if err != nil {
			log.Fatalf("DecodeUUID: Failed to decode base64 %q: %v", encid, err)
		}
		u, err := uuidv7.FromBytes(b)
		if err != nil {
			log.Fatalf("DecodeUUID: Failed to parse UUID %q: %v", encid, err)
		}
		return u
	}
}

// TODO write some tests out of this
// func main() {
// 	newid := NewUUID()
// 	fmt.Println("newid",newid)
// 	shoid := ShortEncodeUUID(newid)
// 	fmt.Println("shoid",shoid)
// 	sdeid := DecodeUUID(shoid)
// 	fmt.Println("sdeid",sdeid)
// 	lonid := LongEncodeUUID(newid)
// 	fmt.Println("lonid",lonid)
// 	ldeid := DecodeUUID(lonid)
// 	fmt.Println("ldeid",ldeid)
// }

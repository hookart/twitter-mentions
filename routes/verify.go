package routes

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt"
	"github.com/hookart/twitter-mentions/models"
	"github.com/spf13/viper"
)

var jwtPubKey *rsa.PublicKey

func InitJWTKey() {
	file, err := os.Open(viper.GetString("key"))
	if err != nil {
		log.Fatal(err)
	}
	stats, _ := file.Stat()

	data := make([]byte, stats.Size())
	_, err = file.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	blockPub, _ := pem.Decode(data)
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	jwtPubKey = genericPublicKey.(*rsa.PublicKey)
}

type VerifyRequest struct {
	JWT string `json:"jwt"`
}

type VerifyResponse struct {
	IsTwitterVerified         bool   `json:"is_twitter_verified"`
	TwitterVerificationString string `json:"twitter_hash"`
	Tweet                     string `json:"tweet"`
	Twitter                   string `json:"twitter"`
}

func Verify(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	request := VerifyRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		log.Println("json parse error", err)
		w.WriteHeader(500)
		return
	}

	log.Println("verifying jwt token", request.JWT)
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(request.JWT, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return jwtPubKey, nil
	})

	if err != nil {
		log.Println("An error occurred parsing the jwt ", err)
		w.WriteHeader(500)
		w.Write([]byte("JWT parse error"))
		return
	}

	if !token.Valid {
		log.Println("An error occured parsing the jwt ", err)
		w.WriteHeader(401)
		w.Write([]byte("Invalid Token"))
		return
	}

	// var ens string
	var address string

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["ens_name"], claims["wallet_public_key"])
		address = (claims["wallet_public_key"]).(string)
		// ens = claims["ens"].(string)
	} else {
		fmt.Println(err)
	}

	db := models.GetDBConnection()

	resp := VerifyResponse{}
	account := models.Account{PublicKey: address}
	db.FirstOrCreate(&account, &models.Account{PublicKey: address})

	if account.Verified {
		resp.IsTwitterVerified = true
		// ok to share because they have authed with the JWT
		resp.Twitter = account.TwitterHandle
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(&resp)
	} else {
		verification := &models.Verification{}
		err := db.First(&verification, models.Verification{AccountID: account.ID})
		if err == nil {
			resp.TwitterVerificationString = verification.VerificationString
		} else {
			newAddr := common.BytesToAddress(
				crypto.Keccak256([]byte(fmt.Sprintf("hook protocol - %s - %f",
					address,
					rand.Float64())))[12:])
			resp.TwitterVerificationString = newAddr.Hex()
			verification := &models.Verification{AccountID: account.ID, VerificationString: resp.TwitterVerificationString}
			db.Create(&verification)
		}
		resp.Tweet = fmt.Sprintf("Verifying my account for gotrekt.xyz (by @HookProtocol) ... %s", resp.TwitterVerificationString)
		json.NewEncoder(w).Encode(&resp)
	}
}

package tools

import (
	"os"
	"io/ioutil"
	"github.com/dvsekhvalnov/jose2go/keys/rsa"
	"log"
	"crypto/rsa"
	"crypto/rand"
	"encoding/pem"
	"crypto/x509"
	"net/http"
	"github.com/dvsekhvalnov/jose2go"
	"encoding/json"
	"time"
	"io"
	"archive/zip"
	"mime/multipart"
)

type (
	JWTSession struct {
		AccessToken		string	`json:"access_token"`
		UserId			uint32	`json:"user_id"`
		ExpiresIn		time.Time `json:"expires_in"`
		Device			string	`json:"device"`
		Username		string	`json:"username"`
		Admin			bool	`json:"admin"`
	}

	JWTSessionAccessRights struct {
		Type	string `json:"type"`
		Id		uint32 `json:"id"`
		Role	string `json:"role"`
	}
)

// ShouldServeStatic
// Checks whether a static file exists or not
func ShouldServeStatic (path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}

	return true
}

func SetJWTCookie(name string, object *JWTSession, writer http.ResponseWriter, publicKey *rsa.PublicKey) error {
	// Encodes the object into a JSON String
	payload, err := json.Marshal(object)
	if err != nil {
		return err
	}

	// Encripts the JSON String into the JWT token
	token, err := jose.Encrypt(string(payload), jose.RSA_OAEP, jose.A256GCM, publicKey)

	// Sets the Cookie
	SetCookie(name, token, writer)
	return nil
}

func SetCookie(name, value string, w http.ResponseWriter) {
	production := GetSettings().Server.Production

	// Generates the Cookie to the Request
	cookie := http.Cookie{
		Name: name,
		Value: value,
		Path: "/",
		HttpOnly: production,
		Secure: production,
	}

	// Adds the Cookie to the request
	http.SetCookie(w, &cookie)
}

func LoadKey(privatePath, publicPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKeyBytes, err := ioutil.ReadFile(privatePath)
	if(err != nil) {
		log.Printf("- File: '%s' is Missing, creating a new one...\n", privatePath)

		// Keys not found. Generating new Keys
		return GenerateKeys(privatePath, publicPath)
	} else {
		// Convert bytes into a public/private key pair to use.
		privateKey, err := Rsa.ReadPrivate(privateKeyBytes)
		if err != nil {
			log.Printf("Error loading private key.")
			return nil, nil, err
		}

		// Reads the Public Key File
		publicKeyBytes, err := ioutil.ReadFile(publicPath)
		if(err != nil) {
			log.Printf("- File: '%s' is Missing, creating a new one...\n", privatePath)

			// Public Key not found, regenerating from private key
			PubASN1, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
			if err != nil {
				log.Printf("Error parsing public key.")
				return nil, nil, err
			}

			// Encode the Public Key into bytes
			publicKeyBytes = pem.EncodeToMemory(&pem.Block{
				Type:  "RSA PUBLIC KEY",
				Bytes: PubASN1,
			})

			// Saves the Public Key
			ioutil.WriteFile(publicPath, publicKeyBytes, 0644)
		}

		// Convert bytes into a public/private key pair to use.
		publicKey, err := Rsa.ReadPublic(publicKeyBytes)
		if err != nil {
			log.Printf("Error loading public key.")
			return nil, nil, err
		}

		return privateKey, publicKey, nil
	}
}

func GenerateKeys(privatePath, publicPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	// Generating a new private Key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	// Encode the Private Key into bytes
	pemData := pem.EncodeToMemory(
		&pem.Block{
			Type: "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	// Extracts the Public Key from the private Key
	PubASN1, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	// Encode the Public Key into bytes
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: PubASN1,
	})

	// Save the Private and Public Keys
	ioutil.WriteFile(privatePath, pemData, 0644)
	ioutil.WriteFile(publicPath, pubBytes, 0644)

	return privateKey, &privateKey.PublicKey, nil
}

func ParseBody(reader io.Reader, obj interface{}) (int, map[string]interface{}) {
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&obj)
	if err != nil {
		return http.StatusForbidden, map[string]interface{}{
			"error":   "InvalidData",
			"message": "The data provided is not valid.",
		}
	}

	return http.StatusOK, nil
}

func FirstDayOfISOWeek(year int, week int, timezone *time.Location) time.Time {
	date := time.Date(year, 0, 0, 0, 0, 0, 0, timezone)
	isoYear, isoWeek := date.ISOWeek()

	// iterate back to Monday
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}

	// iterate forward to the first day of the first week
	for isoYear < year {
		date = date.AddDate(0, 0, 7)
		isoYear, isoWeek = date.ISOWeek()
	}

	// iterate forward to the first day of the given week
	for isoWeek < week {
		date = date.AddDate(0, 0, 7)
		isoYear, isoWeek = date.ISOWeek()
	}

	return date
}

func ZipFiles(zipFilePath string, fileHeaders []*multipart.FileHeader) (*os.File, error) {
	// Creates the ZIP file
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return nil, err
	}
	defer zipFile.Close()

	// Creates a writer to write into the file
	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	for index, fileHeader := range fileHeaders {
		// Get the File
		file, err := fileHeaders[index].Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		// Convert the file header into a valid ZIP header
		header := &zip.FileHeader{
			Name:         fileHeader.Filename,
			Method:       zip.Store,
			ModifiedTime: uint16(time.Now().UnixNano()),
			ModifiedDate: uint16(time.Now().UnixNano()),
		}

		// Set the compression level
		header.Method = zip.Deflate

		// Create a valid Header of the file to add to the file
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return nil, err
		}

		// Copy the file into the zip
		_, err = io.Copy(writer, file)
		if err != nil {
			return nil, err
		}
	}

	return zipFile, nil
}

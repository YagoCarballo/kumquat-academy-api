package tools

import (
	"io/ioutil"
	"os"

	"github.com/naoina/toml"
	"log"
	"strconv"
)

type (
	Settings struct {
		Title       string
		Description string
		Database    Database
		Server      Server
		Email       Email
		Api         Api
	}
	Database struct {
		Type   string
		Mysql  MySQL
		Sqlite SQLite
	}
	MySQL struct {
		Username string
		Password string
		Host     string
		Name     string
	}
	SQLite struct {
		Path	 string
	}
	Server struct {
		Port 		int
		Debug		bool
		Production	bool
		PrivateKey	string
		PublicKey	string
		UploadsPath string
	}
	Email struct {
		Server   string
		Port     int
		User     string
		Password string
		Sender   string
	}
	Api struct {
		Prefix  string
		Version int
	}
)

var localSetting Settings

func createSettings(path string) {
	defaultSettings := Settings{
		Title:       "Kumquat Academy",
		Description: "Kumquat Academy - Learning Platform",
		Database: Database{
			Type:	  "MySQL",
			Mysql:	  MySQL{
				Username: "KumquatAcademy",
				Password: "KumquatAcademy",
				Host:     "localhost:3306",
				Name:     "KumquatAcademy",
			},
			Sqlite:	  SQLite{
				Path:	  "./default.sqlite",
			},
		},
		Server: Server{
			Port:  		3000,
			Debug:		false,
			Production:	false,
			PrivateKey: "./privateKey.pem",
			PublicKey: "./publicKey.pub",
			UploadsPath: "./attachments",
		},
		Email: Email{
			Server: "smtp.gmail.com",
			Port: 	587,
			User: 		"test@gmail.com",
			Password: 	"password",
			Sender: 	"Kumquat Academy <do-not-reply@kumquat.academy>",
		},
		Api: Api{
			Prefix:  "/api",
			Version: 1,
		},
	}

	str, _ := toml.Marshal(defaultSettings)
	ioutil.WriteFile(path, []byte(str), 0644)
}

func LoadSettings(path string) error {
	// Opens the Database File
	f, err := os.Open(path)
	if err != nil {
		log.Printf("- File: '%s' is Missing, creating a new one...\n", path)

		// Create a new settings file with default values
		createSettings(path)

		// Trying to open the new settings
		f, err = os.Open(path)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	defer f.Close()

	// Reads the Contents of the file into a Buffer
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println(err)
		return err
	}

	// Parses the Buffer into the Struct
	if err := toml.Unmarshal(buf, &localSetting); err != nil {
		log.Println(err)
		return err
	}

	// Overrides the settings with Environment Variables
	if os.Getenv("APIPORT") != "" {
		port, err := strconv.Atoi(os.Getenv("APIPORT"))
		if err != nil {
			log.Println("Invalid Port detected in the Environment Variables, Falling back to Settings")
		} else {
			localSetting.Server.Port = port
		}
	}

	if os.Getenv("DB_HOST") != "" {
		localSetting.Database.Mysql.Host = os.Getenv("DB_HOST")
	}

	if os.Getenv("DB_USER") != "" {
		localSetting.Database.Mysql.Username = os.Getenv("DB_USER")
	}

	if os.Getenv("DB_PASS") != "" {
		localSetting.Database.Mysql.Password = os.Getenv("DB_PASS")
	}

	if os.Getenv("DB_NAME") != "" {
		localSetting.Database.Mysql.Name = os.Getenv("DB_NAME")
	}

	if os.Getenv("DB_PATH") != "" {
		localSetting.Database.Sqlite.Path = os.Getenv("DB_PATH")
	}

	if os.Getenv("EMAIL_USER") != "" {
		localSetting.Email.User = os.Getenv("EMAIL_USER")
	}

	if os.Getenv("EMAIL_PASSWORD") != "" {
		localSetting.Email.Password = os.Getenv("EMAIL_PASSWORD")
	}

	if os.Getenv("EMAIL_SENDER") != "" {
		localSetting.Email.Sender = os.Getenv("EMAIL_SENDER")
	}

	if os.Getenv("EMAIL_SMTP_SERVER") != "" {
		localSetting.Email.Server = os.Getenv("EMAIL_SMTP_SERVER")
	}

	if os.Getenv("EMAIL_SMTP_PORT") != "" {
		port, err := strconv.Atoi(os.Getenv("EMAIL_SMTP_PORT"))
		if err != nil {
			log.Println("Invalid Email SMTP Port detected in the Environment Variables, Falling back to Settings")
		} else {
			localSetting.Email.Port = port
		}
	}

	if os.Getenv("UPLOADS_PATH") != "" {
		localSetting.Server.UploadsPath = os.Getenv("UPLOADS_PATH")
	}

	if os.Getenv("PRODUCTION") != "" {
		localSetting.Server.Production = (os.Getenv("PRODUCTION") == "true")
	}

	return nil
}

func (settings *Settings) Save() {
	file, err := os.Create("settings.toml")
	if err != nil {
		log.Println(err)
		return
	}

	defer file.Close()

	tomlEncoder := toml.NewEncoder(file)
	tomlEncoder.Encode(*settings)
}

// Returns the settings
func GetSettings() *Settings {
	return &localSetting
}

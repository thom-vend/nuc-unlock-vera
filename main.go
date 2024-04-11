package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/crypto/scrypt"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Url               string   `yaml:"url"`
	AuthHeader        string   `yaml:"auth_header"`
	AuthToken         string   `yaml:"auth_token"`
	HttpMethod        string   `yaml:"http_method"`
	PayloadPassword   string   `yaml:"payload_password"`
	UnlockCmd         string   `yaml:"unlock_cmd"`
	UnlockArgs        []string `yaml:"unlock_args"`
	UnlockPlaceholder string   `yaml:"unlock_placeholder"`
}

func loadConfig(configPath string) Config {
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	var conf Config
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return conf
}
func httpRequest(method string, url string, authHeader string, token string) (string, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set(authHeader, token)
	req.Header.Set("User-Agent", "NucUnlocker 1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// import code from https://bruinsslot.jp/post/golang-crypto/
func Encrypt(key, data []byte) ([]byte, error) {
	key, salt, err := DeriveKey(key, nil)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	ciphertext = append(ciphertext, salt...)

	return ciphertext, nil
}

func Decrypt(key, data []byte) ([]byte, error) {
	salt, data := data[len(data)-32:], data[:len(data)-32]

	key, _, err := DeriveKey(key, salt)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func DeriveKey(password, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}

	key, err := scrypt.Key(password, salt, 16384, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}

	return key, salt, nil
}

func main() {
	// Read arguments with "flag", load config file path with flag `-c <path>`
	configPath := flag.String("c", "nucunlocker.yml", "path to config file")
	mode := flag.String("m", "unlock", "run mode (unlock/encrypt/decrypt)")
	data := flag.String("d", "", "data to encrypt/decrypt")
	password := flag.String("p", "", "password to encrypt/decrypt (for encrypt/decrypt mode)")
	flag.Parse()
	conf := loadConfig(*configPath)

	switch *mode {
	case "unlock":
		log.Println("Unlocking NUC 🤖")
		// make api call
		response, err := httpRequest(conf.HttpMethod, conf.Url, conf.AuthHeader, conf.AuthToken)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		log.Println("Payload received ✅")

		// decrypt response
		ciphertextbyte, err := base64.StdEncoding.DecodeString(response)
		if err != nil {
			log.Fatal(err)
		}
		plaintextbyte, err := Decrypt([]byte(conf.PayloadPassword), ciphertextbyte)
		if err != nil {
			log.Fatal(err)
		}
		secret := string(plaintextbyte)
		log.Println("Payload decrypted ✅")

		// execute output command
		updatedArgs := make([]string, len(conf.UnlockArgs))
		copy(updatedArgs, conf.UnlockArgs)
		for i, arg := range updatedArgs {
			updatedArgs[i] = strings.ReplaceAll(arg, conf.UnlockPlaceholder, secret)
		}
		log.Println("Command prepared ✅")
		cmd := exec.Command(conf.UnlockCmd, updatedArgs...)
		cmd.Env = os.Environ()
		output, errr := cmd.CombinedOutput()
		if errr != nil {
			log.Fatalf("Error: %v, output: %s", errr, output)
		}
		fmt.Printf("Command output: \n%s\n", output)
		log.Println("NUC unlocked ✅")

	case "encrypt":
		fmt.Println("Encrypting clear text")
		ciphertextbyte, err := Encrypt([]byte(*password), []byte(*data))
		if err != nil {
			log.Fatal(err)
		}
		base64ciphertext := base64.StdEncoding.EncodeToString(ciphertextbyte)
		fmt.Printf("Encrypted data: \n----------------\n%s\n----------------\n", string(base64ciphertext))
	case "decrypt":
		fmt.Println("Decrypting cipher text")
		ciphertextbyte, err := base64.StdEncoding.DecodeString(*data)
		if err != nil {
			log.Fatal(err)
		}
		plaintext, err := Decrypt([]byte(*password), ciphertextbyte)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Decrypted data: \n----------------\n%s\n----------------\n", string(plaintext))
	default:
		log.Fatalf("Invalid mode")
		os.Exit(1)
	}
}

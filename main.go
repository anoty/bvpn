package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Ovpn     string
	Cfg      string
	Pass     string
	Username string
	Password string
	Secret   string
}

var conf Config

func toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func toUint32(bytes []byte) uint32 {
	return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) +
		(uint32(bytes[2]) << 8) + uint32(bytes[3])
}

func oneTimePassword(key []byte, value []byte) uint32 {
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(value)
	hash := hmacSha1.Sum(nil)

	offset := hash[len(hash)-1] & 0x0F

	hashParts := hash[offset : offset+4]
	hashParts[0] = hashParts[0] & 0x7F
	number := toUint32(hashParts)

	pwd := number % 1000000
	return pwd
}

func googleAuthenticatorCode() (s string) {
	inputNoSpaces := strings.Replace(conf.Secret, " ", "", -1)
	inputNoSpacesUpper := strings.ToUpper(inputNoSpaces)
	key, err := base32.StdEncoding.DecodeString(inputNoSpacesUpper)
	if err != nil {
		log.Fatal(err)
	}
	epochSeconds := time.Now().Unix()
	pwd := oneTimePassword(key, toBytes(epochSeconds/30))
	s = fmt.Sprintf("%06d", pwd)
	return
}

func resetPassFile(code string) {
	pass := conf.Pass
	f, err := os.Create(pass)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	f.WriteString(conf.Username + "\r\n" + conf.Password + code + "\r\n")
}

func connect() {
	argv := []string{"--config", conf.Cfg, "--connect-retry", "0", "--auth-user-pass", conf.Pass}
	cmd := exec.Command(conf.Ovpn, argv[0:]...)
	o, _ := cmd.Output()
	fmt.Println(string(o))
	cmd.Start()
	cmd.Wait()
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var confPath string
	flag.StringVar(&confPath, "conf", "", "config path")
	flag.Parse()
	if confPath == "" {
		log.Fatal("conf err")
	}
	if _, err := toml.DecodeFile(confPath, &conf); err != nil {
		log.Fatal(err)
	}
	for {
		resetPassFile(googleAuthenticatorCode())
		log.Println("pass ok")
		log.Println("bvpn ok")
		connect()
		log.Println("bvpn retry")
	}
}

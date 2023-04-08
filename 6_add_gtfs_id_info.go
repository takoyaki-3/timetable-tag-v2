package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	json "github.com/takoyaki-3/go-json"
)

type InfoType struct {
	VersionID           string `json:"version_id"`
	ByStopHashValueSize int    `json:"by_stop_hash_value_size"`
	ByTripHashValueSize int    `json:"by_trip_hash_value_size"`
}

func main() {
	// RSA秘密鍵を読み込む
	privateKeyBytes, err := ioutil.ReadFile("key.pem")
	if err != nil {
		log.Fatalln(err)
	}

	// Set the root directory path
	root := "dist"

	// Get the list of directories in the root directory
	dirs, err := ioutil.ReadDir(root)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Loop through the directories in the root directory
	for _, dir := range dirs {
		// Check if the item is a directory
		if dir.IsDir() {
			gtfsID := dir.Name()

			// Get the list of subdirectories in the directory
			subdirs, err := ioutil.ReadDir(filepath.Join(root, dir.Name()))
			if err != nil {
				fmt.Println(err)
				continue
			}

			versionInfo := []InfoType{}

			// Loop through the subdirectories in the directory
			for _, subdir := range subdirs {
				// Check if the item is a directory
				if subdir.IsDir() {
					// Print the subdirectory name and full path
					dirPath := filepath.Join(root, dir.Name(), subdir.Name())
					fmt.Printf("  Found subdirectory: %s (%s) %s\n", subdir.Name(), gtfsID, dirPath)

					// JSON Load
					info := InfoType{}
					json.LoadFromPath(dirPath+"/info.json", &info)
					info.VersionID = subdir.Name()
					versionInfo = append(versionInfo, info)
					fmt.Println(info)
				}
			}

			json.DumpToFile(versionInfo, "./dist/"+dir.Name()+"/info.json")
			AddSing("./dist/"+dir.Name()+"/info.json", privateKeyBytes)
		}
	}
}

// 電子署名用
// 電子署名を施す
func Sing(dataBytes, privateKeyBytes []byte) ([]byte, error) {
	privateKeyBlock, _ := pem.Decode(privateKeyBytes)

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	// ファイルのハッシュ値を計算する
	hash := sha256.Sum256(dataBytes)

	// ハッシュ値に署名する
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return nil, err
	}

	return signature, nil
}

func AddSing(path string, privateKeyBytes []byte) error {
	// ファイルを読み込む
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	signature, err := Sing(file, privateKeyBytes)
	if err != nil {
		return err
	}

	// 署名をファイルに書き込む
	err = ioutil.WriteFile(path+".sig", signature, 0644)
	return err
}

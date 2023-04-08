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
	"time"

	json "github.com/takoyaki-3/go-json"
)

type DataItem struct {
	GtfsID             string `json:"gtfs_id"`
	AgencyID           string `json:"agency_id"`
	Name               string `json:"name"`
	License            string `json:"license"`
	LicenseURL         string `json:"licenseUrl"`
	ProviderURL        string `json:"providerUrl"`
	ProviderName       string `json:"providerName"`
	ProviderAgencyName string `json:"providerAgencyName"`
	Memo               string `json:"memo"`
	UpdatedAt          string `json:"updatedAt"`
}

type DataList struct {
	Data []DataItem `json:"data_list"`
}

type Data []struct {
	AgentName          string `json:"事業者名,omitempty"`
	URL                string `json:"事業者名_url,omitempty"`
	Prefecture         string `json:"都道府県,omitempty"`
	GTFS               string `json:"GTFSフィード名,omitempty"`
	License            string `json:"ライセンス,omitempty"`
	LicenseURL         string `json:"ライセンス_url,omitempty"`
	URLs               string `json:"URLs,omitempty"`
	GTFSURL            string `json:"GTFS_url,omitempty"`
	StartDate          string `json:"最新GTFS開始日,omitempty"`
	EndDate            string `json:"最新GTFS終了日,omitempty"`
	UpdateDate         string `json:"最終更新日,omitempty"`
	Detail             string `json:"詳細,omitempty"`
	GtfsID             string `json:"gtfs_id,omitempty"`
	AlertURL           string `json:"Alert_url,omitempty"`
	TripUpdateURL      string `json:"TripUpdate_url,omitempty"`
	VehiclePositionURL string `json:"VehiclePosition_url,omitempty"`
}

func main() {
	// RSA秘密鍵を読み込む
	privateKeyBytes, err := ioutil.ReadFile("key.pem")
	if err != nil {
		log.Fatalln(err)
	}

	//
	data := Data{}
	json.LoadFromPath("data.json", &data)

	datalist := DataList{}

	for _, v := range data {
		t, err := time.Parse("2006-01-02", v.UpdateDate)
		if err != nil {
			continue
			log.Fatalln(err)
		}

		datalist.Data = append(datalist.Data, DataItem{
			GtfsID: v.GtfsID,
			// AgencyID: v.,
			Name:               v.AgentName,
			License:            v.License,
			LicenseURL:         v.LicenseURL,
			ProviderURL:        v.GTFSURL,
			ProviderName:       "",
			ProviderAgencyName: "",
			UpdatedAt:          t.Format("2006-01-02T15_04_05+09_00"),
		})
	}

	json.DumpToFile(datalist, "dist/datalist.json")
	err = AddSing("dist/datalist.json", privateKeyBytes)
	fmt.Println(err)
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

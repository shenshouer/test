package main

import(

	"crypto/aes"
	"crypto/cipher"
	"bytes"
	//"encoding/hex"
	"log"
	"fmt"
	//"encoding/hex"
	"encoding/hex"
	"path/filepath"
	"os"
	"regexp"
	"io/ioutil"
	"encoding/json"
)

func main() {
	log.SetFlags(log.Flags()|log.Lshortfile)

	//origData := []byte(`{"agents":"RB250AM539D7PJ,RBF50FM53H0UBR","enname":"GS05","serial":"eeeeeee"}`)
	//
	//key := []byte("1qaz2wsx3edc4rfv");
	//fmt.Println(len(key))
	//
	//cry, err := AesEncrypt(origData, key)
	//if err !=nil{
	//	log.Fatal(err)
	//}
	//
	//fmt.Println(string(cry[:]))
	//data, err := AesDecrypt(cry, key)
	//if err != nil{
	//	log.Fatal(err)
	//}
	//
	//fmt.Println(string(data[:]))



	src := "/Users/sope/workspaces/go/src/test"
	validate(src)
}

type ValidateInfo struct {
	Agents			string 		`json:"agents"`
	Enname			string 		`json:"enname"`
	Generate_time	string 		`json:"generate_time"`
	Period			string		`json:"period"`
	Serial			string 		`json:"serial"`
}

func validate(src string)(ok bool, err error){
	var r3 = regexp.MustCompile(`bjbus-udisk$`)
	key := []byte("1qaz2wsx3edc4rfv");
	vi :=  &ValidateInfo{}

	filepath.Walk(src, func(path string, info os.FileInfo, err error)error{
		if !info.IsDir() && r3.MatchString(path){
			if fileData, err := ioutil.ReadFile(path); err != nil{
				log.Fatal(err)
			}else{
				tt, err := hex.DecodeString(string(fileData[:]))
				if err != nil{
					log.Fatal(err)
				}

				data, err := AesDecrypt(tt, key)
				if err != nil{
					log.Fatal(err)
				}

				fmt.Println(string(data[:]))

				if err := json.Unmarshal(data, vi); err != nil{
					log.Fatal(err)
				}
			}
		}
		return nil
	})

	log.Println(vi)
	return
}

func AesEncrypt(origData, key []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil

}



func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil

}


func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}


func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}


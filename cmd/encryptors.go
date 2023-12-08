package main

import (
	"bytes"
	"errors"
	"file-server-v3/pkg/mars"
	"file-server-v3/pkg/ntru"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"io"
	"log"
	"os"
	"path/filepath"
)

func DecryptMars(file *os.File, name string, key []byte, progressBar *widget.ProgressBar) error {
	out, err := os.OpenFile(filepath.Join(DecFilesPath, name), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return err
	}
	defer out.Close()

	var buff bytes.Buffer

	_, err = io.Copy(&buff, file)
	if err != nil {
		log.Println(err)
		return err
	}

	data := buff.Bytes()

	decContent := mars.Decrypt(data, key, progressBar)

	d := bytes.NewReader(decContent)

	_, err = io.Copy(out, d)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func EncryptMars(read fyne.URIReadCloser, key []byte, progressBar *widget.ProgressBar) error {
	if read == nil {
		return errors.New("file not chosen")
	}
	file, err := os.Open(read.URI().Path())
	if file == nil || err != nil {
		log.Println("file not chosen")
	}
	defer file.Close()

	encClients, err := os.ReadDir(EncFilesPath)
	if err != nil {
		log.Println(err)
		return err
	}

	clientAlreadyExists := false
	for _, v := range encClients {
		if v.Name() == clientMars.Name {
			clientAlreadyExists = true
		}
	}

	if !clientAlreadyExists {
		err = os.Mkdir(filepath.Join(EncFilesPath, clientMars.Name), 0750)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	out, err := os.OpenFile(filepath.Join(EncFilesPath, clientMars.Name, read.URI().Name()), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return err
	}
	defer out.Close()

	var buff bytes.Buffer

	_, err = io.Copy(&buff, file)
	if err != nil {
		log.Println(err)
		return err
	}

	data := buff.Bytes()

	encContent := mars.Encrypt(data, key, progressBar)

	r := bytes.NewReader(encContent)

	_, err = io.Copy(out, r)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func DecryptNtrue(file *os.File, name string, privateKey *ntru.PrivateKey, partSize int64, progressBar *widget.ProgressBar, win fyne.Window) error {
	fileInfo, err := file.Stat()
	if err != nil {
		log.Println(err)
		return err
	}
	fileSize := fileInfo.Size()
	fmt.Println(fileSize)

	parts := fileSize / partSize
	remainder := fileSize % partSize

	out, err := os.OpenFile(filepath.Join(DecFilesPath, name), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return err
	}
	defer out.Close()

	var buff bytes.Buffer

	for i := int64(0); i < parts; i++ {
		buff.Reset()
		_, err = io.CopyN(&buff, file, partSize)
		if err != nil {
			return err
		}

		data := buff.Bytes()

		decContent, err := ntru.Decrypt(privateKey, data)
		if err != nil {
			log.Println(err)
			return err
		}

		d := bytes.NewReader(decContent)

		_, err = io.Copy(out, d)
		if err != nil {
			log.Println(err)
			return nil
		}

		progressBar.SetValue(float64(i+1) / (float64(parts) + float64(2)))

	}
	progressBar.SetValue(1)
	if remainder > 0 {
		buff.Reset()
		_, err = io.CopyN(&buff, file, remainder)
		if err != nil {
			return err
		}

		data := buff.Bytes()

		decContent, err := ntru.Decrypt(privateKey, data)
		if err != nil {
			log.Println(err)
			return err
		}

		d := bytes.NewReader(decContent)

		_, err = io.Copy(out, d)
		if err != nil {
			log.Println(err)
			return nil
		}

		progressBar.SetValue(1)
	}

	return nil
}

func EncryptNtrue(read fyne.URIReadCloser, publicKey *ntru.PublicKey, partSize int64, progressBar *widget.ProgressBar, win fyne.Window) error {
	file, err := os.Open(read.URI().Path())
	if err != nil {
		return err
	}
	defer file.Close()

	encClients, err := os.ReadDir(EncFilesPath)
	if err != nil {
		log.Println(err)
		return err
	}

	clientAlreadyExists := false
	for _, v := range encClients {
		if v.Name() == clientNtrue.Name {
			clientAlreadyExists = true
		}
	}

	if !clientAlreadyExists {
		err = os.Mkdir(filepath.Join(EncFilesPath, clientNtrue.Name), 0750)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	parts := fileSize / partSize
	remainder := fileSize % partSize

	out, err := os.OpenFile(filepath.Join(EncFilesPath, clientNtrue.Name, read.URI().Name()), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return err
	}
	defer out.Close()

	var buff bytes.Buffer

	for i := int64(0); i < parts; i++ {
		buff.Reset()
		_, err = io.CopyN(&buff, file, partSize)
		if err != nil {
			return err
		}

		data := buff.Bytes()

		encContent, err := ntru.Encrypt(reader, publicKey, data)
		if err != nil {
			log.Println(err)
			return err
		}

		r := bytes.NewReader(encContent)

		_, err = io.Copy(out, r)
		if err != nil {
			log.Println(err)
			return nil
		}

		progressBar.SetValue(float64(i+1) / (float64(parts) + float64(2)))

	}

	if remainder > 0 {
		buff.Reset()
		_, err = io.CopyN(&buff, file, remainder)
		if err != nil {
			return err
		}

		data := buff.Bytes()

		log.Println(len(data))

		encContent, err := ntru.Encrypt(reader, publicKey, data)
		if err != nil {
			log.Println(err)
			return err
		}

		r := bytes.NewReader(encContent)

		_, err = io.Copy(out, r)
		if err != nil {
			log.Println(err)
			return nil
		}

		progressBar.SetValue(1)
	}

	return nil
}

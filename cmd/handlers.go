package main

import (
	"bufio"
	"file-server-v3/pkg/mars"
	"file-server-v3/pkg/ntru"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var clients []string
var files []string
var directories []string
var file File

func RefreshClients(selectReceiver *widget.Select, selectSender *widget.Select, currentClient string) {
	clientsFromServer, err := os.ReadDir(ClientsPath)
	if err != nil {
		log.Println(err)
	}

	dataFromServer, err := os.ReadDir(EncFilesPath)
	if err != nil {
		log.Println(err)
	}

	for _, v := range clientsFromServer {
		if (!slices.Contains(clients, v.Name())) && (v.Name() != currentClient) {
			clients = append(clients, v.Name())
		}

	}

	for _, v := range dataFromServer {
		if !slices.Contains(directories, v.Name()) && (v.Name() != clientMars.Name) {
			directories = append(directories, v.Name())
		}
	}
	selectSender.SetOptions(directories)

	selectReceiver.SetOptions(clients)
}

func RefreshFiles(selectFile *widget.Select, algo *string) {
	senderDirMars, err := os.ReadDir(filepath.Join(EncFilesPath, directoryMars.Name))
	senderDirNtrue, err := os.ReadDir(filepath.Join(EncFilesPath, directoryNtrue.Name))

	if err != nil {
		log.Println(err)
	}

	switch *algo {
	case "NTRUEncrypt":
		files = nil
		for _, v := range senderDirNtrue {
			if !slices.Contains(files, v.Name()) && strings.Contains(v.Name(), ".") {
				files = append(files, v.Name())
			}
		}
		selectFile.SetOptions(files)

	case "MARS":
		files = nil
		for _, v := range senderDirMars {
			if !slices.Contains(files, v.Name()) && strings.Contains(v.Name(), ".") {
				files = append(files, v.Name())
			}
		}
		selectFile.SetOptions(files)
	}

}

func UploadFileMars(win fyne.Window, name string, progressBar *widget.ProgressBar) {
	dialog.ShowFileOpen(func(read fyne.URIReadCloser, err error) {
		nameOfReceiver := name

		inFile, err := os.OpenFile(filepath.Join(ClientsPath, nameOfReceiver), os.O_RDONLY, 0)
		if err != nil {
			log.Println(err)
		}
		defer inFile.Close()

		publicKeyReceiver, err := io.ReadAll(inFile)
		if err != nil {
			log.Println(err)
			return
		}

		inP, err := os.OpenFile(filepath.Join(ConstantsPath, "p"), os.O_RDONLY, 0)
		if err != nil {
			log.Println(err)
		}
		defer inP.Close()

		p, err := io.ReadAll(inP)
		if err != nil {
			log.Println(err)
			return
		}

		key := mars.GetCommonSecretKey(publicKeyReceiver, clientMars.PrivateKey, p)

		err = EncryptMars(read, key, progressBar)
		if err != nil {
			log.Println(err)
			return
		}
		dialog.ShowInformation("Успех", fmt.Sprintf("Файл %s был успешно зашифрован и отправлен на сервер", read.URI().Name()), win)
	}, win)
}

func DownloadFileMars(win fyne.Window, name string, progressBar *widget.ProgressBar) {
	senderPublicKey, err := os.Open(filepath.Join(ClientsPath, directoryMars.Name))
	if err != nil {
		log.Println(err)
		return
	}
	defer senderPublicKey.Close()

	p, err := os.Open(filepath.Join(ConstantsPath, "p"))
	if err != nil {
		log.Println(err)
		return
	}
	defer p.Close()

	publicKeyReceiver := bufio.NewReader(senderPublicKey)

	publicKey, _ := io.ReadAll(publicKeyReceiver)

	pReceiver := bufio.NewReader(p)

	P, _ := io.ReadAll(pReceiver)

	key := mars.GetCommonSecretKey(publicKey, clientMars.PrivateKey, P)

	nameOfFile := name

	file, err := os.OpenFile(filepath.Join(EncFilesPath, directoryMars.Name, nameOfFile), os.O_RDONLY, 0)

	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	err = DecryptMars(file, nameOfFile, key, progressBar)
	if err != nil {
		log.Print(err)
		return
	}
	dialog.ShowInformation("Успех", fmt.Sprintf("Файл %s был успешно зашифрован и отправлен на сервер", name), win)
}

func UploadFileNtrue(win fyne.Window, name string, progressBar *widget.ProgressBar) {
	dialog.ShowFileOpen(func(read fyne.URIReadCloser, err error) {

		var partSize int64
		partSize = 247

		nameOfReceiver := name

		in, err := os.OpenFile(filepath.Join(ClientsPath, nameOfReceiver), os.O_RDONLY, 0)
		if err != nil {
			log.Println(err)
		}

		pubKeyReceiver, err := io.ReadAll(in)
		if err != nil {
			log.Println(err)
		}

		newPublicKey, err := ntru.NewPublicKey(pubKeyReceiver)
		if err != nil {
			log.Println(err)
			return
		}

		err = EncryptNtrue(read, newPublicKey, partSize, progressBar, win)
		if err != nil {
			log.Println(err)
			return
		}
		dialog.ShowInformation("Успех", fmt.Sprintf("Файл %s был успешно зашифрован и отправлен на сервер", read.URI().Name()), win)

	}, win)
}

func DownloadFileNtrue(win fyne.Window, name string, progressBar *widget.ProgressBar) {
	var partSize int64
	partSize = 2062

	nameOfFile := name

	file, err := os.OpenFile(filepath.Join(EncFilesPath, directoryNtrue.Name, nameOfFile), os.O_RDONLY, 0)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	privateKey := clientNtrue.PrivateKey

	err = DecryptNtrue(file, nameOfFile, privateKey, partSize, progressBar, win)
	if err != nil {
		log.Println(err)
		return
	}

	dialog.ShowInformation("Успех", fmt.Sprintf("Файл %s был успешно расшифрован и сохранен на устройстве", name), win)

}

package main

import (
	"bufio"
	"crypto/rand"
	"file-server-v3/pkg/mars"
	"file-server-v3/pkg/ntru"
	"file-server-v3/pkg/ntru/params"
	"flag"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	ClientsPath   = "internal/clients"
	EncFilesPath  = "internal/files/encrypted"
	DecFilesPath  = "internal/files/decrypted"
	ConstantsPath = "internal/constants"
)

type ClientNtrue struct {
	Name         string
	Online       bool
	ReceiverName string
	PrivateKey   *ntru.PrivateKey
	PublicKey    *ntru.PublicKey
}

type ClientMars struct {
	Name         string
	Online       bool
	ReceiverName string
	PublicKey    []byte
	PrivateKey   []byte
}

type File struct {
	Name string
}

type Directory struct {
	Name string
}

var clientNtrue ClientNtrue
var clientMars ClientMars
var directoryMars Directory
var directoryNtrue Directory
var reader = rand.Reader

type UI struct {
	UploadBtn      *widget.Button
	DownloadBtn    *widget.Button
	SubmitBtn      *widget.Button
	RefreshBtn     *widget.Button
	NameEntry      *widget.Entry
	ProgressBar    *widget.ProgressBar
	SelectReceiver *widget.Select
	SelectSender   *widget.Select
	SelectFile     *widget.Select
}

func MakeUI(win fyne.Window, algo *string) *UI {
	nameEntry := widget.NewEntry()
	nameEntry.Resize(fyne.NewSize(200, 40))
	nameEntry.Move(fyne.NewPos(30, 30))
	nameEntry.SetPlaceHolder("введите имя...")

	selectReceiver := widget.NewSelect(
		clients,
		func(s string) {
			clientNtrue.ReceiverName = s
			clientMars.ReceiverName = s
		},
	)
	selectReceiver.PlaceHolder = "Выберите имя получателя..."
	selectReceiver.Resize(fyne.NewSize(200, 40))
	selectReceiver.Move(fyne.NewPos(30, 100))

	selectSender := widget.NewSelect(
		directories,
		func(s string) {
			switch *algo {
			case "NTRUEncrypt":
				directoryNtrue.Name = s
			case "MARS":
				directoryMars.Name = s
			}

		},
	)
	selectSender.PlaceHolder = "Выберите имя отправителя..."
	selectSender.Resize(fyne.NewSize(200, 40))
	selectSender.Move(fyne.NewPos(470, 100))

	var selectFile *widget.Select
	switch *algo {
	case "NTRUEncrypt":
		selectFile = widget.NewSelect(
			files,
			func(s string) {
				file.Name = s
			},
		)
		selectFile.PlaceHolder = "Выберите файл..."
		selectFile.Resize(fyne.NewSize(200, 40))
		selectFile.Move(fyne.NewPos(250, 100))
	case "MARS":
		selectFile = widget.NewSelect(
			files,
			func(s string) {
				file.Name = s
			},
		)
		selectFile.PlaceHolder = "Выберите файл..."
		selectFile.Resize(fyne.NewSize(200, 40))
		selectFile.Move(fyne.NewPos(250, 100))

	}

	submitBtn := widget.NewButton("Отправить", func() {
		var nameOfClient = ""
		nameOfClient = nameEntry.Text

		if len(nameOfClient) == 0 {
			log.Println("Empty name of client")
			return
		}

		switch *algo {
		case "NTRUEncrypt":
			createNewNtrueClient(nameOfClient)
		case "MARS":
			createNewMarsClient(nameOfClient)
		}
		RefreshClients(selectReceiver, selectSender, nameEntry.Text)
		RefreshFiles(selectFile, algo)

		dialog.ShowInformation("Подтверждение", fmt.Sprintf("Добро пожаловать, %s. Выбран алгоритм шифрования %s", nameOfClient, *algo), win)
		nameEntry.Disable()
	})
	submitBtn.Resize(fyne.NewSize(100, 40))
	submitBtn.Move(fyne.NewPos(250, 30))

	refreshBtn := widget.NewButton("Обновить", func() {
		RefreshClients(selectReceiver, selectSender, nameEntry.Text)
		RefreshFiles(selectFile, algo)
	})
	refreshBtn.Resize(fyne.NewSize(70, 40))
	refreshBtn.Move(fyne.NewPos(380, 30))

	progressBar := widget.NewProgressBar()
	progressBar.Resize(fyne.NewSize(640, 30))
	progressBar.Move(fyne.NewPos(30, 240))

	uploadBtn := widget.NewButton("Загрузить файл", func() {

		if len(nameEntry.Text) != 0 {
			submitBtn.Disable()
		}
		switch *algo {
		case "NTRUEncrypt":
			UploadFileNtrue(win, clientNtrue.ReceiverName, progressBar, selectSender, selectReceiver, selectFile, refreshBtn)
		case "MARS":
			UploadFileMars(win, clientMars.ReceiverName, progressBar)
		}

	})
	uploadBtn.Resize(fyne.NewSize(200, 40))
	uploadBtn.Move(fyne.NewPos(30, 170))

	downloadBtn := widget.NewButton("Скачать файл", func() {
		selectSender.Disable()
		selectReceiver.Disable()
		selectFile.Disable()
		refreshBtn.Disable()
		uploadBtn.Disable()

		if len(nameEntry.Text) != 0 {
			submitBtn.Disable()
		}
		switch *algo {
		case "NTRUEncrypt":
			DownloadFileNtrue(win, file.Name, progressBar)
		case "MARS":
			DownloadFileMars(win, file.Name, progressBar)
		}

		selectSender.Enable()
		selectReceiver.Enable()
		selectFile.Enable()
		refreshBtn.Enable()
		uploadBtn.Enable()

	})
	downloadBtn.Resize(fyne.NewSize(200, 40))
	downloadBtn.Move(fyne.NewPos(250, 170))

	return &UI{
		UploadBtn:      uploadBtn,
		DownloadBtn:    downloadBtn,
		SubmitBtn:      submitBtn,
		NameEntry:      nameEntry,
		ProgressBar:    progressBar,
		SelectReceiver: selectReceiver,
		SelectSender:   selectSender,
		SelectFile:     selectFile,
		RefreshBtn:     refreshBtn,
	}
}
func createNewNtrueClient(name string) {
	file, err := os.Create(filepath.Join(ClientsPath, name))
	if err != nil {
		log.Println("error creating user")
	}
	defer file.Close()
	file.Sync()

	keypair, err := ntru.GenerateKey(reader, params.EES1499EP1)
	if err != nil {
		log.Println(err)
		return
	}
	publicKey := keypair.PublicKey.Bytes()
	clientNtrue.PublicKey = &keypair.PublicKey
	clientNtrue.PrivateKey = keypair
	clientNtrue.Name = name
	clientNtrue.Online = true

	fmt.Println("client_name: ", name)

	_, err = file.Write(publicKey)
	if err != nil {
		log.Println("error writing key to the client's file")
	}
}

func createNewMarsClient(name string) {
	nameOfClient := name
	constants, err := os.ReadDir(ConstantsPath)
	if err != nil {
		log.Println(err)
		return
	}

	if len(constants) == 0 {
		generatePG()
	}

	p, err := os.OpenFile(filepath.Join(ConstantsPath, "p"), os.O_RDONLY, 0)
	if err != nil {
		log.Println(err)
		return
	}
	defer p.Close()

	g, err := os.OpenFile(filepath.Join(ConstantsPath, "g"), os.O_RDONLY, 0)
	if err != nil {
		log.Println(err)
		return
	}
	defer g.Close()

	pReceiver := bufio.NewReader(p)
	gReceiver := bufio.NewReader(g)

	P, err := io.ReadAll(pReceiver)
	if err != nil {
		log.Println(err)
		return
	}

	G, err := io.ReadAll(gReceiver)
	if err != nil {
		log.Println(err)
		return
	}
	publicKey, privateKey := mars.GenerateKeyPair(P, G)

	clientMars.PublicKey = publicKey
	clientMars.PrivateKey = privateKey
	clientMars.Name = nameOfClient
	clientMars.Online = true

	client, err := os.OpenFile(filepath.Join(ClientsPath, clientMars.Name), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return
	}
	defer client.Close()

	_, err = client.Write(clientMars.PublicKey)
	if err != nil {
		log.Println(err)
		return
	}
}

func generatePG() {
	p, err := os.Create(filepath.Join(ConstantsPath, "p"))
	if err != nil {
		log.Println(err)
		return
	}
	defer p.Close()

	g, err := os.Create(filepath.Join(ConstantsPath, "g"))
	if err != nil {
		log.Println(err)
		return
	}
	defer g.Close()

	P, G := mars.GetCommonPrimeNumbers()

	_, err = p.Write(P)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = g.Write(G)
	if err != nil {
		log.Println(err)
		return
	}
}

func runUI(algo *string) error {
	a := app.New()

	win := a.NewWindow("FileServer")
	win.Resize(fyne.NewSize(700, 300))
	win.CenterOnScreen()

	ui := MakeUI(win, algo)

	cont := container.NewWithoutLayout(ui.UploadBtn, ui.DownloadBtn, ui.SubmitBtn, ui.NameEntry, ui.ProgressBar, ui.SelectReceiver, ui.SelectFile, ui.SelectSender, ui.RefreshBtn)

	win.SetContent(cont)

	win.ShowAndRun()

	return nil
}

func riskyOperation() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered: ", r)
		}
	}()
	panic("ok")
}

func main() {
	defer riskyOperation()
	algo := flag.String("Algorithm", "MARS", "Chosen algorithm for encryption")
	flag.Parse()

	_ = remove(ClientsPath)
	_ = remove(EncFilesPath)
	_ = remove(ConstantsPath)

	if err := runUI(algo); err != nil {
		panic(err)
	}
}

func remove(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

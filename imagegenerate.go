package main

import (
    "bufio"
    "bytes"
    "fmt"
    ico "github.com/biessek/golang-ico"
    findfont "github.com/flopp/go-findfont"
    "github.com/fogleman/gg"
    "syscall"
    "unsafe"
)

// Funktion zur Erkennung des Windows-Themes
func isLightTheme() bool {
    var hKey syscall.Handle
    lightTheme := uint32(1)

    err := syscall.RegOpenKeyEx(syscall.HKEY_CURRENT_USER, syscall.StringToUTF16Ptr(`Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`), 0, syscall.KEY_READ, &hKey)
    if err != nil {
        return true // Standardmäßig auf helles Thema setzen, falls ein Fehler auftritt
    }
    defer syscall.RegCloseKey(hKey)

    var buf [4]byte
    var bufLen uint32 = 4
    err = syscall.RegQueryValueEx(hKey, syscall.StringToUTF16Ptr("AppsUseLightTheme"), nil, nil, (*byte)(unsafe.Pointer(&buf[0])), &bufLen)
    if err != nil {
        return true // Standardmäßig auf helles Thema setzen, falls ein Fehler auftritt
    }

    lightTheme = *(*uint32)(unsafe.Pointer(&buf[0]))
    return lightTheme == 1
}

// Modifizierte Funktion zur Generierung des Tray-Icons
func generateImage(weekNumber int) []byte {
    const iconSize = 64
    const fontSize = 50
    dc := gg.NewContext(iconSize, iconSize)
    setFont(dc, "segoeui.ttf", fontSize)

    // Textfarbe je nach Windows-Theme setzen
    if isLightTheme() {
        dc.SetRGB(0, 0, 0) // Dunkler Text für helles Thema
    } else {
        dc.SetRGB(1, 1, 1) // Weißer Text für dunkles Thema
    }

    dc.DrawStringAnchored(fmt.Sprintf("%d", weekNumber), iconSize/2, iconSize/2, 0.5, 0.5)
    return writeContextToByteArray(dc)
}

// Funktion zum Setzen der Schriftart
func setFont(dc *gg.Context, fontName string, fontSize float64) {
    fontPath, err := findfont.Find(fontName)
    if err != nil {
        panic(err)
    }
    if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
        panic(err)
    }
}

// Funktion zum Schreiben des Image-Contexts als Byte-Array
func writeContextToByteArray(dc *gg.Context) []byte {
    var b bytes.Buffer
    foo := bufio.NewWriter(&b)
    ico.Encode(foo, dc.Image())
    foo.Flush()
    return b.Bytes()
}

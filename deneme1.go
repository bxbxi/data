package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
)

type FormData struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

var dataFile = "data.json"

// Ana sayfa (form) HTML
const formHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>Data Entry</title>
</head>
<body>
    <h1>Veri Girişi Formu</h1>
    <form method="POST" action="/submit">
        <label>İsim:</label><br>
        <input type="text" name="name"><br><br>
        <label>Email:</label><br>
        <input type="email" name="email"><br><br>
        <label>Yaş:</label><br>
        <input type="number" name="age"><br><br>
        <input type="submit" value="Gönder">
    </form>
</body>
</html>
`

// Ana sayfa için handler
func homePage(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprint(w, formHTML)
}

// Form verilerini işleyip JSON dosyasına kaydeden handler
func submitForm(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Geçersiz istek metodu", http.StatusMethodNotAllowed)
        return
    }

    name := r.FormValue("name")
    email := r.FormValue("email")
    age := r.FormValue("age")

    // Yaşı integer olarak çevir
    var ageInt int
    fmt.Sscanf(age, "%d", &ageInt)

    formData := FormData{
        Name:  name,
        Email: email,
        Age:   ageInt,
    }

    file, err := os.OpenFile(dataFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        http.Error(w, "Dosya açılamadı", http.StatusInternalServerError)
        return
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    if err := encoder.Encode(formData); err != nil {
        http.Error(w, "Veri kaydedilemedi", http.StatusInternalServerError)
        return
    }

    fmt.Fprint(w, "Veri başarıyla kaydedildi!")
}

func main() {
    http.HandleFunc("/", homePage)
    http.HandleFunc("/submit", submitForm)

    fmt.Println("Sunucu çalışıyor: http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
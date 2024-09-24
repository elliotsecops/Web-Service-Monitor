package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "net/smtp"
    "os"
    "time"
)

type Website struct {
    URL          string
    Status       bool
    ResponseTime time.Duration
}

type Config struct {
    Websites        []string `json:"websites"`
    Interval        int      `json:"interval"` // in seconds
    MaxRetries      int      `json:"maxRetries"`
    ResponseTimeMax time.Duration `json:"responseTimeMax"` // in milliseconds
}

func checkWebsite(website *Website, maxRetries int, responseTimeMax time.Duration) {
    for i := 0; i < maxRetries; i++ {
        start := time.Now()
        resp, err := http.Get(website.URL)
        if err == nil && resp.StatusCode == http.StatusOK {
            website.Status = true
            website.ResponseTime = time.Since(start)
            if website.ResponseTime > responseTimeMax {
                log.Printf("ALERT: %s is slow (Response time: %v)\n", website.URL, website.ResponseTime)
                err := sendEmailAlert("admin@example.com", "Website Slow Alert", fmt.Sprintf("%s is slow (Response time: %v)", website.URL, website.ResponseTime))
                if err != nil {
                    log.Printf("Error sending email alert: %v", err)
                }
            } else {
                log.Printf("Response time for %s: %v", website.URL, website.ResponseTime)
            }
            return
        }
        time.Sleep(time.Second) // Wait one second between retries
    }
    website.Status = false
}

func monitorWebsites(websites []Website, interval int, maxRetries int, responseTimeMax time.Duration) {
    for {
        log.Println("Starting website verification...")
        for i := range websites {
            go checkWebsite(&websites[i], maxRetries, responseTimeMax)
        }
        time.Sleep(5 * time.Second)

        for _, site := range websites {
            if !site.Status {
                log.Printf("ALERT: %s is not responding\n", site.URL)
                err := sendEmailAlert("admin@example.com", "Website Down Alert", fmt.Sprintf("%s is not responding", site.URL))
                if err != nil {
                    log.Printf("Error sending email alert: %v", err)
                }
            } else {
                log.Printf("%s is functioning correctly\n", site.URL)
            }
        }
        log.Printf("Waiting %d seconds for the next verification...\n", interval)
        time.Sleep(time.Duration(interval) * time.Second)
    }
}

func loadConfig(filename string) (Config, error) {
    var config Config
    file, err := os.Open(filename)
    if err != nil {
        return config, err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    err = decoder.Decode(&config)
    return config, err
}

func sendEmailAlert(to, subject, body string) error {
    from := os.Getenv("EMAIL_FROM")
    password := os.Getenv("EMAIL_PASSWORD")
    smtpHost := os.Getenv("SMTP_HOST")
    smtpPort := os.Getenv("SMTP_PORT")

    message := []byte("Subject: " + subject + "\r\n" +
        "\r\n" +
        body + "\r\n")

    auth := smtp.PlainAuth("", from, password, smtpHost)

    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
    if err != nil {
        return err
    }
    return nil
}

func main() {
    config, err := loadConfig("config.json")
    if err != nil {
        log.Fatal("Error loading configuration:", err)
    }

    logFile, err := os.OpenFile("monitor.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer logFile.Close()
    log.SetOutput(logFile)

    var websites []Website
    for _, url := range config.Websites {
        websites = append(websites, Website{URL: url})
    }

    monitorWebsites(websites, config.Interval, config.MaxRetries, config.ResponseTimeMax*time.Millisecond)
}

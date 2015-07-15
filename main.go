package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/jordan-wright/email"
	"github.com/zachlatta/postman/mail"
)

type Recipient map[string]string

type RandTempDesc struct {
	Name  string
	Items []string
}

var (
	htmlTemplatePath, textTemplatePath        string
	csvPath                                   string
	randTemp                                  string
	smtpURL, smtpUser, smtpPassword, smtpPort string
	sender, subject                           string
	attach                                    string
	files                                     []string
	debug                                     bool
	workerCount                               int
	freqPerMinute                             int
	freqMinutes                               int
)

var flags, requiredFlags []*flag.Flag

func main() {
	flag.StringVar(&htmlTemplatePath, "html", "", "html template path")
	flag.StringVar(&textTemplatePath, "text", "", "text template path")
	flag.StringVar(&csvPath, "csv", "", "path to csv of contact list")
	flag.StringVar(&randTemp, "rand", "", "path to json file that descripte which column should pick item from list randomly")
	flag.StringVar(&smtpURL, "server", "", "url of smtp server")
	flag.StringVar(&smtpPort, "port", "", "port of smtp server")
	flag.StringVar(&smtpUser, "user", "", "smtp username")
	flag.StringVar(&smtpPassword, "password", "", "smtp password")
	flag.StringVar(&sender, "sender", "", "email to send from")
	flag.StringVar(&subject, "subject", "", "subject of email")
	flag.BoolVar(&debug, "debug", false, "print emails to stdout instead of sending")
	flag.StringVar(&attach, "attach", "", "attach a list of comma separated files")
	flag.IntVar(&workerCount, "c", 8, "number of concurrent requests to have")
	flag.IntVar(&freqPerMinute, "freq", 12, "number of requests in x minutes")
	flag.IntVar(&freqMinutes, "fmin", 1, "number of requests in x minutes, x value")

	requiredFlagNames := []string{"text", "csv", "server", "port", "user",
		"password", "sender", "subject"}
	flag.VisitAll(func(f *flag.Flag) {
		flags = append(flags, f)

		for _, name := range requiredFlagNames {
			if name == f.Name {
				requiredFlags = append(requiredFlags, f)
			}
		}
	})

	flag.Usage = usage

	flag.Parse()

	if attach != "" {
		files = strings.Split(attach, ",")
	} else {
		files = []string{}
	}

	checkAndHandleMissingFlags(requiredFlags)

	var throttle <-chan time.Time
	if freqPerMinute > 0 {
		throttle = time.Tick(time.Duration((freqMinutes*60*1e6)/freqPerMinute) * time.Microsecond)
	}

	csv, err := os.Open(csvPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening CSV:", err.Error())
		os.Exit(2)
	}
	defer csv.Close()

	jsonFile, err := os.Open(randTemp)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening randTemp:", err.Error())
		os.Exit(2)
	}
	defer jsonFile.Close()

	var randDesc []RandTempDesc
	dec := json.NewDecoder(jsonFile)
	err = dec.Decode(&randDesc)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Decode json failed", err.Error())
		os.Exit(2)
	}

	recipients, emailField, err := readCSV(csvPath, randDesc)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading CSV:", err.Error())
		os.Exit(2)
	}

	mailer := mail.NewMailer(
		smtpUser,
		smtpPassword,
		smtpURL,
		smtpPort,
	)

	jobs := make(chan Recipient, len(*recipients))
	success := make(chan *email.Email)
	fail := make(chan error)

	// Start workers
	for i := 0; i < workerCount; i++ {
		go func() {
			if freqPerMinute > 0 {
				<-throttle
			}
			for recipient := range jobs {
				sendMail(recipient, *emailField, &mailer, debug, success, fail)
			}
		}()
	}

	// Send jobs to workers
	for _, recipient := range *recipients {
		jobs <- recipient
	}
	close(jobs)

	for i := 0; i < len(*recipients); i++ {
		select {
		case msg := <-success:
			if !debug {
				log.Printf("\rEmailed recipient %d of %d...", i+1, len(*recipients))
			} else {
				bytes, err := msg.Bytes()
				if err != nil {
					log.Printf("Error parsing email: %v", err)
				}
				log.Printf("%s\n\n\n", string(bytes))
			}
			log.Println("sending recipient success", msg.To)
		case err := <-fail:
			log.Println("\nError sending email:", err.Error())
			continue
		}
	}
	fmt.Println()
}

func checkAndHandleMissingFlags(requiredFlags []*flag.Flag) {
	var flagsMissing []*flag.Flag
	for _, f := range requiredFlags {
		if f.Value.String() == "" {
			flagsMissing = append(flagsMissing, f)
		}
	}

	missingCount := len(flagsMissing)
	if missingCount > 0 {
		if missingCount == len(requiredFlags) {
			usage()
		}

		missingFlags(flagsMissing)
	}
}

const usageTemplate = `Postman is a utility for sending batch emails.

Usage:

  postman [flags]

Flags:
{{range .}}
  -{{.Name | printf "%-11s"}} {{.Usage}}{{end}}

`

const missingFlagsTemplate = `Missing required flags:
{{range .}}
  -{{.Name | printf "%-11s"}} {{.Usage}}{{end}}

`

func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func printUsage(w io.Writer) {
	tmpl(w, usageTemplate, flags)
}

func usage() {
	printUsage(os.Stderr)
	os.Exit(2)
}

func printMissingFlags(w io.Writer, missingFlags []*flag.Flag) {
	tmpl(w, missingFlagsTemplate, missingFlags)
}

func missingFlags(missingFlags []*flag.Flag) {
	printMissingFlags(os.Stderr, missingFlags)
	os.Exit(2)
}

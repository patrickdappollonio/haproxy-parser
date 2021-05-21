package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"text/tabwriter"
)

type IPData struct {
	URL   string
	IP    string
	Count int
}

var (
	reURL        = regexp.MustCompile(`\"uri\"\:\"([^"]+)\"`)
	reRemoteAddr = regexp.MustCompile(`\"remote_addr\"\:\"([^"]+)\"`)
)

func main() {
	st, err := os.Stdin.Stat()
	if err != nil {
		errexit("Unable to read from stdin: %s", err.Error())
		return
	}

	if (st.Mode() & os.ModeNamedPipe) != os.ModeNamedPipe {
		errexit(
			"This program doesn't accept input from stdin. You must pipe the output of the HAProxy Access Logs:\n" +
				"Example: kubectl logs -n ingress-haproxy haproxy-ingress-controller-abcde -c access-log | haproxy-parser",
		)
		return
	}

	var buf bytes.Buffer

	if _, err := io.Copy(&buf, os.Stdin); err != nil {
		errexit("Unable to read contents from piped output: %s", err.Error())
		return
	}

	scanner := bufio.NewScanner(bytes.NewReader(buf.Bytes()))
	records := []IPData{}

	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()

		if len(line) == 0 {
			continue
		}

		if len(line) >= 1 && line[0] != '{' {
			continue
		}

		url, ip := "", ""

		if match := reURL.FindAllStringSubmatch(line, -1); len(match) == 1 && len(match[0]) == 2 {
			url = match[0][1]
		}

		if match := reRemoteAddr.FindAllStringSubmatch(line, -1); len(match) == 1 && len(match[0]) == 2 {
			ip = match[0][1]
		}

		if ip == "" {
			ip = "(truncated)"
		}

		urlpos := -1
		for i := 0; i < len(records); i++ {
			if records[i].URL == url && records[i].IP == ip {
				urlpos = i
			}
		}

		if urlpos == -1 {
			records = append(records, IPData{
				URL:   url,
				IP:    ip,
				Count: 1,
			})
		} else {
			records[urlpos].Count = records[urlpos].Count + 1
		}
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Count > records[j].Count
	})

	wr := tabwriter.NewWriter(os.Stdout, 6, 1, 4, ' ', 0)

	printrow(wr, "REMOTE ADDRESS", "COUNT", "ENDPOINT")
	for _, v := range records {
		printrow(wr, v.IP, v.Count, v.URL)
	}

	if err := wr.Flush(); err != nil {
		errexit("Unable to flush table output: %s", err.Error())
		return
	}
}

func printrow(w io.Writer, a ...interface{}) {
	var x bytes.Buffer

	for i := 0; i < len(a); i++ {
		x.WriteString(fmt.Sprintf("%v", a[i]))

		if i+1 != len(a) {
			x.WriteString("\t")
		}
	}

	x.WriteString("\n")
	x.WriteTo(w)
}

func errexit(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

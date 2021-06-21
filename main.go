package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "githubtokenechange_total_ops",
		Help: "The total number token exchanges",
	})
	opsFailed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "githubtokenechange_total_errors",
		Help: "The total number of errors",
	})
	githubDAO = NewDAO()
)

func main() {
	checks()
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", handler)
	log.Printf("listening on: %s\n", os.Getenv("PORT"))
	log.Println(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil))
}

func checks() {
	if os.Getenv("PORT") == "" {
		panic("must set PORT variable")
	}
	if os.Getenv("GITHUB_REDIRECT_URL") == "" {
		panic("must set GITHUB_REDIRECT_URL variable")
	}
	if os.Getenv("GITHUB_CLIENT_ID") == "" {
		panic("must set GITHUB_CLIENT_ID variable")
	}
	if os.Getenv("GITHUB_CLIENT_SECRET") == "" {
		panic("must set GITHUB_CLIENT_SECRET variable")
	}
	if os.Getenv("ALLOWLIST_REDIRECT_URLS") == "" {
		panic("must set ALLOWLIST_REDIRECT_URLS variable")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// we only want to process the root request, this ignores all other queries - for example, from the browser.
		// as we don't want to count those as errors.
		return
	}
	opsProcessed.Inc()
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "no github code found in request", http.StatusBadRequest)
		opsFailed.Inc()
		return
	}
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	redirectURL := os.Getenv("GITHUB_REDIRECT_URL")
	allowlistString := os.Getenv("ALLOWLIST_REDIRECT_URLS")

	user, err := githubDAO.GetUser(clientID, clientSecret, code, redirectURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		opsFailed.Inc()
		return
	}

	allowlist := strings.Split(allowlistString, ",")
	url := "http://localhost:3827"
	requestedRedirectURL := r.URL.Query().Get("redirect")

	if (len(requestedRedirectURL) > 0) {
		// check if in allowlist
		for _, entry := range allowlist {
			entry = strings.TrimSpace(entry)
			if(len(entry) > 0) {
				match, _ := regexp.MatchString("^" + entry, requestedRedirectURL)
				if(match) {
					url = requestedRedirectURL
					break
				}
			}
		}
	}

	html := `
		<script>
			let user = ` + string(user) + `
			window.location.href = '` + string(url) + `' + '?user=' + encodeURIComponent(JSON.stringify(user))
		</script>
		`
	_, err = w.Write([]byte(html))
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		opsFailed.Inc()
		return
	}
}

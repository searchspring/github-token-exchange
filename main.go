package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	user, err := githubDAO.GetUser(clientID, clientSecret, code, redirectURL)
	if err != nil {
		http.Error(w, err.Error(), 403)
		opsFailed.Inc()
		return
	}
	html := `
		<script>
			let user = ` + string(user) + `
			window.location.href = 'http://localhost:3827/?user=' + encodeURIComponent(JSON.stringify(user))
		</script>
		`
	w.Write([]byte(html))
}

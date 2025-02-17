package front

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/WLM1ke/poptimizer/opt/internal/domain"
	"github.com/WLM1ke/poptimizer/opt/internal/domain/data/raw"
	"github.com/WLM1ke/poptimizer/opt/internal/domain/data/securities"
	"github.com/WLM1ke/poptimizer/opt/internal/domain/portfolio/port"
	"github.com/WLM1ke/poptimizer/opt/pkg/lgr"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:embed static
var static embed.FS

// Frontend представляет web-интерфейс приложения.
type Frontend struct {
	mux    *chi.Mux
	files  fs.FS
	logger *lgr.Logger

	json domain.JSONViewer

	tickers   *securities.Service
	dividends *raw.Service
	accounts  *port.AccountsService
	portfolio *port.PortfolioService
}

// NewFrontend создает обработчики, отображающие frontend.
//
// Основная страничка расположена в корне. Отдельные разделы динамические отображаются с помощью Alpine.js.
func NewFrontend(logger *lgr.Logger, client *mongo.Client, pub domain.Publisher) http.Handler {
	static, err := fs.Sub(static, "static")
	if err != nil {
		logger.Panicf("can't load frontend data -> %s", err)
	}

	front := Frontend{
		mux:   chi.NewRouter(),
		files: static,

		logger: logger,

		json: domain.NewMongoJSON(client),

		tickers: securities.NewService(
			domain.NewRepo[securities.Table](client),
			pub,
		),
		dividends: raw.NewService(
			domain.NewRepo[securities.Table](client),
			domain.NewRepo[raw.Table](client),
			pub,
		),
		accounts:  port.NewAccountsService(domain.NewRepo[port.Portfolio](client)),
		portfolio: port.NewPortfolioService(domain.NewRepo[port.Portfolio](client)),
	}

	front.registerJSON()

	front.registerMainPage()
	front.registerTickersHandlers()
	front.registerDividendsHandlers()
	front.registerAccountsHandlers()
	front.registerPortfolioHandlers()

	return &front
}

func (f Frontend) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	f.mux.ServeHTTP(writer, request)
}

func (f Frontend) registerMainPage() {
	f.mux.Handle("/{file}", http.StripPrefix("/", http.FileServer(http.FS(f.files))))

	index := template.Must(template.ParseFS(f.files, "index.html"))

	f.mux.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html;charset=UTF-8")

		err := index.Execute(writer, nil)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)

			f.logger.Warnf("can't render template -> %s", err)
		}
	})
}

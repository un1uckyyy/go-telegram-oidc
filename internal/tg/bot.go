package tg

import (
	"fmt"
	"net/http"
	"time"

	"github.com/un1uckyyy/go-telegram-oidc/internal/config"
	"github.com/un1uckyyy/go-telegram-oidc/internal/jwt"
	"github.com/un1uckyyy/go-telegram-oidc/internal/oidc"
	"github.com/un1uckyyy/go-telegram-oidc/pkg/logger"
	"github.com/un1uckyyy/go-telegram-oidc/pkg/storage"
	"github.com/un1uckyyy/go-telegram-oidc/pkg/ticket"
	tele "gopkg.in/telebot.v4"
)

type UserStorage interface {
	GetUser(telegramId int64) (userId string, ok bool)
	SetUser(userId string, telegramId int64) error
}

type Service struct {
	oidcService *oidc.OidcService
	db          UserStorage
	bot         *tele.Bot
}

func (t *Service) Start() {
	t.bot.Start()
}

func NewService() (*Service, error) {
	token := config.Instance.TgToken

	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	ticketService := ticket.NewKeyValue()
	dbService := storage.NewDb()
	oidcService, err := oidc.NewOidcService(ticketService)
	if err != nil {
		return nil, err
	}

	service := &Service{
		oidcService: oidcService,
		db:          dbService,
		bot:         bot,
	}

	service.bot.Use(service.authMiddleware)
	service.bot.Handle(tele.OnText, textHandler)

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/auth", service.authHandler)
		if err := http.ListenAndServe(":8080", mux); err != nil {
			logger.ErrorLogger.Fatal(err)
		}
	}()

	return service, nil
}

func textHandler(c tele.Context) error {
	text := fmt.Sprintf("echo '%v' from %v %v", c.Text(), c.Sender().ID, c.Sender().Username)
	return c.Send(text)
}

func (t *Service) authMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		telegramId := c.Sender().ID
		_, ok := t.db.GetUser(telegramId)
		if !ok {
			url := t.oidcService.GetAuthUrl(telegramId)
			message := fmt.Sprintf("Please, log in by <a href=\"%s\">link</a>", url)
			return c.Send(message, tele.ModeHTML)
		}
		return next(c)
	}
}

func (t *Service) authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	if state == "" || code == "" {
		logger.ErrorLogger.Printf("received empty state or code. state: '%s' code: '%s'", state, code)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := t.oidcService.CompleteAuth(r.Context(), state, code)
	if err != nil {
		logger.ErrorLogger.Printf("failed to complete auth: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.InfoLogger.Println("JWT token:", token.AccessToken)
	sub, err := jwt.GetSubjectFromJwt(token.AccessToken)
	if err != nil {
		logger.ErrorLogger.Printf("failed to get subject from token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	telegramId, ok := t.oidcService.UserTracker.PopTmpUserInfo(state)
	if !ok {
		logger.ErrorLogger.Printf("failed to pop user info")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = t.db.SetUser(sub, telegramId)

	user := &tele.User{ID: telegramId}

	_, err = t.bot.Send(user, "You are successfully logged in!")
	if err != nil {
		logger.ErrorLogger.Println("failed to send message about successful authentication", err)
	}

	botUrl := fmt.Sprintf("https://t.me/%s", t.bot.Me.Username)
	http.Redirect(w, r, botUrl, http.StatusSeeOther)
}

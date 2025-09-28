package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/SomeHowMicroservice/gateway/common"
	"github.com/SomeHowMicroservice/gateway/config"
	"github.com/SomeHowMicroservice/gateway/event"
	"github.com/SomeHowMicroservice/gateway/initialization"
	"github.com/SomeHowMicroservice/gateway/mq"
	"github.com/SomeHowMicroservice/gateway/socket"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

var (
	authAddr    = "localhost:8081"
	userAddr    = "localhost:8082"
	productAddr = "localhost:8083"
	postAddr    = "localhost:8084"
	chatAddr    = "localhost:8085"
)

type Server struct {
	cfg        *config.Config
	httpServer *http.Server
	clients    *initialization.GRPCClients
	hub        *socket.Hub
	sseManager *event.Manager
	router     *message.Router
	wmProduct  *initialization.WatermillSubscriber
	wmPost     *initialization.WatermillSubscriber
}

func NewServer(cfg *config.Config) (*Server, error) {
	authAddr = fmt.Sprintf("%s:%d", cfg.App.ServerHost, cfg.Services.AuthPort)
	userAddr = fmt.Sprintf("%s:%d", cfg.App.ServerHost, cfg.Services.UserPort)
	productAddr = fmt.Sprintf("%s:%d", cfg.App.ServerHost, cfg.Services.ProductPort)
	postAddr = fmt.Sprintf("%s:%d", cfg.App.ServerHost, cfg.Services.PostPort)
	chatAddr = fmt.Sprintf("%s:%d", cfg.App.ServerHost, cfg.Services.ChatPort)

	ca := common.ClientAddresses{
		AuthAddr:    authAddr,
		UserAddr:    userAddr,
		ProductAddr: productAddr,
		PostAddr:    postAddr,
		ChatAddr:    chatAddr,
	}

	clients, err := initialization.InitClients(ca)
	if err != nil {
		return nil, fmt.Errorf("init clients thất bại: %w", err)
	}

	hub := socket.NewHub(clients.ChatClient)
	go hub.Run()

	sseManager := event.NewManager()
	go sseManager.Run()

	logger := watermill.NewStdLogger(false, false)
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, err
	}

	router.AddMiddleware(
		middleware.CorrelationID,
		middleware.Retry{
			MaxRetries:      5,
			InitialInterval: time.Microsecond,
			Multiplier:      1.5,
			MaxInterval:     5 * time.Microsecond,
			Logger:          logger,
		}.Middleware,
		middleware.Recoverer,
	)

	wmProduct, err := initialization.InitWatermill(cfg, logger, common.ProductExchange)
	if err != nil {
		return nil, err
	}

	wmPost, err := initialization.InitWatermill(cfg, logger, common.PostExchange)
	if err != nil {
		return nil, err
	}

	mq.RegisterProductImageUploadedConsumer(router, wmProduct.Subscriber, sseManager)
	mq.RegisterPostImageUploadedConsumer(router, wmPost.Subscriber, sseManager)

	go func() {
		if err := router.Run(context.Background()); err != nil {
			log.Printf("Lỗi chạy message router: %v", err)
		}
	}()

	httpServer, err := NewHttpServer(cfg, clients, hub, sseManager)
	if err != nil {
		return nil, err
	}

	return &Server{
		cfg,
		httpServer,
		clients,
		hub,
		sseManager,
		router,
		wmProduct,
		wmPost,
	}, nil
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) {
	log.Println("Đang shutdown service...")

	if s.router != nil {
		s.router.Close()
	}
	if s.wmPost != nil {
		s.wmPost.Close()
	}
	if s.wmProduct != nil {
		s.wmProduct.Close()
	}
	if s.clients != nil {
		s.clients.Close()
	}
	if s.hub != nil {
		for _, clients := range s.hub.Conversations {
			for client := range clients {
				close(client.Send)
			}
		}
	}
	if s.sseManager != nil {
		for _, client := range s.sseManager.Clients {
			close(client.Send)
		}
	}
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Printf("Shutdown http server thất bại: %v", err)
			return
		}
	}

	log.Println("Shutdown service thành công")
}

func (s *Server) GracefulShutdown(ch <-chan error) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-ch:
		log.Printf("Chạy service thất bại: %v", err)
	case <-ctx.Done():
		log.Println("Có tín hiệu dừng server")
	}

	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.Shutdown(shutdownCtx)
}

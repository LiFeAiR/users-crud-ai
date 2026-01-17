package server

import (
	"context"
	"log"
	"net/http"

	"github.com/LiFeAiR/crud-ai/internal/handlers"
	"github.com/LiFeAiR/crud-ai/internal/repository"
	gw "github.com/LiFeAiR/crud-ai/pkg/server/grpc"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// Server представляет HTTP сервер
type Server struct {
	portHTTP    string
	portProm    string
	connStr     string
	secretKey   string
	db          *repository.DB
	baseHandler *handlers.BaseHandler
}

// NewServer создает новый экземпляр сервера
func NewServer(portHttp, portProm string, connStr, secretKey string) *Server {
	return &Server{
		portHTTP:  portHttp,
		portProm:  portProm,
		connStr:   connStr,
		secretKey: secretKey,
	}
}

// DB возвращает подключение к базе данных
func (s *Server) DB() *repository.DB {
	return s.db
}

func (s *Server) BaseHandler() gw.CrudServiceServer {
	return s.baseHandler
}

func (s *Server) Close() error {
	return s.db.Close()
}

// Start запускает HTTP сервер
func (s *Server) Start(ctx context.Context) error {
	// Подключаемся к базе данных
	db, err := repository.NewDB(s.connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	s.db = db

	userRepo := repository.NewUserRepository(db)
	orgRepo := repository.NewOrganizationRepository(db)
	permRepo := repository.NewPermissionRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	baseHandler := handlers.NewBaseHandler(userRepo, orgRepo, permRepo, roleRepo, s.secretKey)
	s.baseHandler = baseHandler
	defer s.Close()

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)

	gw.RegisterCrudServiceServer(grpcServer, s.BaseHandler())

	var group errgroup.Group

	//lis, err := net.Listen("tcp", "localhost:8081")
	//if err != nil {
	//	return err
	//}

	//group.Go(func() error {
	//	return grpcServer.Serve(lis)
	//})

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}))

	group.Go(func() error {
		err := gw.RegisterCrudServiceHandlerServer(ctx, mux, s.BaseHandler())
		if err != nil {
			return err
		}

		log.Printf("CrudService Listening on :%s...", s.portHTTP)
		return http.ListenAndServe(":"+s.portHTTP, mux)
	})

	group.Go(func() error {
		log.Printf("Promhttp Handler Listening on :%s...", s.portProm)
		return http.ListenAndServe(":"+s.portProm, promhttp.Handler())
	})

	return group.Wait()
}

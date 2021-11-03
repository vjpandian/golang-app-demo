package http

import (
	"context"

	"github.com/aristat/golang-example-app/app/dataloader"

	products_router "github.com/aristat/golang-example-app/app/http_routers/products-router"

	"github.com/aristat/golang-example-app/app/auth"

	"github.com/aristat/golang-example-app/app/graphql"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/opentracing/opentracing-go"

	"github.com/aristat/golang-example-app/app/logger"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

var mux *chi.Mux

// Cfg
func Cfg(cfg *viper.Viper) (Config, func(), error) {
	c := Config{}
	e := cfg.UnmarshalKey("http", &c)
	if e != nil {
		return c, func() {}, nil
	}
	c.Debug = cfg.GetBool("debug")
	return c, func() {}, nil
}

// CfgTest
func CfgTest() (Config, func(), error) {
	return Config{}, func() {}, nil
}

// Mux
func Mux(managers Managers, log logger.Logger, tracer opentracing.Tracer) (*chi.Mux, func(), error) {
	if mux != nil {
		return mux, func() {}, nil
	}

	mux = chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(Logger(log))
	mux.Use(Tracer(tracer))
	mux.Use(dataloader.LoaderMiddleware)

	managers.products.Router.Run(mux)
	managers.graphql.Routers(mux.With(managers.authMiddleware.JWTHandler))

	return mux, func() {}, nil
}

// ServiceManagers
type Managers struct {
	products       *products_router.Manager
	authMiddleware *auth.Middleware
	graphql        *graphql.GraphQL
}

var ProviderManagers = wire.NewSet(
	wire.Struct(new(Managers), "*"),
)

// Provider
func Provider(ctx context.Context, mux *chi.Mux, log logger.Logger, cfg Config) (*Http, func(), error) {
	g := New(ctx, mux, log, cfg)
	return g, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, Cfg, Mux, ProviderManagers)
	ProviderTestSet       = wire.NewSet(Provider, CfgTest)
)

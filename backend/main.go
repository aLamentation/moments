package main

import (
	"embed"
	"fmt"
	"net/http"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kingwrcy/moments/db"
	_ "github.com/kingwrcy/moments/docs"
	"github.com/kingwrcy/moments/handler"
	"github.com/kingwrcy/moments/log"
	myMiddleware "github.com/kingwrcy/moments/middleware"
	"github.com/kingwrcy/moments/vo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
	_ "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

var version string
var commitId string

var staticFiles *embed.FS

func newEchoEngine(_ do.Injector) (*echo.Echo, error) {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins:     []string{"https://moments.l-zs.com", "https://moments.alamentation.xyz"}, // 前端实际部署地址
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Content-Type", "Authorization", "x-api-token"},
        AllowCredentials: true,
    }))


	return e, nil
}

// @title	Moments API
// @version	0.2.1
func main() {
	injector := do.New()
	var cfg vo.AppConfig

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		fmt.Printf("读取配置文件异常:%s", err)
		return
	}

	if version == "" {
		version = "unknown"
	}

	if commitId == "" {
		commitId = "unknown"
	}

	do.ProvideValue(injector, &cfg)
	do.Provide(injector, log.NewLogger)

	myLogger := do.MustInvoke[zerolog.Logger](injector)

	myLogger.Info().Msgf("version = %s", version)
	myLogger.Info().Msgf("commitId = %s", commitId)

	handleEmptyConfig(myLogger, &cfg)
	cfg.Version = version
	cfg.CommitId = commitId

	backupDatabase(myLogger, &cfg)

	migrateFriendLink(myLogger, &cfg)

	do.Provide(injector, db.NewDB)
	do.Provide(injector, newEchoEngine)
	do.Provide(injector, handler.NewBaseHandler)

	tx := do.MustInvoke[*gorm.DB](injector)

	e := do.MustInvoke[*echo.Echo](injector)
	e.Use(myMiddleware.Auth(injector))

	setupRouter(injector)

	if staticFiles != nil {
		e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:       "public",
			HTML5:      true,
			IgnoreBase: false,
			Browse:     false,
			Filesystem: http.FS(staticFiles),
			Skipper: func(c echo.Context) bool {
				return strings.HasPrefix(c.Request().URL.Path, "/swagger/")
			},
		}))
	}

if cfg.CorsOrigin != "" {
    allowOrigins := strings.Split(cfg.CorsOrigin, ",")
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowCredentials: true,
        AllowOrigins:     allowOrigins,
        AllowMethods:     []string{
            http.MethodGet,
            http.MethodPost,
            http.MethodPut,
            http.MethodDelete,
            http.MethodOptions, // 👈 必须加
        },
        AllowHeaders: []string{
            echo.HeaderOrigin,
            echo.HeaderContentType,
            echo.HeaderAccept,
            echo.HeaderAuthorization,
            "X-API-TOKEN", // 👈 自定义头
        },
    }))
}


	migrateTo3(tx, myLogger)
	migrateIframeVideoUrl(tx, myLogger)

	e.HideBanner = true
	err = e.Start(fmt.Sprintf(":%d", cfg.Port))
	if err == nil {
		myLogger.Info().Msgf("服务端启动成功,监听:%d端口...", cfg.Port)
	} else {
		myLogger.Fatal().Msgf("服务启动失败,错误原因:%s", err)
	}
}

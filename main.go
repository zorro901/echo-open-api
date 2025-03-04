package main

import (
	"echo-open-api/api"
	"github.com/getkin/kin-openapi/openapi3filter"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	middleware "github.com/oapi-codegen/echo-middleware"
)

// ServerImpl はサーバーインターフェースを実装します
type ServerImpl struct{}

// PostEcho はエコーリクエストを処理します
func (s *ServerImpl) PostEcho(ctx echo.Context) error {
	var request api.EchoRequest
	if err := ctx.Bind(&request); err != nil {
		return err
	}

	// 文字列のコピーを作成してそのポインタを使用
	message := request.Message
	response := api.EchoResponse{
		Message: &message,
	}

	return ctx.JSON(http.StatusOK, response)
}

func main() {
	e := echo.New()

	// OpenAPI仕様を取得
	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatal("OpenAPIの仕様を読み込めませんでした:", err)
	}

	// 検証に使用しないため、Serversフィールドをクリア
	swagger.Servers = nil

	// オプションを設定
	options := &middleware.Options{
		ErrorHandler: func(c echo.Context, err *echo.HTTPError) error {
			return c.JSON(err.Code, map[string]string{
				"message": err.Message.(string),
			})
		},
		Options: openapi3filter.Options{},
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/health"
		},
		SilenceServersWarning: true,
	}

	// OpenAPI検証ミドルウェアを設定
	//e.Use(middleware.OapiRequestValidator(swagger))
	e.Use(middleware.OapiRequestValidatorWithOptions(swagger, options))

	// サーバー実装を作成
	server := &ServerImpl{}

	// ハンドラを登録
	api.RegisterHandlers(e, server)

	// サーバーを起動
	log.Println("サーバーを起動します: http://localhost:8080")
	e.Logger.Fatal(e.Start(":8080"))
}

# oapi-codegenを使用したEchoサーバー実装例

以下はoapi-codegenを使用してEchoでサーバーを実装する方法です：

## OpenAPI仕様

まずは簡単な仕様から始めましょう：

```yaml
openapi: "3.0.0"
info:
  version: 1.0.0
  title: 最小限のping APIサーバー
paths:
  /ping:
    get:
      responses:
        '200':
          description: レスポンス
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Pong'
components:
  schemas:
    # 基本タイプ
    Pong:
      type: object
      required:
        - ping
      properties:
        ping:
          type: string
          example: pong
```

## 生成されるコード

この仕様は以下のようなコードを生成します：

```go
// Pong はPongモデルを定義します。
type Pong struct {
	Ping string `json:"ping"`
}

// ServerInterface はすべてのサーバーハンドラを表します。
type ServerInterface interface {
	// (GET /ping)
	GetPing(ctx echo.Context) error
}

// これはecho.Routeのためのシンプルなインターフェースであり、
// echo.Echoとecho.Groupの両方に存在する関数を指定します。
// パス登録のためにどちらかを使用できるようにするためです。
type EchoRouter interface {
	// ...
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	// ...
}

// RegisterHandlers は各サーバールートをEchoRouterに追加します。
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// RegisterHandlersWithBaseURL はハンドラを登録し、パスにBaseURLを先頭に付加します。
// これによりパスはプレフィックスの下で提供できます。
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {
	// ...
	router.GET(baseURL+"/ping", wrapper.GetPing)
}
```

## 実装

このHTTPサーバーを実装するには、`api/impl.go`というファイルを作成します：

```go
package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Server はServerInterfaceを実装します
type Server struct{}

// NewServer は新しいサーバーインスタンスを作成します
func NewServer() Server {
	return Server{}
}

// GetPing はGET /pingリクエストを処理します
func (Server) GetPing(ctx echo.Context) error {
	resp := Pong{
		Ping: "pong",
	}

	return ctx.JSON(http.StatusOK, resp)
}
```

## メインアプリケーション

すべてを接続するメインアプリケーションを作成します：

```go
package main

import (
	"log"

	"あなたのモジュール/api"
	"github.com/labstack/echo/v4"
)

func main() {
	// 新しいEchoインスタンスを作成
	e := echo.New()
	
	// ServerInterfaceを実装するサーバーインスタンスを作成
	server := api.NewServer()

	// ハンドラーを登録
	api.RegisterHandlers(e, server)

	// サーバーを起動
	log.Fatal(e.Start("0.0.0.0:8080"))
}
```

## 設定

サーバーコードを生成するため、`cfg.yaml`という設定ファイルを作成します：

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/oapi-codegen/oapi-codegen/HEAD/configuration-schema.json
package: api
generate:
  echo-server: true
  models: true
output: server.gen.go
```

そして、コードに生成コマンドを追加します：

```go
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config cfg.yaml api.yaml
```

## リクエスト検証（オプション）

リクエスト検証には、echo-middlewareを使用できます：

```go
import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	middleware "github.com/oapi-codegen/echo-middleware"
	"あなたのモジュール/api"
)

func main() {
	// OpenAPI仕様を読み込む
	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatalf("Swaggerスペックの読み込みエラー: %s", err)
	}
	// 検証が失敗するのを防ぐため、swagger仕様のserversの配列をクリア
	swagger.Servers = nil

	// 新しいEchoインスタンスを作成
	e := echo.New()

	// リクエスト検証のためのミドルウェアを使用
	e.Use(middleware.OapiRequestValidator(swagger))
	
	// サーバーインスタンスを作成
	server := api.NewServer()

	// ハンドラーを登録
	api.RegisterHandlers(e, server)

	// サーバーを起動
	log.Fatal(e.Start("0.0.0.0:8080"))
}
```

この例はoapi-codegenでEchoサーバーを実装する簡潔な概要を提供しています。あなたの特定のアプリケーションに必要に応じて、追加のエンドポイントやより複雑なリクエスト/レスポンス構造を拡張できます。
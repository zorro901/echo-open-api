package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	// api.gen.goのimport
	"echo-server/api"
)

// ServerInterfaceを実装する構造体
type Server struct{}

// PostEchoの実装
func (s *Server) PostEcho(ctx echo.Context) error {
	// リクエストボディの受け取り
	var req api.PostEchoJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	// レスポンスの返却
	return ctx.JSON(http.StatusOK, map[string]string{
		"message": req.Message,
	})
}

func main() {
	e := echo.New()

	// ミドルウェアの設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// サーバーインスタンスの作成
	server := &Server{}

	// 生成されたRegisterHandlers関数を使用してルートを登録
	api.RegisterHandlers(e, server)

	// サーバーの起動
	if err := e.Start(":8082"); err != nil {
		log.Fatal(err)
	}
}

package handlertest

// Package handlertest はhandler配下のパッケージのテストで使われる機能をまとめたものです。

// 各テストパッケージで多くの共通する機能があったためパッケージとして独立させた.
// 問題点:
//  1. テスト以外でもコンパイルされるというオーバーヘッドがある
//  2. DBをmapを使ってモックしてみたが, テストを作る上でコストが高すぎる気がする
//     DBのように動く機能を用意せずに, 毎テスト期待する挙動をベタ書きするのも大変(でもこれが主流かも)
//     特に今回は一つの関数内で複数回同じ機能を使用したりしているので, 呼び出されるごとに挙動を変えるような動きをつけたかった
//      -> そういうメソッドを作ればいけるか
// 改善点:
//  各テストパッケージではその中でテスト処理を完結させる(他テストパッケージをimportしない)
//  テストパッケージが大きくなるようなコードを書かない(今回だとhandler系パッケージをさらに分けるべき？)
//  DBモックは次回毎テストごとに挙動を書くやり方でそのコストと見やすさを見てみる
// 留意点:
//  _test.goのみで構成されたディレクトリ, パッケージは作れない(テスト対象パッケージがあってのテストパッケージ)
//  _testパッケージはimportの仕様が(今の所)分からない
// 疑問点:
//  複数のテストパッケージで共有したい機能をテスト時のみコンパイルするような形で実現することは可能か？
//  DBのモックの例を色々知りたい

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"

	"yatter-backend-go/app/app"
)

// 引数が多く機能にまとまりもないので, 使い勝手が良くない
// そもそもC structがいるのか, 使うならもっと活用できそう
func Setup(t *testing.T, db *DBMock, newRouter func(app *app.App) http.Handler) *C {
	db = fillDB(db)
	app := &app.App{Dao: &DaoMock{db: db}}
	server := httptest.NewServer(newRouter(app))
	return &C{
		App:    app,
		Server: server,
	}
}

func fillDB(db *DBMock) *DBMock {
	if db == nil {
		db = new(DBMock)
	}
	if db.Account == nil {
		db.Account = make(AccountTableMock)
	}
	if db.Status == nil {
		db.Status = make(StatusTableMock)
	}
	if db.RelationShip == nil {
		db.RelationShip = make(RelationShipTableMock)
	}
	return db
}

type C struct {
	App    *app.App
	Server *httptest.Server
}

func (c *C) Close() {
	c.Server.Close()
}

func (c *C) Get(apiPath string) (*http.Response, error) {
	return c.Server.Client().Get(c.asURL(apiPath))
}

func (c *C) GetWithAuth(apiPath string, authUser string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.asURL(apiPath), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authentication", "username "+authUser)
	return c.Server.Client().Do(req)
}

func (c *C) PostJSON(apiPath string, payload string) (*http.Response, error) {
	return c.Server.Client().Post(c.asURL(apiPath), "application/json", bytes.NewReader([]byte(payload)))
}

func (c *C) PostJsonWithAuth(apiPath string, payload string, authUser string) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.asURL(apiPath), bytes.NewReader([]byte(payload)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authentication", "username "+authUser)
	return c.Server.Client().Do(req)
}

func (c *C) DeleteWithAuth(apiPath string, authUser string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", c.asURL(apiPath), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authentication", "username "+authUser)
	return c.Server.Client().Do(req)
}

func (c *C) asURL(apiPath string) string {
	baseURL, _ := url.Parse(c.Server.URL)
	baseURL.Path = path.Join(baseURL.Path, apiPath)
	return baseURL.String()
}

複数のテストパッケージに共有される為だけのutil_test.go的なのは作れない
理由1. _test.goは特定のテスト対象に対してのみ作れるので,
      どこかのテストパッケージに属する必要がある
理由2. importパスとしてテストパッケージを指定するやり方がわからない
解決:  通常のパッケージとしてutilを含める(テストでしか使わないが毎度コンパイルされる)
結論:  テストパッケージで行うことはそのパッケージ内で完結させる
結論:  しょうがなく通常のパッケージを用意
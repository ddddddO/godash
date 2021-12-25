fmt:
	go fmt ./...

test:
	go test ./...

run_worker:
	go run cmd/worker/*.go run

# datasourceの接続情報をサーバー側へ保存
run_cli_s:
	go run cmd/cli/main.go settings "postgres://postgres:passw0rd@localhost:15432/dvdrental"

# query実行
run_cli_q:
	go run cmd/cli/main.go query "select first_name, last_name from actor limit 1"

# 準備：https://qiita.com/yomon8/items/10b6cd47dda3fd3921c0
PASSWORD=passw0rd
VERSION=12.4
run_pg:
	docker run --rm -d --name sampledb \
	-p 15432:5432 \
	-e POSTGRES_PASSWORD=$(PASSWORD) \
	-v /home/ochi/pg_sample:/tmp/data \
	postgres:$(VERSION)-alpine

conn_pg:
	PGPASSWORD=$(PASSWORD) psql -U postgres -p 15432 -d dvdrental -h localhost

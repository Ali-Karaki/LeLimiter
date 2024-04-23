run:
	go run ./cmd/main.go
	
run-redis:
	docker-compose up -d && docker-compose exec redis redis-cli

stats:
	awk -f ./utils/logger/stats.awk ./utils/logger/logs/nolua.log ./utils/logger/logs/lua.log ./utils/logger/logs/luasha.log > ./utils/logger/logs/stats.log
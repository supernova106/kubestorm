build:
	docker build -t supernova106/kubestorm-ui:latest -f Dockerfile .
push:
	docker push supernova106/kubestorm-ui:latest
deploy:
	helm upgrade --install kubestorm-ui ./charts/kubestorm-ui --namespace=kubestorm
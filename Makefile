release:
	goreleaser --rm-dist
deploy:
	./scripts/create_namespace.sh
	helm upgrade --install kubestorm ./charts/kubestorm --namespace=kubestorm


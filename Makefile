deploy:
	cd service; make deploy
	cd ui; make deploy
	cd infra; make deploy

image:
	codevis -i ./ --whitelist-extension go,hcl,tf,css,js,templ,

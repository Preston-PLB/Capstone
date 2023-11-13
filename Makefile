deploy:
	cd service; make deploy
	cd ui; make deploy
	cd infra; make deploy

image: SHELL := /bin/bash
image:
	[[ -d "/tmp/capstone" ]] || mkdir /tmp/capstone
	cp -R infra/ /tmp/capstone/
	cp -R service/ /tmp/capstone/
	cp -R ui/ /tmp/capstone/
	rm -rf /tmp/capstone/ui/templates/*_templ.go
	codevis -i /tmp/capstone --whitelist-extension go,hcl,tf,templ -o ./out.png
	rm -rf /tmp/capstone/*

.PHONY: staging production

staging:
	docker build . -t localhost:5000/app/appsrv:staging
	docker push localhost:5000/app/appsrv:staging

production:
	docker build . -t localhost:5000/app/appsrv
	docker push localhost:5000/app/appsrv
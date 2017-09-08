dist=dist
exe=buildMonitorServer

default: build

build: setup server client

clean:
	rm -rf $(dist); cd client; yarn run clean; cd ..

setup: clean
	mkdir $(dist); mkdir $(dist)/client; mkdir $(dist)/server

server: buildLinux
	cp run.sh $(dist); \
	chmod 755 $(dist)/run.sh; \
	chmod 755 $(dist)/server/$(exe)

client: buildClient
	 cp -R client/dist/* $(dist)/client

buildClient: ensureClient
	cd client; \
	rm -rf dist; \
	yarn run build; cd ..

buildLinux: ensureServer
	GOOS=linux go build -o ./$(dist)/server/$(exe) ./server

ensureServer: FORCE
	cd server; dep ensure; cd ..

ensureClient: FORCE
	cd client; yarn install --silent; cd ..

docker: clean buildClient
	docker build -t build-monitor-v2:latest .; \
	docker image prune --force

dockerRun:
	docker run \
	--name build-monitor-v2 \
	--hostname build-monitor-v2 \
	-p 3030:3030 \
	--network host \
	--rm --detach \
	build-monitor-v2

# https://github.com/typicode/json-server
jsonServer:
	json-server --watch db.json --routes routes.json --port 3031

FORCE:
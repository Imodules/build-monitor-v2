dist=dist
exe=build-monitor-v2

default: setup server client

clean:
	rm -rf $(dist); cd client; yarn run clean; cd ..

setup: clean
	mkdir $(dist); mkdir $(dist)/client; mkdir $(dist)/server

server: buildLinux
	cp run.sh $(dist); \
	chmod 755 $(dist)/run.sh

client: buildClient
	 cp -R client/dist/* $(dist)/client

buildClient: ensureClient
	cd client; yarn run build; cd ..

buildLinux: ensureServer
	GOOS=linux go build -o ./$(dist)/server/$(exe) ./server

ensureServer: FORCE
	cd server; dep ensure; cd ..

ensureClient: FORCE
	cd client; yarn install --silent; cd ..

FORCE:
# Set the secrets
salt := "LEKOI3IBQ476CD335QYTKNXUPDFNPIQR"
pepper := "J3KCJGE3ISR2V2Y56SRVPI5S7F5WZQB3"
secretkey := "J23RBAUAWAGBT2UYABUXG5TUVSXWNYSM"
secretms1 := "NAV4AE5E655CN6VSR6EPW4ALDITCTNWS"
secretms2 := "RSWB5UZBPVOBQZ2JDY6SNYJ4LB432N3I"

build: version := dev
tests: version := test

.PHONY: package clean

%:
	@:

format:
	go fmt `go list ./...`
	go vet `go list ./...`

tests:
	go test -race -ldflags "\
		-extldflags '-static' \
		-X 'main.Version=${version}' \
		-X 'mediator/mediatorscript.salt=${salt}' \
		-X 'mediator/mediatorscript.pepper=${pepper}' \
		-X 'mediator/mediatorscript.secretKey=${secretkey}' \
		-X 'mediator/totp.secretMS1=${secretms1}' \
		-X 'mediator/totp.secretMS2=${secretms2}' \
		" `go list ./... | grep -v /vendor/ | grep -v /clones/`

build:
	@for elt in mediator-client mediator-server mediator-cli; do \
		echo "→ Building $${elt}" ; \
		go build -race -ldflags "\
			-extldflags '-static' \
			-X 'main.Version=${version}' \
			-X 'mediator/mediatorscript.salt=${salt}' \
			-X 'mediator/mediatorscript.pepper=${pepper}' \
			-X 'mediator/mediatorscript.secretKey=${secretkey}' \
			-X 'mediator/totp.secretMS1=${secretms1}' \
			-X 'mediator/totp.secretMS2=${secretms2}' \
		" -o "./bin/$${elt}"   "./cmd/$${elt}" ; \
	done

run:
	go run main.go

package: clean

	@$(MAKE) \
		version=$(shell date '+%Y%m%d%H%M%S') \
		salt=$(shell tr -cd '[:alnum:]' < /dev/urandom | tr '[:lower:]' '[:upper:]' | head -c32) \
		pepper=$(shell tr -cd '[:alnum:]' < /dev/urandom  | tr '[:lower:]' '[:upper:]' | head -c32) \
		secretkey=$(shell tr -cd '[:alnum:]' < /dev/urandom | tr '[:lower:]' '[:upper:]' | head -c32) \
		secretms1=$(shell tr -cd '[:alnum:]' < /dev/urandom | head -c20 | base32) \
		secretms2=$(shell tr -cd '[:alnum:]' < /dev/urandom | head -c20 | base32) \
		build

	@echo "→ Copying files"
	@mkdir -p package/mediator/
	@cp bin/* package/mediator/
	@cp ./cmd/mediator-server/mediator-server_dist.yml package/mediator/
	@cp ./cmd/mediator-client/mediator-client_dist.yml package/mediator/
	@sed -i "s/EnterYourSecretHere/$(shell tr -cd '[:alnum:]' < /dev/urandom | head -c32)/" "package/mediator/mediator-server_dist.yml"
	@cp ./cmd/mediator-server/mediator-server.service package/mediator/
	@cp -r ./cmd/scripts/ package/mediator/

	@echo "→ Creating the archive"
	@tar -czf package/mediator.tar.gz -C package mediator
	@rm -rf package/mediator

clean:
	@echo "→ Removing bin/ and package/"
	@rm -rf bin/
	@rm -rf package/

dest = $(filter-out $@,$(MAKECMDGOALS))
upload:
	@echo "→ Copying bin/mediator-client to ${dest}"
	@scp bin/mediator-client ${dest}:/tmp
	@echo "→ Pushing bin/mediator-client to Securechange on ${dest}"
	@ssh ${dest} 'sudo tos scripts sc push --overwrite /tmp/mediator-client && rm /tmp/mediator-client'

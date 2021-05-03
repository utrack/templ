install-snapshot:
	# Clear logs.
	rm cmd/lsp/*.txt || true
	# Install the latest version.
	cd cmd && go build -o ~/bin/templ

build-snapshot:
	goreleaser build --snapshot --rm-dist

release: 
	if [ "${GITHUB_TOKEN}" == "" ]; then echo "No github token, run:"; echo "export GITHUB_TOKEN=`pass github.com/goreleaser_access_token`"; exit 1; fi
	./push-tag.sh
	goreleaser --rm-dist

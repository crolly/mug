
app_bins = bin/course/create bin/course/read bin/course/update bin/course/delete bin/course/list bin/user/create bin/user/read bin/user/update bin/user/delete bin/user/list 
bin/% : functions/%.go
		env GOOS=linux go build -ldflags="-s -w" -o $@

build: vendor | $(app_bins)

vendor: Gopkg.toml
	    dep ensure
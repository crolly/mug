
app_bins = bin/course/create bin/course/read bin/course/update bin/course/delete bin/course/list 
app_debugs = debug/course/create debug/course/read debug/course/update debug/course/delete debug/course/list 
bin/% : functions/%/main.go
		env GOOS=linux go build -ldflags="-s -w" -o $@ $<

debug/%: functions/%/main.go
		 env GOARCH=amd64 GOOS=linux go build -gcflags='-N -l' -o $@ $<

build: vendor | $(app_bins)

debug: vendor | $(app_debugs)

vendor: Gopkg.toml
	    dep ensure
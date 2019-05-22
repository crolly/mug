app_bins = bin/test bin/testtwo bin/course/create bin/course/read bin/course/update bin/course/delete bin/course/list bin/user/create bin/user/read bin/user/update bin/user/delete bin/user/list 
app_debugs = debug/test debug/testtwo debug/course/create debug/course/read debug/course/update debug/course/delete debug/course/list debug/user/create debug/user/read debug/user/update debug/user/delete debug/user/list 

bin/% : functions/%/main.go
		env GOOS=linux go build -ldflags="-s -w" -o $@ $<

debug/%: functions/%/main.go
		 env GOARCH=amd64 GOOS=linux go build -gcflags='-N -l' -o $@ $<

build: vendor | $(app_bins)

debug: vendor | $(app_debugs)

vendor: Gopkg.toml
	    dep ensure
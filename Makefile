
build:
	npm run build
	env GOOS=linux GOARCH=amd64 GOARM=7 go build -i -o ./dist/mongodb-atlas-datasource_linux_amd64 ./pkg
	env GOOS=windows GOARCH=386 go build -i -o ./dist/mongodb-atlas-datasource_windows.exe ./pkg

go get -u -v github.com/quited/toaster/...
mkdir -Force build
go build -v -o build/launcher.exe github.com/quited/toaster/launcher
go build -v -o build/demoservice.exe github.com/quited/toaster/launcher/demoservice
cp $env:GOPATH/src/github.com/quited/toaster/launcher/config.json build/


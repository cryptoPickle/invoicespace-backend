if hash gin 2>/dev/null; then
   gin --appPort 8080 -all -i main.go
else
  go get github.com/codegangsta/gin
  gin --appPort 8080 -all -i main.go
fi
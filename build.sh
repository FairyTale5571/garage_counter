echo "Start building extension grc"
echo "grc_x64.dll in progress..."
GOARCH=amd64 go build -o release/grc_x64.dll -buildmode=c-shared .
echo "grc_x64.dll builded"

echo "grc.dll in progress..."
GOARCH=386 CGO_ENABLED=1 go build -o release/grc.dll -buildmode=c-shared .
echo "grc.dll builded"

rm release/*.h
echo "Auto-generated headers removed"

echo "Building done, find dll's in release folder"
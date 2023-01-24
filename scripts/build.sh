list=("ip" "notify" "static" "proxy" "file-urls")

system="darwin"
ext=""
if [ $system == "windows" ]; then
  ext=".exe"
fi

rm -rf dist

for val in ${list[@]}
do
  echo $val
  CGO_ENABLED=0 GOOS=$system go build -o dist/$val$ext src/$val/main.go
done

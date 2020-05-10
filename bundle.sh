cd bundle
if [ -f bundle.app ]; then
  rm bundle.app
fi
go build -o bundle.app
if [ -f bundle.app ]; then
  ./bundle.app "../"
fi
cd ..

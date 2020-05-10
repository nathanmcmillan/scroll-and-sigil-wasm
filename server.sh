if [ -f ss.app ]; then
  rm ss.app
fi
cd server
go build -o ss.app
cd ..
if [ -f server/ss.app ]; then
  mv server/ss.app .
  ./ss.app 3000 temp $@
fi

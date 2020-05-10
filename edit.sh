if [ -f sse.app ]; then
  rm sse.app
fi
cd editor
go build -o sse.app
cd ..
if [ -f editor/sse.app ]; then
  mv editor/sse.app .
  ./sse.app $@
fi

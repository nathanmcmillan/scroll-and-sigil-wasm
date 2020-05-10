if [ -f public/ss.wasm ]; then
  rm public/ss.wasm
fi
cd client
GOARCH=wasm GOOS=js go build -o ss.wasm
cd ..
if [ -f client/ss.wasm ]; then
  mv client/ss.wasm public/ss.wasm
  cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" public/wasm.js
fi

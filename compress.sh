rm public/mini.js
cat public/mega.js | terser -c -m > public/mini.js

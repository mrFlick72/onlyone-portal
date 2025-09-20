rm -rf ../dist
cd ../app
npm install
npm run-script build
rm -rf node_module
npm ci --omit=dev
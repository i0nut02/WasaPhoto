# WasaPhotoWasaPhoto

https://github.com/i0nut02/WasaPhoto/assets/99051485/921e94da-6b24-4ba6-86fa-56a0708c672c

## Backend
### How to build 
If you're not using the WebUI, or if you don't want to embed the WebUI into the final executable, then:
```
go build ./cmd/webapi/
```
If you're using the WebUI and you want to embed it into the final executable:
```
./open-npm.sh
# (here you're inside the NPM container)
npm run build-embed
exit
# (outside the NPM container)
go build -tags webui ./cmd/webapi/
```

## Frontend
### How to run (in development mode)
You can launch the backend only using:
```
go run ./cmd/webapi/
```
If you want to launch the WebUI, open a new tab and launch:
```
./open-npm.sh
# (here you're inside the NPM container)
npm run dev
```

## Deployment
### How to build the images 
Backend
```
docker build -t wasa-photos-backend:latest -f Dockerfile.backend .
```
Frontend 
```
docker build -t wasa-photos-frontend:latest -f Dockerfile.frontend .
```
### How to run the container images
Backend
```
docker run -it --rm -p 3000:3000 wasa-photos-backend:latest
```
Frontend
```
docker run -it --rm -p 8081:80 wasa-photos-frontend:latest
```

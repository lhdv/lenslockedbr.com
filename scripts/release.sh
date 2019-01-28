#!/bin/bash

# Change to the directory with our code that we plan to work from
cd "$GOPATH/src/lenslockedbr.com"

echo "===== Releasing lenslockedbr.com ====="
echo "  Deleting the local binary if it exists (so it isn't uploaded)..."
rm *.exe
echo "  Done!"

sleep 2

echo "  Packing the code"
cd "$GOPATH/src/"
rm lenslockedbr.com.tar.gz
tar -cvzf lenslockedbr.com.tar.gz --exclude='lenslockedbr.com/.git/*' --exclude='lenslockedbr.com/images/*' --exclude='lenslockedbr.com/*.exe' lenslockedbr.com\

sleep 2

echo "  Deleting existing code..."
ssh root@leandr0.net -p 2233 " rm -rf /root/go/src/lenslockedbr.com"
echo "  Code deleted successfully!"

sleep 2

echo "  Uploading and extract code..."
cd "$GOPATH/src/"
scp -P 2233 lenslockedbr.com.tar.gz root@leandr0.net:/root/go/src/
ssh root@leandr0.net -p 2233 " tar -zxvf /root/go/src/lenslockedbr.com.tar.gz -C /root/go/src/"
#echo "  Code uploaded successfully!"

sleep 2

echo "  Go getting deps..."
ssh root@leandr0.net -p 2233 "export GOPATH=/root/go; /usr/local/go/bin/go get golang.org/x/crypto/bcrypt"

ssh root@leandr0.net -p 2233 "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/gorilla/mux"

ssh root@leandr0.net -p 2233 "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/gorilla/schema"

ssh root@leandr0.net -p 2233 "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/gorilla/csrf"

ssh root@leandr0.net -p 2233 "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/lib/pq"

ssh root@leandr0.net -p 2233 "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/jinzhu/gorm"

ssh root@leandr0.net -p 2233 "export GOPATH=/root/go; /usr/local/go/bin/go get gopkg.in/mailgun/mailgun-go.v1"

ssh root@leandr0.net -p 2233 "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"

ssh root@leandr0.net -p 2233 "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"

ssh root@leandr0.net -p 2233 "export GOPATH=/root/go; /usr/local/go/bin/go get golang.org/x/oauth2"

sleep 2

echo "  Building the code on remote server..."
ssh root@leandr0.net -p 2233 "export GOPATH=/root/go;  cd /root/app;  /usr/local/go/bin/go build -o ./server /root/go/src/lenslockedbr.com/*.go"
echo "  Code built successfully!..."

sleep 2

echo "  Moving assets..."
ssh root@leandr0.net -p 2233 " cd /root/app;  cp -R /root/go/src/lenslockedbr.com/assets ."
echo "  Assets moved successfully!..."

echo "  Moving views..."
ssh root@leandr0.net -p 2233 " cd /root/app;  cp -R /root/go/src/lenslockedbr.com/views ."
echo "  Views moved successfully!..."

echo "  Moving Caddyfile..."
ssh root@leandr0.net -p 2233 " cd /root/app;  cp /root/go/src/lenslockedbr.com/Caddyfile ."
echo "  Caddyfile moved successfully!..."

sleep 2

echo "  Restarting the server..."
ssh root@leandr0.net -p 2233 " service leandr0.net restart"
echo "  Server restarted successfully!..."

sleep 2

echo "  Restarting Caddy server..."
ssh root@leandr0.net -p 2233 " service caddy restart"
echo "  Caddy restarted successfully!..."

echo "===== Done releasing lenslockedbr.com ====="

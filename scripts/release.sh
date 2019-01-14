#!/bin/bash

# Change to the directory with our code that we plan to work from
cd "$GOPATH/src/lenslockedbr.com"

echo "===== Releasing lenslockedbr.com ====="
echo "  Deleting the local binary if it exists (so it isn't uploaded)..."
rm lenslockedbr.com.exe
echo "  Done!"

echo "  Deleting existing code..."
ssh root@157.230.137.192 "rm -rf /root/go/src/lenslockedbr.com"
echo "  Code deleted successfully!"

echo "  Uploading code..."
# The \ at the end of the line tells bash that our command isn't done
# and wraps to the next line.
rsync -avr --exclude '.git/*' --exclude 'tmp/*' --exclude 'images/*' ./ root@157.230.137.192:/root/go/src/lenslockedbr.com/
#echo "  Code uploaded successfully!"

echo "  Go getting deps..."
ssh root@157.230.137.192 "export GOPATH=/root/go; /usr/local/go/bin/go get golang.org/x/crypto/bcrypt"

ssh root@157.230.137.192 "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/gorilla/mux"

ssh root@157.230.137.192 "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/gorilla/schema"

ssh root@157.230.137.192 "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/gorilla/csrf"

ssh root@157.230.137.192 "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/lib/pq"

ssh root@157.230.137.192 "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/jinzhu/gorm"

echo "  Building the code on remote server..."
ssh root@157.230.137.192 "export GOPATH=/root/go; cd /root/app; /usr/local/go/bin/go build -o ./server $GOPATH/src/lenslockedbr.com/*.go"
echo "  Code built successfully!..."


echo "  Moving assets..."
ssh root@157.230.137.192 "cd /root/app; cp -R /root/go/src/lenslockedbr.com/assets ."
echo "  Assets moved successfully!..."

echo "  Moving views..."
ssh root@157.230.137.192 "cd /root/app; cp -R /root/go/src/lenslockedbr.com/views ."
echo "  Views moved successfully!..."

echo "  Moving Caddyfile..."
ssh root@157.230.137.192 "cd /root/app; cp /root/go/src/lenslockedbr.com/Caddyfile ."
echo "  Caddyfile moved successfully!..."

echo "  Restarting the server..."
ssh root@157.230.137.192 "service leandr0.net restart"
echo "  Server restarted successfully!..."

echo "  Restarting Caddy server..."
ssh root@157.230.137.192 "service caddy restart"
echo "  Caddy restarted successfully!..."

echo "===== Done releasing lenslockedbr.com ====="

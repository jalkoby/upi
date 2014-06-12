upi
===

Simple file hosting build on top of Martini(Go lang) 

== Setup process

~ git clone git@github.com:jalkoby/upi.git
~ cd upi
~ vim .envrc
# content of .envrc
export GOPATH="$(pwd)/.godeps"
export UPI_PASSWORD="*password*"
export AWS_ACCESS_KEY_ID="*aws key*"
export AWS_SECRET_ACCESS_KEY=""
export UPI_UPLOAD="$(pwd)/uploads"
export UPI_PUBLIC_HOST="http://localhost:3000"
export AWS_BUCKET="****"
# content of .envrc
~ direnv allow
~ gpm install
~ go get github.com/codegangsta/gin
~ gin

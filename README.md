## Get all dependencies

    go get -v

## DB Migrations
#### Install migration tool goose

    go get -v goose

### Run migration script

	goose -dir migrations mysql "$DBUSER:$DBPASSWORD@/user_sso" up

## Google API credential
    $ touch creds.json

Add your Google API application Id and secret in creds.json with format

    {
        "cid": "{$APPLICATIONID}",
        "csecret" : "{$APPLICATIONSECRET}"
    }

## Run Application

    go run main.go -dbUrl="$DBUSER:$DBPASSWORD@tcp(127.0.0.1:3306)/user_sso" -port="9090" -baseUrl="http://localhost"
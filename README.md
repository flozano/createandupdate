# createandupdate

Stupid, non-idiomatic, low-quality stuff.

Tool to register a user then create and update objects one by one repeatedly
and infinitely for an App in Kii Cloud.

## Usage

To run without install:

    $ go run createandupdate.go {URL} {APP_ID} {APP_KEY}

or with install:

    $ go get github.com/flozano/createandupdate
    $ createandupdate {URL} {APP_ID} {APP_KEY}

Where `{URL}` should be replaced by endpoint of Kii Cloud, like
`https://api.kii.com` without last slash (`/`).

`{APP_ID}` and `{APP_KEY}` should be replaced by information of your app.

Don't forget to stop after the test, otherwise this tool will create too many
records on Kii Cloud.

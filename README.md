# Gypsy Cab

Gypsy Cab is a standalone image processing web-service written in Go. In short, send images at it, it will perform configurable processing recipes, and finally upload it to S3. It uses RethinkDB as a datastore.

A work in progress still... but I use it in production for a handful of apps.

Some of the awesome Go libraries used:

* [go-martini/martini](http://github.com/go-martini/martini)
* [disintegration/gift](http://github.com/disintegration/gift)
* [dancannon/gorethink](http://github.com/dancannon/gorethink)

### TODO

* Validate unique job `Key` per user
* Validate unique item ids
* Test coverage `JobProcess` API endpoint
* Cleanup tmp files after success
* Document how to use install it!
* Document how to use use it!
* Document how to use deploy it!
* Publish Objective-C client
* Publish web client

### Basic Use

`curl -F image='@test.jpg' http://localhost:3000/api/v1/jobs/JOB-IDENTIFIER/STORE-TOKEN`

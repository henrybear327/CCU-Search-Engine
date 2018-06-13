# RSS 

Manager, downloader, and parser

# TODO

* Fetcher   
    * Fetch periodically
    * Able to handle restart
    * Fetch the external link
    * Respect the update delay
* Storage
    * Duplication detection

# Packages

* `go get github.com/mmcdole/gofeed`: RSS 

# CLI

A command line interface for managing RSS subscriptions

# TODO

* Design database schema
    * multi-user
* Design command line argument 
    * `-showallfeeds`
    * `-adduser`, `-deluser`
* User
    * add user `-adduser [username]`
    * delete user (and it's subscriptions) `-deluser [username]`
* Subscription 
    * show all unique feeds in database `-showallfeeds *`
    * feed (for a certain user)
        * add `-user [username] -addfeed [rss link]`
            * deplicated feed checking
            * default group "default"
        * delete `-user [username] -delfeed [link]`
            * search all groups
        * show all `-user [username] -showfeeds *`
* OPML parser (from sites like Feedly)
    * Auto-import to rss feed under certain user `-user [username] -importOPML [path to file]`
        * auto-create/merge groups
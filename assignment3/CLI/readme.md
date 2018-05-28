# CLI

A command line interface for managing RSS subscriptions

# TODO

* Design database schema
    * multi-user
    * multi-news-group
* User
    * add user
    * delete user (and it's subscriptions)
* Subscription (for user)
    * feed
        * add
            * deplicated feed checking
        * delete
        * show all
    * group
        * create
        * delete
        * add feed
        * delete feed
        * show all
* OPML parser (from sites like Feedly)
    * Auto-import to rss feed under certain user
        * auto-create/merge groups
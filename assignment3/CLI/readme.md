# CLI

A command line interface for managing RSS subscriptions

# TODO

* Design database schema
    * multi-user
    * multi-news-group
* Design command line argument 
    * `-showallfeeds`
    * `-adduser`, `-deluser`
    * has `-user`
        * `-newgroup`, `-delgroup`
        * has `-group`
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
    * group (for a certain user) 
        * create `-user [username] -newgroup [group name]`
        * delete `-user [username] -delgroup [group name]`
        * add feed `-user [username] -group [group name] -addfeed [feedname]`
        * delete feed `-user [username] -group [group name] -delfeed [feedname]`
        * show all `-user [username] -group [group name] -showfeeds *`
* OPML parser (from sites like Feedly)
    * Auto-import to rss feed under certain user `-user [username] -importOPML [path to file]`
        * auto-create/merge groups
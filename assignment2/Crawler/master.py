import configparser
import datetime
import sys

import URLManager
import fetcher

if __name__ == '__main__':
    config = configparser.ConfigParser()
    config.read('crawler.config')

    print("Master started. Initial page {}".format(config["DEFAULT"]["initial_page"]))

    # init
    start_time = datetime.datetime.now()

    url_manager = URLManager.URLManager()
    url_manager.insert_url(config["DEFAULT"]["initial_page"], 0, 0)

    fetcher = fetcher.Fetcher(url_manager)

    end_time = datetime.datetime.now()
    delta = end_time - start_time
    print("Init time " + str(delta))

    # start crawling
    while url_manager.has_next_url():  # TODO: change to empty over 10 seconds
        print("queue size", url_manager.get_size())
        next_url = url_manager.get_next_url() # TODO: handle "" case

        print("fetching ", next_url)
        fetcher.get_page(next_url)
        print("done")

        sys.stdout.flush()
        sys.stderr.flush()

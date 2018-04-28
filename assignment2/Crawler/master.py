import URLManager
import fetcher

import datetime

# Parameters

# initial_page = "https://www.google.com"
initial_page = "https://www.npr.org/"
checking_url = "npr.org"

max_retry = 3

if __name__ == '__main__':
    print("Master started. Initial page {}".format(initial_page))

    # init
    start_time = datetime.datetime.now()

    url_manager = URLManager.URLManager(max_retry)
    url_manager.insert_new_urls([initial_page])

    fetcher = fetcher.Fetcher(checking_url, url_manager)

    end_time = datetime.datetime.now()
    delta = end_time - start_time
    print("Init time " + str(delta))

    # start crawling
    while url_manager.has_next_url():  # TODO: change to empty over 10 seconds
        print("queue size", url_manager.get_size())
        next_url = url_manager.get_next_url()

        print("fetching ", next_url)
        fetcher.get_page(next_url)
        print("done")

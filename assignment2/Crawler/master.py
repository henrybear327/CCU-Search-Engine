import URLQueue
import fetcher

import datetime

# initial_page = "https://www.google.com"
initial_page = "https://www.npr.org/"
checking_url = "npr.org"

if __name__ == '__main__':
    print("Master started. Initial page {}".format(initial_page))

    # init
    start_time = datetime.datetime.now()

    url_queue = URLQueue.URLQueue()
    url_queue.insert_url([initial_page])

    fetcher = fetcher.Fetcher(checking_url, url_queue)

    end_time = datetime.datetime.now()
    delta = end_time - start_time
    print("Init time " + str(delta))

    # start crawling
    while url_queue.has_next_url():  # TODO: change to empty over 10 seconds
        print("queue size", url_queue.get_size())
        next_url = url_queue.get_next_url()

        print("fetching ", next_url)
        fetcher.get_page(next_url)
        print("done")

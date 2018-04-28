import queue
from collections import namedtuple
import sys


class URLManager:
    def __init__(self, max_retry):
        self.url_queue = queue.Queue()
        self.in_queue = set()  # urls in queue
        self.fetched = set()  # urls fetched

        self.queueData = namedtuple('QueueData', ['url', 'attempts', 'level'])

        self.max_retry = max_retry

    def __del__(self):
        # save data
        pass

    def has_next_url(self):
        return self.url_queue.empty() is False

    def get_next_url(self):
        if self.has_next_url() is False:
            return ""
        return self.url_queue.get()

    def add_fetched_url(self, url):
        self.fetched.add(url)

        # remove it from the set now, so we can avoid url being added back to queue during parallel fetching...
        self.in_queue.discard(url)

    def insert_url(self, url, attempts, level):
        # check for in queue or not
        if url in self.in_queue or url in self.fetched:
            return
        if attempts >= self.max_retry:
            sys.stderr.write("Max retries exceeded " + url + "\n")
            return

        self.in_queue.add(url)

        # enqueue
        self.url_queue.put(self.queueData(url, attempts, level))

    def insert_new_urls(self, urls, level):
        for url in urls:
            # check for in queue or not
            if url in self.in_queue or url in self.fetched:
                continue
            self.in_queue.add(url)

            # enqueue
            self.url_queue.put(self.queueData(url, 0, level))

    def get_size(self):
        return self.url_queue.qsize()

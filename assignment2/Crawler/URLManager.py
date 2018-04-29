import configparser
import queue
import sys
from collections import namedtuple


class URLManager:
    def __init__(self):
        self.url_queue = queue.Queue()
        self.in_queue = set()  # urls in queue
        self.fetched = set()  # urls fetched

        config = configparser.ConfigParser()
        config.read('crawler.config')
        self.max_retry = int(config["RULES"]["max_retry"])
        self.assumed_non_content_depth = int(config["RULES"]["assumed_non_content_depth"])
        self.max_overall_depth = int(config["RULES"]["max_overall_depth"])

        # fetched_set_file = config["STORAGE"]["fetched_set_file"]

        # with open(fetched_set_file, "w+") as inputFile:
        #     pass
        #
        # with open(fetched_set_file, "r") as inputFile:
        #     for line in inputFile:
        #         line = line.replace("\n", "")
        #         if line == "":
        #             continue
        #         self.fetched.add(line)
        #         print("Add to fetched set", line)
        #
        # self.fetchedFile = open(fetched_set_file, 'a')

        self.queueData = namedtuple('QueueData', ['url', 'attempts', 'depth'])

    def __del__(self):
        # self.fetchedFile.close()
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

        if url.depth < self.assumed_non_content_depth:
            return
        line = url.url + "\n"
        # self.fetchedFile.write(line)
        print("add fetched url", url)

    def add_retry_url(self, url, attempts, depth):
        self.in_queue.discard(url)
        self.insert_url(url, attempts, depth)

    def insert_url(self, url, attempts, depth):
        # check for in queue or not
        if url in self.in_queue or url in self.fetched:
            return
        if attempts >= self.max_retry:
            sys.stderr.write("Max retries exceeded " + url + "\n")
            return
        if depth >= self.max_overall_depth:
            sys.stderr.write("Max overall depth exceeded " + url + "\n")
            return

        self.in_queue.add(url)

        # enqueue
        self.url_queue.put(self.queueData(url, attempts, depth))

    def insert_new_urls(self, urls, depth):
        for url in urls:
            self.insert_url(url, 0, depth)

    def get_size(self):
        return self.url_queue.qsize()

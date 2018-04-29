import configparser
import queue
import sys
from collections import namedtuple

"""
Decides what url can go into queue 

Depth management here
"""


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
        self.max_internal_depth = int(config["RULES"]["max_internal_depth"])
        self.max_external_depth_ = int(config["RULES"]["max_external_depth_"])
        self.checking_url = config["SITE"]["checking_url"]

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

        self.queueData = namedtuple('QueueData', ['url', 'attempts', 'internal_depth', 'external_depth'])

    def __del__(self):
        # self.fetchedFile.close()
        pass

    def is_current_site_url(self, url):
        if str(url).find(self.checking_url) != -1:
            return True
        else:
            # print("External site url ", url)
            return False

    def has_next_url(self):
        return self.url_queue.empty() is False

    def get_next_url(self):
        if self.has_next_url() is False:
            return ""
        return self.url_queue.get()

    def add_fetched_url(self, url):
        self.fetched.add(url.url)

        # remove it from the set now, so we can avoid url being added back to queue during parallel fetching...
        self.in_queue.discard(url.url)

        # when restarting, crawler just need to focus on topic pages and start from there
        # if url.depth < self.assumed_non_content_depth:
        #     return
        # line = url.url + "\n"
        # self.fetchedFile.write(line)

        print("add fetched url", url)

    def add_retry_url(self, url):
        self.in_queue.discard(url.url)
        self.insert_url(url.url, url.attempts + 1, url.internal_depth, url.external_depth)

    def insert_url(self, url, attempts, internal_depth, external_depth):
        # check for in queue or not
        if url in self.in_queue or url in self.fetched:
            return
        if attempts >= self.max_retry:
            sys.stderr.write("Max retries exceeded " + url + "\n")
            return

        if internal_depth + external_depth >= self.max_overall_depth:
            sys.stderr.write("Max overall depth exceeded " + url + "\n")
            return
        if internal_depth >= self.max_internal_depth:
            sys.stderr.write("Max internal depth exceeded " + url + "\n")
            return
        if external_depth >= self.max_external_depth_:
            sys.stderr.write("Max external depth exceeded " + url + "\n")
            return

        self.in_queue.add(url)

        # enqueue
        self.url_queue.put(self.queueData(url, attempts, internal_depth, external_depth))

    def insert_new_urls(self, urls, parent_url):
        for url in urls:
            internal_depth = parent_url.internal_depth
            external_depth = parent_url.external_depth
            if self.is_current_site_url(url):
                internal_depth += 1
            else:
                external_depth += 1

            self.insert_url(url, 0, internal_depth, external_depth)

    def get_size(self):
        return self.url_queue.qsize()

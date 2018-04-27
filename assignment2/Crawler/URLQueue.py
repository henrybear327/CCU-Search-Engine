import queue


class URLQueue:
    def __init__(self):
        self.url_queue = queue.Queue()
        self.seen = set()

    def has_next_url(self):
        return self.url_queue.empty() is False

    def get_next_url(self):
        if self.has_next_url() is False:
            return ""
        return self.url_queue.get()

    def insert_url(self, urls):
        for url in urls:
            # check for fetched or not
            if url in self.seen:
                continue
            self.seen.add(url)

            # enqueue
            self.url_queue.put(url)

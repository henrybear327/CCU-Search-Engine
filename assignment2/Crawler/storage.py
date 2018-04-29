import configparser
import datetime
import pprint
import hashlib

from pymongo import MongoClient


class Storage:
    def __init__(self):
        """
        Database name = npr_test
        Collection name = data
        """
        config = configparser.ConfigParser()
        config.read('crawler.config')
        self.path = config["FOLDER"]["page_source"]

        self.database_name = config["STORAGE"]["database_name"]
        self.collection_name = config["STORAGE"]["collection_name"]

        self.client = MongoClient('localhost', 27017)
        self.db = self.client[self.database_name]
        self.collection = self.db[self.collection_name]

    def get_sha1(self, url: str):
        hash = hashlib.sha1()
        hash.update(url.encode('utf-8'))
        url_sha1 = hash.hexdigest()
        print(url_sha1)

        return url_sha1

    def insert_record(self, url: str, title: str, page_source: str):
        url_sha1 = self.get_sha1(url)

        record = {
            "url": url,
            "title": title,
            "page_source": page_source,
            "timestamp": datetime.datetime.now(),
            "url_sha1": url_sha1
        }

        post_id = self.collection.insert_one(record).inserted_id
        print("Inserted a record", post_id)

        self.write_source_code_to_file(url_sha1, page_source)

    def display_all_records(self):
        print("displaying all records")
        for record in self.collection.find():
            pprint.pprint(record)

    def search_record(self, url):
        print("search url", url)

        url_sha1 = self.get_sha1(url)
        ret = self.collection.find({"url_sha1": url_sha1})
        count = ret.count()
        # for record in ret:
        #     pprint.pprint(record['timestamp'])
        for record in self.collection.find({"url": url}):
            pprint.pprint(record)
        print(count, "record(s)")

    def clear_collection(self):
        print("clearing collection")
        self.db.drop_collection(self.collection_name)

        # TODO: remove data in folder

    def write_source_code_to_file(self, url, page_source):
        with open(self.path + url + ".html", "w") as outputFile:
            outputFile.write(page_source)

    def test_run(self):
        storage.insert_record("apple.com", "apple", "<html>")
        storage.insert_record("google.com", "google", "<html>")
        storage.display_all_records()
        storage.search_record("apple.com")
        storage.search_record("google.com")
        storage.clear_collection()
        storage.display_all_records()


if __name__ == '__main__':
    storage = Storage()
    storage.clear_collection()
    storage.display_all_records()
    # storage.search_record("https://www.npr.org/")

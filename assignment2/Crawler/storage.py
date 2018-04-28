import configparser
import datetime
import pprint

from pymongo import MongoClient


class Storage:
    def __init__(self):
        """
        Database name = npr_test
        Collection name = data
        """
        self.client = MongoClient('localhost', 27017)
        self.db = self.client.npr_test
        self.collection = self.db.data

        config = configparser.ConfigParser()
        config.read('crawler.config')
        self.path = config["FOLDER"]["page_source"]

    def insert_record(self, url, title, page_source):
        record = {
            "url": url,
            "title": title,
            "page_source": page_source,
            "date": datetime.datetime.now()
        }

        post_id = self.collection.insert_one(record).inserted_id
        print("Inserted a record", post_id)

        self.write_source_code_to_file(url, page_source)

    def display_all_records(self):
        print("displaying all records")
        for record in self.collection.find():
            pprint.pprint(record)

    def search_record(self, url):
        for record in self.collection.find({"url": url}):
            pprint.pprint(record)

    def clear_collection(self):
        print("clearing collection")
        self.db.drop_collection("data")

        # TODO: remove data in folder

    def write_source_code_to_file(self, url, page_source):
        with open(self.path + url + ".html", "w") as outputFile:
            outputFile.write(page_source)


if __name__ == '__main__':
    storage = Storage()
    storage.insert_record("apple.com", "apple", "<html>")
    storage.display_all_records()
    storage.clear_collection()
    storage.display_all_records()

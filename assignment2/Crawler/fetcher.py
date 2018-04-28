from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.common.exceptions import NoSuchElementException, TimeoutException

from urllib.parse import urlparse

import datetime
import parser
import URLManager

import sys


class Fetcher:
    def __init__(self, checking_url, url_manager: URLManager):
        """

        :param checking_url:
        :param url_manager:
        """
        chrome_options = Options()
        chrome_options.add_argument("--headless")
        chrome_options.add_argument("--window-size=1920x1080")
        self.driver = webdriver.Chrome(chrome_options=chrome_options)
        self.parser = parser.Parser(checking_url, url_manager)
        self.url_manager = url_manager

    def __del__(self):
        self.driver.quit()

    def get_page(self, url):
        # get content
        start_time = datetime.datetime.now()

        try:
            self.driver.get(url.url)
        except TimeoutException:
            sys.stderr.write("Timeout " + url.url + "\n")
            self.url_manager.insert_url(url.url, url.attempts + 1, url.level)
            return

        filename = "{}-{}.png".format(datetime.datetime.today(), self.driver.title)
        self.driver.save_screenshot("../image/" + filename)

        # parse links
        try:
            links = self.driver.find_elements_by_tag_name('a')
        except NoSuchElementException:
            sys.stderr.write(url.url + " has no links\n")
            return

        end_time = datetime.datetime.now()
        delta = end_time - start_time
        print("get content", delta)

        self.url_manager.add_fetched_url(url)
        self.parser.parse(self.driver.page_source, links, url.level + 1)

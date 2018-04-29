import configparser
import datetime
import sys
import requests
from requests.exceptions import Timeout, ConnectTimeout
from selenium import webdriver
from selenium.common.exceptions import NoSuchElementException, TimeoutException
from selenium.webdriver.chrome.options import Options

import URLManager
import parser
import platform


class Fetcher:
    def __init__(self, url_manager: URLManager):
        config = configparser.ConfigParser()
        config.read('crawler.config')
        self.backend = config["RULES"]["backend"]
        self.timeout = float(config["RULES"]["timeout"])

        self.parser = parser.Parser(url_manager)
        self.url_manager = url_manager

        if self.backend == "chrome":
            chrome_options = Options()
            chrome_options.add_argument("--headless")
            chrome_options.add_argument("--window-size=1920x1080")
            chrome_options.add_argument("--proxy=null")

            self.driver = webdriver.Chrome(chrome_options=chrome_options)
        elif self.backend == "requests":
            pass
        else:
            print("unknown backend")
            sys.exit(1)

    def __del__(self):
        self.driver.quit()

    def get_page(self, url):
        # get content
        start_time = datetime.datetime.now()
        if self.backend == "chrome":
            try:
                self.driver.get(url.url)
            except TimeoutException as e:
                sys.stderr.write("Timeout " + url.url + "\n")
                sys.stderr.write(e.msg)
                self.url_manager.add_retry_url(url.url, url.attempts + 1, url.level)
                return

            filename = "{}-{}.png".format(datetime.datetime.today(), self.driver.title)
            self.driver.save_screenshot("../image/" + filename)

            # # parse links (selenium parser)
            # try:
            #     selenium_links = self.driver.find_elements_by_tag_name('a')
            # except NoSuchElementException as e:
            #     sys.stderr.write(url.url + " has no links\n")
            #     sys.stderr.write(e.msg)
            #     return

            end_time = datetime.datetime.now()
            delta = end_time - start_time
            print("get content", delta)

            self.url_manager.add_fetched_url(url)
            # self.parser.parse(url.url, self.driver.title, self.driver.page_source, url.level + 1, links=selenium_links)
            self.parser.parse(url.url, self.driver.title, self.driver.page_source, url.level + 1, links=None)
        elif self.backend == "requests":
            try:
                r = requests.get(url.url, timeout=self.timeout)
            except Timeout as e:
                sys.stderr.write("Timeout " + url.url + "\n")
                self.url_manager.add_retry_url(url.url, url.attempts + 1, url.level)
                return

            end_time = datetime.datetime.now()
            delta = end_time - start_time
            print("get content", delta)

            self.url_manager.add_fetched_url(url)
            self.parser.parse(url.url, "", r.text, url.level + 1, links=None)

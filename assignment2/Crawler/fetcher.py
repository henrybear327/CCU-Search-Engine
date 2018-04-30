import configparser
import datetime
import sys
import requests
from requests.exceptions import Timeout, InvalidSchema
from selenium import webdriver
from selenium.common.exceptions import NoSuchElementException, TimeoutException
from selenium.webdriver.chrome.options import Options

import URLManager
import parser

"""
Simply fetch the data
"""


class Fetcher:
    def __init__(self, url_manager: URLManager):
        config = configparser.ConfigParser()
        config.read('crawler.config')
        self.backend = config["RULES"]["backend"]

        self.parser = parser.Parser(url_manager)
        self.url_manager = url_manager

        self.timeout = float(config["RULES"]["timeout"])

        chrome_options = Options()
        chrome_options.add_argument("--headless")
        chrome_options.add_argument("--window-size=1920x1080")
        chrome_options.add_argument("--proxy=null")

        self.driver = webdriver.Chrome(chrome_options=chrome_options)

    def __del__(self):
        self.driver.quit()

    def is_response_200(self, url):
        try:
            r = requests.get(url.url, timeout=self.timeout)

            if r.status_code != 200:
                sys.stderr.write("Page access error " + url.url + "\n")
                return False
        except Timeout:
            sys.stderr.write("Timeout " + url.url + "\n")
            self.url_manager.add_retry_url(url)
            return False
        except:
            print("Unexpected requests.get() error:", sys.exc_info()[0])
            return False
        return True

    def get_page(self, url):
        start_time = datetime.datetime.now()
        if url.external_depth > -1 or self.backend == "chrome":
            print("Use chrome headless")

            if not self.is_response_200(url):
                return

            try:
                self.driver.get(url.url)
            except TimeoutException as e:
                sys.stderr.write("Timeout " + url.url + "\n")
                sys.stderr.write(e.msg)
                self.url_manager.add_retry_url(url)
                return
            except:
                print("Unexpected self.driver.get() error:", sys.exc_info()[0])
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
            # self.parser.parse(self.driver.title, self.driver.page_source, url, links=selenium_links)
            self.parser.parse(self.driver.title, self.driver.page_source, url, links=None)
        elif self.backend == "requests":
            print("Use requests")

            try:
                r = requests.get(url.url, timeout=self.timeout)

                if r.status_code != 200:
                    sys.stderr.write("Page access error " + url.url + "\n")
                    return
            except Timeout:
                sys.stderr.write("Timeout " + url.url + "\n")
                self.url_manager.add_retry_url(url)
                return
            except InvalidSchema:
                sys.stderr.write("invalid link " + url.url + "\n")
                return
            except:
                print("Unexpected requests.get() error:", sys.exc_info()[0])
                return

            end_time = datetime.datetime.now()
            delta = end_time - start_time
            print("get content", delta)

            self.url_manager.add_fetched_url(url)
            self.parser.parse("", r.text, url, links=None)
        else:
            print("What the fuck to fetch with?")
            sys.exit(1)

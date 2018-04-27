from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.common.exceptions import NoSuchElementException
# from collections import namedtuple

from urllib.parse import urlparse

import datetime
import parser


class Fetcher:
    def __init__(self, checking_url, url_queue):
        chrome_options = Options()
        chrome_options.add_argument("--headless")
        chrome_options.add_argument("--window-size=1920x1080")
        self.driver = webdriver.Chrome(chrome_options=chrome_options)

        # self.FetchedData = namedtuple('FetchedData', ['page_source', 'title'])

        self.parser = parser.Parser(checking_url, url_queue)

    def __del__(self):
        self.driver.quit()

    def get_page(self, url):
        # get content
        start_time = datetime.datetime.now()

        self.driver.get(url)

        filename = "{}-{}.png".format(datetime.datetime.today(), self.driver.title)
        self.driver.save_screenshot("../image/" + filename)

        # parse links
        try:
            links = self.driver.find_elements_by_tag_name('a')
        except NoSuchElementException:
            return

        end_time = datetime.datetime.now()
        delta = end_time - start_time
        print("get content", delta)

        self.parser.parse(self.driver.page_source, links)

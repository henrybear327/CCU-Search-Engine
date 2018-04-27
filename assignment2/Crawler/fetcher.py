from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.common.exceptions import NoSuchElementException, StaleElementReferenceException
# from collections import namedtuple

from urllib.parse import urlparse

import datetime
import parser


class Fetcher:
    def __init__(self):
        chrome_options = Options()
        chrome_options.add_argument("--headless")
        chrome_options.add_argument("--window-size=1920x1080")
        self.driver = webdriver.Chrome(chrome_options=chrome_options)

        # self.FetchedData = namedtuple('FetchedData', ['page_source', 'title'])

        self.parser = parser.Parser()

    def __del__(self):
        self.driver.quit()

    def split_url_parameters(self, href):
        url = urlparse(href)
        # print(url.scheme, url.netloc, url.path, url.params, url.query)
        return url

    def get_new_links(self, url):
        # get content
        start_time = datetime.datetime.now()

        self.driver.get(url)
        # driver.save_screenshot(driver.title+".png")

        end_time = datetime.datetime.now()
        delta = end_time - start_time
        print("get content", delta)

        # parse links
        start_time = datetime.datetime.now()

        soup_links = self.parser.get_all_links(self.driver.page_source)

        selenium_links = []
        try:
            links = self.driver.find_elements_by_tag_name('a')
        except NoSuchElementException:
            return selenium_links

        for link in links:
            try:
                if link.get_attribute("href") is None:
                    continue
            except StaleElementReferenceException:
                continue

            # print(link.text)
            split_href = self.split_url_parameters(link.get_attribute("href"))
            if split_href.netloc == "":
                continue  # void(0) case

            href = split_href.scheme + "://" + split_href.netloc + split_href.path
            # print(link.text, href)
            selenium_links.append(href)

        end_time = datetime.datetime.now()
        delta = end_time - start_time
        print("get links", delta)

        print(len(soup_links), len(selenium_links))
        return selenium_links
        # return self.FetchedData(self.driver.page_source, self.driver.title)


if __name__ == '__main__':
    fetcher = Fetcher()
    print(fetcher.get_new_links("https://edition.cnn.com/"))

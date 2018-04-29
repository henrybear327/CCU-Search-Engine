import configparser
import datetime
import sys
from urllib.parse import urljoin
from urllib.parse import urlparse

from bs4 import BeautifulSoup
from selenium.common.exceptions import StaleElementReferenceException

import URLManager
import storage

"""
Generate links for queue
"""


class Parser:
    def __init__(self, url_manager: URLManager):
        config = configparser.ConfigParser()
        config.read('crawler.config')

        self.url_manager = url_manager
        self.storage = storage.Storage()

    def parse(self, url, title, page_source, depth, links=None):
        """
        1. generate new links
        2. TODO: content extraction
        :param url:
        :param title:
        :param page_source:
        :param depth:
        :param links:
        :return:
        """
        # new_links = self.get_all_links(url, links)
        new_links_soup = self.get_all_links_soup(url, page_source)
        # print("selenium", len(new_links), "soup", len(new_links_soup))

        # sys.stderr.write("============================\n")
        # sys.stderr.write(url + "\n")
        # for link in new_links_soup:
        #     if link not in new_links:
        #         sys.stderr.write("Soup has " + link + "\n")
        # for link in new_links:
        #     if link not in new_links_soup:
        #         sys.stderr.write("Selenium has " + link + "\n")
        # sys.stderr.write("============================\n")

        if len(new_links_soup) > 0:
            # self.url_manager.insert_new_urls(new_links, depth)
            self.url_manager.insert_new_urls(new_links_soup, depth)

            self.storage.insert_record(url, title, page_source)

    def get_all_links(self, base_url, links):
        start_time = datetime.datetime.now()

        selenium_links = []
        for link in links:
            try:
                if link.get_attribute("href") is None:
                    continue
            except StaleElementReferenceException as e:
                sys.stderr.write(e.msg)
                continue

            # print(link.text)
            split_href = urlparse(link.get_attribute("href"))
            if split_href.netloc == "":
                continue  # void(0) case

            href = split_href.scheme + "://" + split_href.netloc + split_href.path
            # print(href)
            # print(link.text, split_href.scheme, split_href.netloc, split_href.path)
            if href != base_url:
                selenium_links.append(href)

        end_time = datetime.datetime.now()
        delta = end_time - start_time
        print("get selenium links", delta)

        # print(selenium_links)
        return selenium_links

    def get_all_links_soup(self, base_url, page_source):
        start_time = datetime.datetime.now()

        # original_url = self.trim_trailing_slash(base_url)
        #
        # base_url = urlparse(base_url)
        # base_url = self.trim_trailing_slash(base_url.scheme + "://" + base_url.netloc)

        soup = BeautifulSoup(page_source, "lxml")
        links = soup.find_all('a', href=True)

        result = []
        for link in links:
            href = str(link['href']).strip()
            if href == "void(0)":
                continue
            if href.startswith("mailto://"):
                continue

            href = urljoin(base_url, href)

            url = urlparse(href)
            href = url.scheme + "://" + url.netloc + url.path

            if href != base_url:
                result.append(href)
                # print(href)

        # print(len(links))
        end_time = datetime.datetime.now()
        delta = end_time - start_time
        print("get soup links", delta)

        return result

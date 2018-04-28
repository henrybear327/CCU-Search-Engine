from bs4 import BeautifulSoup
from urllib.parse import urlparse
from selenium.common.exceptions import StaleElementReferenceException
import datetime
import URLManager


class Parser:
    def __init__(self, checking_url, url_manager: URLManager):
        self.checking_url = checking_url
        self.url_manager = url_manager

    def parse(self, page_source, links, level):
        new_links = self.get_all_links(links)
        self.url_manager.insert_new_urls(new_links, level)

    def split_url_parameters(self, href):
        url = urlparse(href)
        # print(url.scheme, url.netloc, url.path, url.params, url.query)
        return url

    def get_all_links(self, links):
        start_time = datetime.datetime.now()

        # soup_links = self.parser.get_all_links(self.driver.page_source)

        selenium_links = []
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
            # print(href)
            # print(link.text, split_href.scheme, split_href.netloc, split_href.path)
            if str(href).find(self.checking_url) != -1:
                selenium_links.append(href)
            else:
                print("Rejected url ", href)

        end_time = datetime.datetime.now()
        delta = end_time - start_time
        print("get links", delta)

        # print("Two methods link count check", len(soup_links), len(selenium_links))
        return selenium_links

    # def get_all_links(self, page_source):
    #     soup = BeautifulSoup(page_source, "lxml")
    #     links = soup.find_all('a', href=True)
    #     for link in links:
    #         href = link['href']
    #         print(href)
    #         link_name = link.text
    #
    #         url = urlparse(href)
    #         print(link_name, url.scheme, url.netloc, url.path, url.params, url.query)
    #         print(url)
    #     return links

    def get_page_content(self, page_source):
        pass

from bs4 import BeautifulSoup
from urllib.parse import urlparse
from selenium.common.exceptions import StaleElementReferenceException
import datetime
import URLManager
import sys


class Parser:
    def __init__(self, checking_url, url_manager: URLManager):
        self.checking_url = checking_url
        self.url_manager = url_manager

    def parse(self, url, page_source, links, level):
        # new_links = self.get_all_links(url, links)
        new_links_soup = self.get_all_links_soup(url, page_source)
        # print("selenium", len(new_links), "soup", len(new_links_soup))
        #
        # sys.stderr.write("============================\n")
        # sys.stderr.write(url + "\n")
        # for link in new_links_soup:
        #     if link not in new_links:
        #         sys.stderr.write("Soup has " + link + "\n")
        # for link in new_links:
        #     if link not in new_links_soup:
        #         sys.stderr.write("Selenium has " + link + "\n")
        # sys.stderr.write("============================\n")
        #
        # self.url_manager.insert_new_urls(new_links, level)
        self.url_manager.insert_new_urls(new_links_soup, level)

    def split_url_parameters(self, href):
        url = urlparse(href)
        # print(url.scheme, url.netloc, url.path, url.params, url.query)
        return url

    def is_current_site_url(self, url):
        if str(url).find(self.checking_url) != -1:
            return True
        else:
            print("Rejected url ", url)
            return False

    def trim_trailing_slash(self, url):
        url = str(url)
        if url.endswith("/"):  # trim trailing /
            url = url[0:-1]
        return url

    def get_all_links(self, base_url, links):
        start_time = datetime.datetime.now()

        base_url = self.trim_trailing_slash(base_url)

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

            path = self.trim_trailing_slash(split_href.path)

            href = split_href.scheme + "://" + split_href.netloc + path
            # print(href)
            # print(link.text, split_href.scheme, split_href.netloc, split_href.path)
            if self.is_current_site_url(href) and href != base_url:
                selenium_links.append(href)

        end_time = datetime.datetime.now()
        delta = end_time - start_time
        print("get selenium links", delta)

        # print(selenium_links)
        return selenium_links

    def get_all_links_soup(self, base_url, page_source):
        start_time = datetime.datetime.now()

        original_url = self.trim_trailing_slash(base_url)

        base_url = urlparse(base_url)
        base_url = self.trim_trailing_slash(base_url.scheme + "://" + base_url.netloc)

        soup = BeautifulSoup(page_source, "lxml")
        links = soup.find_all('a', href=True)

        result = []
        for link in links:
            href = str(link['href'])
            # print(href)

            """
            Special cases
            1. <a href="#mainContent" class="skiplink">Skip to main content</a>
            2. <a href=".">
            3. <a href="/book" class="nprhome nprhome-news" data-metrics-action="click npr logo">
            4. <a href="book">
            5. <a href="https://google.com>
            """
            if href.startswith("#") or href.startswith("."):  # case 1, 2
                continue

            url = urlparse(href)
            # print(link.text)
            # ParseResult(scheme='https', netloc='www.npr.org', path='/sections/allsongs/2018/04/27/606066039/janelle-mon-e-strips-the-hardware-for-humanity', params='', query='', fragment='')
            # print(url)

            netloc = self.trim_trailing_slash(url.netloc)
            path = self.trim_trailing_slash(url.path)

            if href.startswith("http"):  # case 5
                # get rid of param
                href = url.scheme + "://" + netloc + path
            elif path != "" and not path.startswith("/"):  # case 4, use current url + path
                path = "/" + path
                href = original_url + path
            elif url.scheme == "" or (path != "" and path.startswith("/")):  # case 3, use base url + path
                href = base_url + path
            else:
                print("parsing failed", link, url)

            if self.is_current_site_url(href) and href != original_url:
                result.append(href)
                # print(href)

        # print(len(links))
        end_time = datetime.datetime.now()
        delta = end_time - start_time
        print("get soup links", delta)

        return result

    def get_page_content(self, page_source):
        pass

from bs4 import BeautifulSoup
from urllib.parse import urlparse


class Parser:
    def __init__(self):
        pass

    def get_all_links(self, page_source):
        soup = BeautifulSoup(page_source, "lxml")
        links = soup.find_all('a', href=True)
        for link in links:
            href = link['href']
            link_name = link.text

            url = urlparse(href)
            # print(link_name, url.scheme, url.netloc, url.path, url.params, url.query)
            print(url)
        return links

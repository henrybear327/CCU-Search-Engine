from selenium import webdriver
from selenium.webdriver.chrome.options import Options
# from collections import namedtuple
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

    def get_new_links(self, url):
        # get content
        start_time = datetime.datetime.now()

        self.driver.get(url)
        # driver.save_screenshot(driver.title+".png")

        end_time = datetime.datetime.now()
        delta = end_time - start_time
        print("get content " + str(delta))

        # parse data
        links = self.parser.get_all_links(self.driver.page_source)

        return links
        # return self.FetchedData(self.driver.page_source, self.driver.title)


if __name__ == '__main__':
    fetcher = Fetcher()
    print(fetcher.get_new_links("https://edition.cnn.com/"))

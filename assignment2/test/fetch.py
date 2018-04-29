import datetime
import sys
import requests
from selenium import webdriver
from selenium.common.exceptions import TimeoutException
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.firefox.options import Options


class Fetcher:
    def __init__(self):
        print("Init")

        # chrome
        chrome_options = Options()
        chrome_options.add_argument("--headless")
        chrome_options.add_argument("--window-size=1920x1080")
        chrome_options.add_argument("--proxy=null")

        self.chrome_driver = webdriver.Chrome(chrome_options=chrome_options)

        # Firefox
        options = Options()
        options.add_argument("--headless")
        self.firefox_driver = webdriver.Firefox(firefox_options=options)

        print("Done")

    def __del__(self):
        self.chrome_driver.quit()
        self.firefox_driver.quit()

    def go(self, driver, type="", url=""):
        start_time = datetime.datetime.now()
        try:
            driver.get(url)
        except TimeoutException as e:
            sys.stderr.write("Timeout " + url + "\n")
            sys.stderr.write(e.msg)
            return

        filename = "{}-{}-{}.png".format(datetime.datetime.today(), driver.title, type)
        self.chrome_driver.save_screenshot(filename)

        end_time = datetime.datetime.now()
        delta = end_time - start_time
        print("get content", delta)

        self.write_source_code_to_file(driver.page_source, type)

    def get_page(self, url):
        # type
        # 0 chrome 1 firefox 2 requests

        # get content
        self.go(self.chrome_driver, "chrome", url)
        self.go(self.firefox_driver, "firefox", url)

        start_time = datetime.datetime.now()

        r = requests.get(url)
        print(r.status_code)
        self.write_source_code_to_file(r.text, "requests")

        end_time = datetime.datetime.now()
        delta = end_time - start_time
        print("get content", delta)

    def write_source_code_to_file(self, page_source, type):
        with open("output-{}.html".format(type), "w") as outputFile:
            outputFile.write(page_source)


if __name__ == '__main__':
    # url = "https://www.ettoday.net/"
    # url = "https://edition.cnn.com/"
    url = "https://news.ycombinator.com/vote?id=16949460&how=up&goto=news"
    fetcher = Fetcher()
    fetcher.get_page(url)

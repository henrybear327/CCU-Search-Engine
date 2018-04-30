import datetime
import sys
import time

from selenium import webdriver
from selenium.common.exceptions import TimeoutException
from selenium.webdriver.chrome.options import Options


class Fetcher:
    def __init__(self):
        print("Init")

        # chrome
        chrome_options = Options()
        chrome_options.add_argument("--headless")
        chrome_options.add_argument("--window-size=1920x1080")
        chrome_options.add_argument("--proxy=null")

        self.chrome_driver = webdriver.Chrome(chrome_options=chrome_options)

        print("Done")

    def __del__(self):
        self.chrome_driver.quit()
        pass

    def get_page(self, url):
        # type
        # 0 chrome 1 firefox 2 requests

        # get content

        start_time = datetime.datetime.now()

        try:
            self.chrome_driver.get(url)
            self.chrome_driver.find_element_by_link_text("Past Contests").click()
            time.sleep(3)
            self.chrome_driver.find_element_by_xpath('//*[@id="wrapper"]/main/div/div/div/table/tr[2]/td[5]/a').click()
            time.sleep(3)
            self.chrome_driver.refresh()
            time.sleep(3)

            idx = 161
            while idx > 0:
                print("scoreboard", idx)
                self.chrome_driver.find_element_by_xpath(
                    '//*[@id="wrapper"]/main/div/div[2]/div[2]/div/ul/li[{}]/a'.format(idx)).click()
                time.sleep(3)

                ret = self.chrome_driver.find_elements_by_partial_link_text("south")
                if len(ret) > 0:
                    print("Found south")
                    return
                idx -= 1
        except TimeoutException as e:
            sys.stderr.write("Timeout " + url + "\n")
            sys.stderr.write(e.msg)
            return

        filename = "{}-{}-{}.png".format(datetime.datetime.today(), self.chrome_driver.title, type)
        self.chrome_driver.save_screenshot(filename)

        end_time = datetime.datetime.now()
        delta = end_time - start_time
        print("get content", delta)

    def write_source_code_to_file(self, page_source):
        with open("output.html".format(type), "w") as outputFile:
            outputFile.write(page_source)


if __name__ == '__main__':
    fetcher = Fetcher()
    url = "https://codejam.withgoogle.com/2018/"
    fetcher.get_page(url)
    # url = "https://codejam.withgoogle.com/2018/challenges/0000000000007764/dashboard"
    # fetcher.get_page(url)
    # url = "https://codejam.withgoogle.com/2018/challenges/0000000000007764/scoreboard"
    # fetcher.get_page(url)

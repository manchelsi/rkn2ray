#!/usr/bin/env pypy3
"""
Парсер https://github.com/zapret-info/z-i/raw/master/dump.csv
для использования в v2ray или xray.
Будет создан файл с правилами domain: full: и regexp:
Дальше перевести его в .dat можно будет при помощи:
https://github.com/v2ray/domain-list-community
Парсить можно и c python3, но c pypy3 раза в 2 бодрее
"""

import csv
import sys
from url_normalize import url_normalize
from urllib.parse import urlparse

csv.field_size_limit(sys.maxsize)

FILENAME = 'z-i/dump.csv'
OUTFILE  = 'zapretinfo'

DOMAINS  = set()
REGEXP   = set()
FULL     = set()

EXCLUDE  = ['youtube.com', 'www.youtube.com']


def exclude(domain=None, url=None):
    """Исключение доменов из списка"""
    urlcheck = ""
    if domain:
        urlcheck = domain
    else:
        urlcheck = urlparse(url).hostname

    for exc in EXCLUDE:
        if urlcheck == exc:
            return True

    return False


def domain_parse(domain):
    """Парсим и заполняем сеты DOMAINS и REGEXP"""
    if domain == "":
        return
    domain = url_normalize(domain)
    domain = urlparse(domain).hostname

    if exclude(domain=domain):
        return

    if ":" in domain:
        return

    if "*." in domain:
        REGEXP.add(domain)
    else:
        DOMAINS.add(domain)


def url_parse(url):
    """Парсим и заполняем сет FULL"""
    if exclude(url=url):
        return
    url = url_normalize(url)
    url = url.partition('//')[2]
    if ":" in url:
        # Здесь теряется ок. 1000 урлов,
        # В основном вида ip:port
        # Но есть и валидные
        # Можно поразбирать еще, но надо ли?
        # В доменах они все равно есть
        return
    FULL.add(url)

def url_split(url):
    """В некоторых строках содержится несколько урлов
    разделим их перед парсером, в домены урлы тоже отправим
    т.к. если есть full:, но нет domain: v2ray не загрузит правила"""
    if url == "":
        return

    url = url.lower()
    if " | " in url:
        urls = url.split(sep = " | ")
        for u in urls:
            domain_parse(u)
            url_parse(u)
    else:
        domain_parse(url)
        url_parse(url)


def read_dump():
    with open(FILENAME, "r", newline="", encoding="windows-1251") as file:
        headers = next(file)
        fieldnames = [None, 'domain', 'url', None, None, None]
        reader = csv.DictReader(file, delimiter=';', fieldnames=fieldnames)
        for row in reader:
            domain_parse(row['domain'])
            url_split(row['url'])


def write_dump():
    with open(OUTFILE, 'w') as zinfo:
        for u in sorted(FULL):
            zinfo.write('full:' + u + '\n')

        for d in sorted(DOMAINS):
            # Если есть www, то добавим домены и без него:
            if d[:4] == 'www.':
                notwww = d.replace('www.', '')
                zinfo.write('domain:' + notwww + '\n')
            zinfo.write('domain:' + d + '\n')

        for r in sorted(REGEXP):
            r = r.replace('.', '\.')
            r = r.replace('*\.', 'regexp:(.*\.)?')
            zinfo.write(r + '\n')


if __name__ == "__main__":
    read_dump()
    write_dump()

#!/usr/bin/env pypy3
# apt install python3-requests
# apt install python3-url-normalize pypy3

import requests
from datetime import datetime,timedelta
import parse
import os

DATEFILE = 'datefile.txt'
URLDATE = 'https://github.com/zapret-info/z-i/commits/master.atom'
URLDUMP = 'https://github.com/zapret-info/z-i/raw/master/dump.csv'
FILEDUMP = 'z-i/dump.csv'

# Похоже, что время коммитов -1h, для соответствия MSK
# добавить 4h
TMZ = timedelta(hours=4)

def newdate():
    # Получить дату, z-i обновляет при новых коммитах
    try:
        r = requests.get(URLDATE, allow_redirects=True)
        for i in r.text.split(sep='\n'):
            if 'Updated:' in i:
                NEWDATE = ' '.join(i.split()[1:3])
                NEWDATE = datetime.fromisoformat(NEWDATE) + TMZ
                break
        return NEWDATE
    except:
        print("Не удалось получить дату")
        exit(0)


def df_wr(NEWDATE):
    # запись даты в файл
    with open(DATEFILE, 'w') as df:
        df.write(str(NEWDATE))


def olddate():
    # взять дату из файла, либо перезаписать 
    # с новой датой
    try:
        with open(DATEFILE) as df:
            OLDDATE = df.read().replace('\n', '')
            OLDDATE = datetime.fromisoformat(OLDDATE)
        return OLDDATE
    except:
        oldoldate = '2023-08-07 14:35:01'
        df_wr(oldoldate)
        return olddate()


def upgrade_dat(ndate):
    try:
        if not os.path.exists("z-i"):
            os.makedirs("z-i")
        print("Скачивается новый дамп")
        r = requests.get(URLDUMP, allow_redirects=True)
        open(FILEDUMP, 'wb').write(r.content)
        df_wr(ndate)
        print("Парсим дамп")
        parse.read_dump()
        print("Запись zapretinfo")
        parse.write_dump()
        print("Запись zapretinfo.dat")
        os.system('./domain-list-community')
    except:
        print("Не удалось скачать дамп")
        exit(0)


def main():
    fdate = olddate()
    ndate = newdate()

    if fdate < newdate():
        upgrade_dat(ndate)
    else:
        print("Новых обновлений нет, было обновлено", fdate)
        exit(0)


if __name__ == "__main__":
        main()

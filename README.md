# rkn2ray

Парсим и генерируем zapretinfo.dat для xray, или v2ray
из дампа https://github.com/zapret-info/z-i/raw/master/dump.csv

dat создается тулзой https://github.com/v2fly/domain-list-community,
собирается "go build"


./run.py скачивает dump.csv, если он обновился, парсит и создает zapretinfo.dat, который можно использовать в xray
примерно так:

```
        "routing": {                   
            "domainStrategy": "AsIs",
            "rules": [                
            {         
                "type": "field",                                                                         
                "domainMatcher": "mph",
                "domains": [           
                    "ext:zapretinfo.dat:zapretinfo"
                ],
                "outboundTag": "allow"
            },       

```

Python modules: requests, url-normalize, PyYAML

Для исключения доменов, добавить нужные домены в config.yaml, секцию EXCLUDES

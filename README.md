# rkn2ray

Парсим и генерируем zapretinfo.dat для xray, или v2ray
из дампа https://github.com/zapret-info/z-i/raw/master/dump.csv


./rkn2ray скачивает dump.csv, если он обновился, парсит и создает zapretinfo.dat, который можно использовать в xray
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

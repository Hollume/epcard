import shutil
import json
import requests

url = 'http://localhost:5050/card'
data = {
    "rarity":   "0",
    "is_trash": True,
    "owned":    False,
    "catching": "줄무늬고등어가 낚였다!",
    "price":    "502$",
    "tax":      "-35$ (-7%)",
    "bonus":    "+0$ (+0%)",
    "money":    "467$",
    "name":     "줄무늬고등어",
    "detail":   "40.15cm\n(평균 40cm)\n502$\n(평균 500$)",
    "place":    "2022-03-01 02시에 '『국영ㆍ낚시터』'에서\n『123』"
}
response = requests.get(url, json=data, stream=True)
print(response.status_code)
with open('img.png', 'wb') as out_file:
    shutil.copyfileobj(response.raw, out_file)
del response
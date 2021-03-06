import requests
import json
import time

# Create account [should return code 201]
link = {
#    "extradata": "ignoreme1",
    "id"       : "ignoreme2",
    "url"      : "https://www.google.com"
}
print("Attempting to create link... [Should return 201]")
res = requests.post("http://localhost:8080/url", json=link)
print(res)
print(res.json())

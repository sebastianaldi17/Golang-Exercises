import requests
import json
import time

# Create link [should return code 201]
link = {
#    "extradata": "ignoreme1",
    "id"       : "ignoreme2",
    "url"      : "https://www.google.com"
}
# print("Attempting to create link... [Should return 201]")
# res = requests.post("http://localhost:8080/url", json=link)
# print(res)
# print(res.json())

# time.sleep(1)

# Get all links
links = requests.get("http://localhost:8080/url")
print(links)
print(links.json())
for k, v in links.json().items():
    print(k, v)
    if v['url'] == "https://www.google.com":
        # Delete link
        # Should return 200 for successful deletion
        print("Attempting to DELETE created link... [Should return 200]")
        res = requests.delete("http://localhost:8080/url/"+v['id'])
        print(res)
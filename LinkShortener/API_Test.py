import requests
import json
import time

# Get all links [should return code 200]
print("Attempting to get all links... [Should return 200]")
res = requests.get("http://localhost:8080/url")
print(res)
print(res.json())

time.sleep(1)

# Create link [should return code 201]
link = {
    # "extradata": "ignoreme",
    "id"       : "customshort",
    "url"      : "https://www.google.com"
}
print("Attempting to create link... [Should return 201]")
res = requests.post("http://localhost:8080/url", json=link)
print(res)

time.sleep(1)

# Delete link [should return code 200]
print("Attempting to DELETE created link... [Should return 200]")
res = requests.delete("http://localhost:8080/url/customshort")
print(res)
import requests
import json
import time

# Create account [should return code 201]
account = {
    "id": 12,
    "first_name": "Se",
    "last_name": "Al",
    "user_name": "Nald"
}
print("Attempting to create account... [Should return 201]")
res = requests.post("http://localhost:8080/accounts", json=account)
print(res)
print(res.json())
time.sleep(1)

# Create second account [should return code 201]
account = {
    "id": 11,
    "first_name": "Se",
    "last_name": "Al",
    "user_name": "Nald"
}
print("Attempting to create second account... [Should return 201]")
res = requests.post("http://localhost:8080/accounts", json=account)
print(res)
print(res.json())
time.sleep(1)

# Create broken account (incomplete field) [should return error code 400]
account = {
    "id": 123,
    "first_name": "Se",
    "user_name": "Nald"
}
print("Attempting to create incomplete account... [Should return 400]")
res = requests.post("http://localhost:8080/accounts", json=account)
print(res)
time.sleep(1)

# Create broken account (too much information) [should return error code 400]
account = {
    "id": 1234,
    "first_name": "Seb",
    "last_name": "Ald",
    "user_name": "Axy",
    "loren": "ipsum"
}
print("Attempting to create account with too much information... [Should return 400]")
res = requests.post("http://localhost:8080/accounts", json=account)
print(res)
time.sleep(1)

# Get the new account [should return code 200 and data]
print("Attempting to GET created account... [Should return 200 with json containing account details]")
res = requests.get("http://localhost:8080/accounts/12")
print(res)
print(res.json())
time.sleep(1)

# Should return 200 for successful deletion
print("Attempting to DELETE created account... [Should return 200]")
res = requests.delete("http://localhost:8080/accounts/12")
print(res)
time.sleep(1)

# Try deleting nonexistant ID (should be 404)
print("Attempting to DELETE nonexistant account... [Should return 404]")
res = requests.delete("http://localhost:8080/accounts/999")
print(res)
time.sleep(1)

# Try getting deleted ID (should be 404)
print("Attempting to DELETE previously deleted account... [Should return 404]")
res = requests.get("http://localhost:8080/accounts/12")
print(res)
time.sleep(1)
import requests

url = "http://localhost:8000/the_good_stuff"
payload = {
    "name": "session1"
}

if __name__ == '__main__':
    req = requests.post(url, payload)

    req.raise_for_status()
    
    print(req.content)
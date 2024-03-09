
import requests

url = "http://localhost:8080/blogs"
headers = {
        "Content-Type": "application/json",
        "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDk1NDcwMTYsInJvbGUiOiJ1c2VyIiwidXNlcm5hbWUiOiJoLnRpd2FyaS5kZXYifQ.o2PYcosVieynIYexNWF7-aEJ2uKRwvBV_B49fjilzvA"
        }

data = {
        "title": "Tet title",
        "description": "Descrption Test Article",
        "content": open("../blogs/md/pydantic.md").read()
        }

response = requests.post(url, json=data, headers=headers)

print(response)


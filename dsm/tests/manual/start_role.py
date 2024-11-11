#! /usr/bin/env python3

import json
import requests

def send_request(url, method, data):
    requests.request(method, url, json=data)

url = "http://localhost:50551/start"
method = "POST"
data = {"roleId": "role-1", "serviceId": "service-1", "imageId": "busybox:latest"}
send_request(url, method, data)
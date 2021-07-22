import JointCloudStorage
import requests

if __name__ == '__main__':
    JointCloudStorage.test()
    access_key = "47bed2954a3647fdbc7a3364778c388f"
    secret_key = "2d61f1ec5d834af6ad7150d756fa33e5"
    user = JointCloudStorage.auth.Auth(access_key,secret_key)

    # request session
    request_session = requests.Session()
    with open("a.txt", 'rb') as this_file:
        req = requests.Request('GET', url="http://192.168.105.13:8085/object/")
        pre_request = request_session.prepare_request(req)
        var = user._sign_request(pre_request)
        # send request
        resp = request_session.send(pre_request)
        print(resp.text)

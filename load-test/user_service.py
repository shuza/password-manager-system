from locust import HttpLocust, TaskSet, task

class UserBehavior(TaskSet):
    def login(self):
        body = {
            "email": "user-1@gmail.com"
            "password": "123456"
        }
        self.client.post("/user/api/v1/auth/login", body)


class WebsiteUser(HttpLocust):
    task_set = UserBehavior

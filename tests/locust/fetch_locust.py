from locust import HttpUser, task

class HelloWorldUser(HttpUser):
    @task
    def fetch(self):
        self.client.get("/fetch")

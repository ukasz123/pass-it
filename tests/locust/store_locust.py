from locust import HttpUser, task
import random

class HelloWorldUser(HttpUser):
    @task
    def fetch(self):
        r = random.randrange(10000000)
        self.client.put("/store/{id}".format(id=r), data={"key": "MIGJAoGBAOsrU1YavLSGb7Fsnw5RCgaBQdO4BmiDmkJIbDc5WviCQzP+e8HTBNC9jseiNBG0yW48LexbMhpvt/rKNR0D7ACtlShzCu8YdbcP0GIn2lobxvP7Ne9tfTYeQLwXsQOX2aFM4ysqAhijWNlKUUq31N1tsavM+H7mdlouOpmmHLnDAgMBAAE=", "payload": "Something"})

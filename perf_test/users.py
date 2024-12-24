import json
# Количество сущностей для создания
num_entities = 100000


for i in range(num_entities):
        with open(f"targets/target{i}.json", "w") as f:
                entity = {
                        "first_name": f"First{i}",
                        "last_name": f"Last{i}",
                        "email": f"test{i}@mail.com",
                        "password": "qwerty123"
                }
                f.write(json.dumps(entity))

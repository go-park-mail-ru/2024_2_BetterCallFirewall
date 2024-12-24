import json

# Количество сущностей для создания
num_entities = 100000

# Открываем файл для записи
with open('entities.txt', 'w') as f:
    for i in range(num_entities):
        f.write(f"POST http://185.241.194.197:8072/api/v1/auth/register\n")
        f.write("Content-Type: application/json\n")
        f.write(f"@./targets/target{i}.json\n\n")

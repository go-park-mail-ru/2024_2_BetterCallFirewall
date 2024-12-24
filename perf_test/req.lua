math.randomseed(os.time())
local counter = 0
-- Генерируем уникального пользователя
function request()
    local user_id = math.random(1000, 999999)
    local first_name = "fname" .. user_id .. counter
    local last_name = "lname" .. user_id .. counter
    local email = first_name .. last_name  .. "@example.com"
    local password = "qwerty12345"

    -- Тело запроса
    local body = string.format('{"first_name":"%s", "last_name":"%s", "password":"%s", "email":"%s"}', first_name, last>    counter = counter + 1
    return wrk.format("POST", "/api/v1/auth/register", {["Content-Type"] = "application/json"}, body)
end

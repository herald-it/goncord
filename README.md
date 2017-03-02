goncord
-------

[![Build Status](https://travis-ci.org/herald-it/goncord.svg?branch=master)](https://travis-ci.org/herald-it/goncord)

Для запуска необходим файл настроек.

Пример файла: 
```yaml
database:
  host: localhost
  dbname: testdb
  tokentable: tokentable
  usertable: usertable
ssl:
  key: ./pass_key
  certificate: ./pass_certificate
router:
  register: 
    path: /register
  login: 
    path: /login
  validate: 
    path: /validate
  logout: 
    path: /logout
  update: 
    path: /update
  resetpassword:
    path: /reset
domain: my.domain.com
ip: 0.0.0.0:8000
```
Секция `ssl` не является обязательной.

Для каждого роута можно настроить доступ только определенного круга лиц.
Например:
```yaml
router:
  register:
    path: /register
    allowedhost:
      - localhost
      - 192.168.0.2
```
#Аргументы запросов
> \* помечены обязательные поля.
-------

##Register - [POST]
`login*` - логин.

`email*` - почтовый ящик.

`password*` - пароль.


##Login - [POST]
`login*` - логин.

`email*` - почтовый ящик.

`password*` - пароль.

Можно отправлять только логин или только почтовый ящик.

##Validate - [POST]
Запрос на данный адрес вернет модель пользователя в формате json.

##Logout - [POST|GET]
Любой запрос на данный адрес удалить текущий токен из БД.
В результате этого пользователь больше не сможет пройти валидацию
и будет считаться не активным.

##Update - [POST]
Принимает модель пользователя в формате json.
Поля: ID, login, email, password обновить невозможно.
Обновить можно только поле: payloads.
Пример запроса на обновление:
```json
user: {"payload":"{'foo': 'bar'}"}
```

##Reset password - [POST]
Сброс пароля.

`new_password*` - новый пароль.

`old_password` - старый пароль.

Если в запросе присутствует поле `old_password`, то будет произведена проверка на соответствие, иначе поле пароля будет заменено на `new_password`.
